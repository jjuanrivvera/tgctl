package commands

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/tgctl/internal/store"
)

// webhookListenCmd runs a local HTTP server that receives Telegram webhook updates and prints
// each one. It's a value-add beyond the raw API (GOAL.md §3c): Telegram pushes updates to a
// public HTTPS URL, so front this with a tunnel (cloudflared/ngrok) or a public host, point
// the webhook at it (manually or with --set-url), and watch updates stream in.
func webhookListenCmd() *cobra.Command {
	var (
		port         int
		path         string
		secret       string
		setURL       string
		deleteOnExit bool
	)
	cmd := &cobra.Command{
		Use:   "listen",
		Short: "Run a local server that receives webhook updates and prints them",
		Long: `Start an HTTP server that receives Telegram webhook updates (the Update objects
Telegram POSTs to your webhook URL) and renders each one with the usual -o/--output.

Telegram only delivers to a public HTTPS endpoint, so put this behind a tunnel
(cloudflared, ngrok) or run it on a public host. Point the webhook at that URL either
yourself (tgctl webhook set --url ...) or with --set-url here. When --secret-token is set,
requests whose X-Telegram-Bot-Api-Secret-Token header doesn't match are rejected. Ctrl-C
stops the server (and deletes the webhook if --delete-on-exit).`,
		Example: `  tgctl webhook listen --port 8080 -o json
  tgctl webhook listen --port 8080 --secret-token s3cr3t \
      --set-url https://my-tunnel.example.com/bot --delete-on-exit
  # local test (no bot needed):
  curl -XPOST localhost:8080 -d '{"update_id":1,"message":{"message_id":7,"text":"hi"}}'`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if setURL != "" {
				client, err := clientFromCmd(cmd)
				if err != nil {
					return err
				}
				// Registered before the deleteOnExit defer below, so LIFO ordering runs the
				// webhook cleanup call first and closes the client (and its store handle) last.
				defer func() { _ = client.Close() }()
				params := map[string]any{"url": setURL}
				if secret != "" {
					params["secret_token"] = secret
				}
				if _, err := client.Call(ctx, "setWebhook", params, false); err != nil {
					return err
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "webhook set to %s\n", setURL)
				if deleteOnExit {
					// ctx is cancelled on Ctrl-C, so clean up on a fresh (non-cancelled) context.
					defer func() {
						dctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
						defer cancel()
						if _, err := client.Call(dctx, "deleteWebhook", nil, false); err != nil {
							fmt.Fprintf(cmd.ErrOrStderr(), "webhook cleanup failed: %v\n", err)
							return
						}
						fmt.Fprintln(cmd.ErrOrStderr(), "webhook deleted")
					}()
				}
			}

			// Best-effort, opened once for the life of the server (not per-request): a nil
			// store (disabled via --no-store, or unavailable) simply means inbound updates
			// aren't recorded, exactly like the read failure path elsewhere (DECISIONS.md).
			st := openStoreForWrite(cmd)
			if st != nil {
				defer func() { _ = st.Close() }()
			}

			mux := http.NewServeMux()
			mux.HandleFunc(routePath(path), webhookUpdateHandler(cmd, secret, st))
			srv := &http.Server{
				Addr:              fmt.Sprintf(":%d", port),
				Handler:           mux,
				ReadHeaderTimeout: 10 * time.Second,
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "listening on :%d%s (Ctrl-C to stop)\n", port, routePath(path))
			errCh := make(chan error, 1)
			go func() { errCh <- srv.ListenAndServe() }()

			select {
			case <-ctx.Done():
				shutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
				defer cancel()
				_ = srv.Shutdown(shutCtx)
				fmt.Fprintln(cmd.ErrOrStderr(), "stopped")
				return nil
			case err := <-errCh:
				if err == http.ErrServerClosed {
					return nil
				}
				return err
			}
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")
	cmd.Flags().StringVar(&path, "path", "/", "URL path to receive updates on")
	cmd.Flags().StringVar(&secret, "secret-token", "", "require this X-Telegram-Bot-Api-Secret-Token header")
	cmd.Flags().StringVar(&setURL, "set-url", "", "register the webhook at this public HTTPS URL before listening")
	cmd.Flags().BoolVar(&deleteOnExit, "delete-on-exit", false, "delete the webhook when the server stops (requires --set-url)")
	return cmd
}

func routePath(p string) string {
	if p == "" {
		return "/"
	}
	return p
}

// webhookUpdateHandler returns the HTTP handler that validates and renders one update. Writes
// to stdout are serialized so concurrent deliveries don't interleave. It acknowledges Telegram
// with 200 immediately on a valid request, then renders. When st is non-nil, an incoming
// message is also recorded (direction 'in') to the local store — the webhook counterpart to
// commands/updates.go's recordInboundUpdates.
func webhookUpdateHandler(cmd *cobra.Command, secret string, st *store.Store) http.HandlerFunc {
	var mu sync.Mutex
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Constant-time compare so the secret check isn't a timing oracle.
		if secret != "" {
			got := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, 5<<20))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if !json.Valid(body) {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK) // ack Telegram fast; rendering must not delay the 200

		mu.Lock()
		defer mu.Unlock()
		if st != nil {
			var update struct {
				Message *telegramMessage `json:"message"`
			}
			if err := json.Unmarshal(body, &update); err == nil && update.Message != nil {
				recordInboundMessage(cmd, st, update.Message)
			}
		}
		if err := render(cmd, json.RawMessage(body)); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "render error: %v\n", err)
		}
	}
}
