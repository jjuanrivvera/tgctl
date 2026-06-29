package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Client is the single typed gateway to the Telegram Bot API. Every command goes through it,
// so auth, retries, rate limiting, and dry-run live in exactly one place (GOAL.md §2).
type Client struct {
	auth    Authenticator
	baseURL string
	http    *http.Client
	limiter *rateLimiter
	retry   retryPolicy

	DryRun    bool
	ShowToken bool
	Verbose   bool

	dryRunW  io.Writer // where the --dry-run curl line is written (default os.Stderr)
	verboseW io.Writer
}

// Option configures a Client.
type Option func(*Client)

// DefaultBaseURL is the public Bot API host. Override for a self-hosted Local Bot API Server.
const DefaultBaseURL = "https://api.telegram.org"

// New builds a Client for the given authenticator. Sensible defaults are applied; override
// with Options.
func New(auth Authenticator, opts ...Option) *Client {
	c := &Client{
		auth:     auth,
		baseURL:  DefaultBaseURL,
		http:     &http.Client{Timeout: 0}, // no hard timeout: getUpdates long-polls; ctx governs cancellation
		limiter:  newRateLimiter(25),       // Telegram tolerates ~30 msg/s overall; stay under it
		retry:    defaultRetryPolicy(),
		dryRunW:  os.Stderr,
		verboseW: os.Stderr,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func WithBaseURL(u string) Option {
	return func(c *Client) {
		if u != "" {
			c.baseURL = u
		}
	}
}
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }
func WithRPS(rps float64) Option           { return func(c *Client) { c.limiter = newRateLimiter(rps) } }
func WithRetryPolicy(p retryPolicy) Option { return func(c *Client) { c.retry = p } }
func WithDryRun(v bool) Option             { return func(c *Client) { c.DryRun = v } }
func WithShowToken(v bool) Option          { return func(c *Client) { c.ShowToken = v } }
func WithVerbose(v bool) Option            { return func(c *Client) { c.Verbose = v } }
func WithDryRunWriter(w io.Writer) Option  { return func(c *Client) { c.dryRunW = w } }

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string { return c.baseURL }

// apiResponse is the Bot API envelope shared by every method.
type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
	Parameters  *RespParams     `json:"parameters"`
}

// Call invokes a Bot API method with JSON params and returns the raw `result`. idempotent
// marks read-only methods so the retry policy may safely replay them on ambiguous failures.
// In dry-run mode it prints the equivalent curl and returns (nil, nil).
func (c *Client) Call(ctx context.Context, method string, params map[string]any, idempotent bool) (json.RawMessage, error) {
	body, err := jsonBody(params)
	if err != nil {
		return nil, err
	}
	req := &preparedRequest{
		method:      method,
		contentType: "application/json",
		body:        body,
		idempotent:  idempotent,
		curlData:    curlJSONArg(params),
	}
	return c.do(ctx, req)
}

// CallInto is Call plus a JSON decode of the result into out.
func (c *Client) CallInto(ctx context.Context, method string, params map[string]any, idempotent bool, out any) error {
	raw, err := c.Call(ctx, method, params, idempotent)
	if err != nil {
		return err
	}
	if len(raw) == 0 || out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

// Upload invokes a method with multipart/form-data. files maps a Bot API field (e.g. "photo")
// to a LOCAL file path; params carries the remaining scalar fields. File paths are validated
// (existence + regular file) by the caller's path-confinement helper before reaching here.
func (c *Client) Upload(ctx context.Context, method string, params map[string]any, files map[string]string, idempotent bool) (json.RawMessage, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for k, v := range params {
		if err := mw.WriteField(k, scalarString(v)); err != nil {
			return nil, err
		}
	}
	var curlParts []string
	for field, path := range files {
		f, err := os.Open(path) //nolint:gosec // G304: path is confined by the caller (see confinePath)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", field, err)
		}
		fw, err := mw.CreateFormFile(field, filepath.Base(path))
		if err != nil {
			_ = f.Close()
			return nil, err
		}
		if _, err := io.Copy(fw, f); err != nil {
			_ = f.Close()
			return nil, err
		}
		_ = f.Close()
		curlParts = append(curlParts, fmt.Sprintf("-F %s=@%s", field, shellQuote(path)))
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	req := &preparedRequest{
		method:      method,
		contentType: mw.FormDataContentType(),
		body:        buf.Bytes(),
		idempotent:  idempotent,
		curlData:    strings.Join(append(curlFields(params), curlParts...), " "),
	}
	return c.do(ctx, req)
}

type preparedRequest struct {
	method      string
	contentType string
	body        []byte
	idempotent  bool
	curlData    string // the rendered curl data args, for --dry-run
}

// do executes a prepared request with rate limiting and the retry policy. It is the single
// network path: dry-run, backoff, 429 adaptation, and envelope parsing all live here.
func (c *Client) do(ctx context.Context, req *preparedRequest) (json.RawMessage, error) {
	url := c.auth.RequestURL(c.baseURL, req.method)

	if c.DryRun {
		c.printCurl(req)
		return nil, nil
	}

	var lastErr error
	for attempt := 0; attempt < c.retry.maxAttempts; attempt++ {
		if err := c.limiter.wait(ctx); err != nil {
			return nil, err
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(req.body))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", req.contentType)

		resp, err := c.http.Do(httpReq) //nolint:bodyclose // body is read+closed in parse()
		if err != nil {
			lastErr = fmt.Errorf("%s: %w", req.method, err)
			if retry, wait := c.retry.decide(attempt, 0, 0, true, req.idempotent); retry {
				if serr := sleepCtx(ctx, wait); serr != nil {
					return nil, serr
				}
				continue
			}
			return nil, lastErr
		}

		result, apiErr, retryAfter := c.parse(resp, req.method)
		if apiErr == nil {
			c.limiter.reward()
			return result, nil
		}
		lastErr = apiErr
		if apiErr.Code == 429 {
			c.limiter.penalize()
		}
		if retry, wait := c.retry.decide(attempt, apiErr.Code, retryAfter, false, req.idempotent); retry {
			if serr := sleepCtx(ctx, wait); serr != nil {
				return nil, serr
			}
			continue
		}
		return nil, lastErr
	}
	return nil, lastErr
}

// parse reads the response body, decodes the Bot API envelope, and turns a non-ok response
// into a typed *APIError. retryAfter is surfaced separately so the retry loop can honor it.
func (c *Client) parse(resp *http.Response, method string) (json.RawMessage, *APIError, int) {
	defer func() { _ = resp.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<20)) // cap at 32MiB to bound memory

	if c.Verbose {
		_, _ = fmt.Fprintf(c.verboseW, "← %s %d %s\n", method, resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var env apiResponse
	if err := json.Unmarshal(raw, &env); err != nil {
		// Non-JSON (e.g. an HTML 502 from a proxy): synthesize an APIError off the status.
		return nil, &APIError{
			StatusCode:  resp.StatusCode,
			Code:        resp.StatusCode,
			Description: fmt.Sprintf("non-JSON response: %s", truncate(strings.TrimSpace(string(raw)), 200)),
			Body:        string(raw),
			Method:      method,
		}, retryAfterHeader(resp)
	}

	if env.OK {
		return env.Result, nil, 0
	}

	ae := &APIError{
		StatusCode:  resp.StatusCode,
		Code:        env.ErrorCode,
		Description: env.Description,
		Parameters:  env.Parameters,
		Body:        string(raw),
		Method:      method,
	}
	ra := ae.RetryAfter()
	if ra == 0 {
		ra = retryAfterHeader(resp)
	}
	return nil, ae, ra
}

// printCurl writes the equivalent curl command for --dry-run. The token is redacted unless
// --show-token is set, so a dry-run is always safe to paste into a bug report.
func (c *Client) printCurl(req *preparedRequest) {
	url := c.auth.RedactedURL(c.baseURL, req.method)
	if c.ShowToken {
		url = c.auth.RequestURL(c.baseURL, req.method)
	}
	var b strings.Builder
	b.WriteString("curl -sS -X POST ")
	b.WriteString(shellQuote(url))
	if req.contentType == "application/json" {
		b.WriteString(" -H 'Content-Type: application/json'")
	}
	if req.curlData != "" {
		b.WriteByte(' ')
		b.WriteString(req.curlData)
	}
	_, _ = fmt.Fprintln(c.dryRunW, b.String())
}

func retryAfterHeader(resp *http.Response) int {
	if v := resp.Header.Get("Retry-After"); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return 0
}

// --- small helpers ---

func jsonBody(params map[string]any) ([]byte, error) {
	if len(params) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(params)
}

func curlJSONArg(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}
	b, _ := json.Marshal(params)
	return "-d " + shellQuote(string(b))
}

func curlFields(params map[string]any) []string {
	out := make([]string, 0, len(params))
	for k, v := range params {
		out = append(out, fmt.Sprintf("-F %s=%s", k, shellQuote(scalarString(v))))
	}
	return out
}

func scalarString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case nil:
		return ""
	case json.RawMessage:
		return string(t)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// shellQuote single-quote-escapes a string so a dry-run curl line is copy-pasteable and safe.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
