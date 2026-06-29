package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func render(t *testing.T, data string, opts Options) (string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	opts.Out = &out
	opts.Err = &errb
	require.NoError(t, Render(json.RawMessage(data), opts))
	return out.String(), errb.String()
}

func TestRender_Table_Object(t *testing.T) {
	out, _ := render(t, `{"message_id":42,"text":"hi","chat":{"id":7,"title":"Room"}}`,
		Options{Format: FormatTable, NoColor: true})
	// Header is uppercased; preferred keys lead; nested object flattens to dotted keys.
	assert.Contains(t, out, "MESSAGE_ID")
	assert.Contains(t, out, "CHAT.ID")
	assert.Contains(t, out, "42")
	assert.Contains(t, out, "Room")
	// message_id should appear before text (preferred order).
	assert.Less(t, strings.Index(out, "MESSAGE_ID"), strings.Index(out, "TEXT"))
}

func TestRender_Table_Array(t *testing.T) {
	out, _ := render(t, `[{"command":"start","description":"Begin"},{"command":"help","description":"Help"}]`,
		Options{Format: FormatTable, NoColor: true})
	assert.Contains(t, out, "COMMAND")
	assert.Contains(t, out, "start")
	assert.Contains(t, out, "help")
}

func TestRender_Table_Scalar(t *testing.T) {
	out, _ := render(t, `true`, Options{Format: FormatTable, NoColor: true})
	assert.Equal(t, "true\n", out)
}

func TestRender_JSON_PreservesBigInt(t *testing.T) {
	out, _ := render(t, `{"id":9007199254740993}`, Options{Format: FormatJSON})
	assert.Contains(t, out, "9007199254740993", "big int must not lose precision")
}

func TestRender_YAML(t *testing.T) {
	out, _ := render(t, `{"id":42,"name":"x"}`, Options{Format: FormatYAML})
	assert.Contains(t, out, "id: 42")
	assert.Contains(t, out, "name: x")
}

func TestRender_CSV(t *testing.T) {
	out, _ := render(t, `[{"command":"a","description":"d1"},{"command":"b","description":"d2"}]`,
		Options{Format: FormatCSV})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	assert.Equal(t, "command,description", lines[0])
	assert.Equal(t, "a,d1", lines[1])
}

func TestRender_CSV_FormulaInjection(t *testing.T) {
	out, _ := render(t, `[{"text":"=SUM(A1:A2)"},{"text":"+1"},{"text":"@cmd"},{"text":"-7"},{"text":"-bad"}]`,
		Options{Format: FormatCSV, Columns: []string{"text"}})
	assert.Contains(t, out, "'=SUM(A1:A2)")
	assert.Contains(t, out, "'+1")
	assert.Contains(t, out, "'@cmd")
	assert.Contains(t, out, "\n-7\n", "a real negative number is left alone")
	assert.Contains(t, out, "'-bad", "a leading - that isn't numeric is neutralized")
}

func TestRender_ID(t *testing.T) {
	out, _ := render(t, `[{"update_id":100,"message":{}},{"update_id":101,"message":{}}]`,
		Options{Format: FormatID})
	assert.Equal(t, "100\n101\n", out)
}

func TestRender_JQ(t *testing.T) {
	out, _ := render(t, `{"result":{"username":"bot"}}`,
		Options{Format: FormatJSON, JQ: ".result.username"})
	assert.Contains(t, out, `"bot"`)
}

func TestRender_JQ_Invalid(t *testing.T) {
	var o, e bytes.Buffer
	err := Render(json.RawMessage(`{}`), Options{Format: FormatJSON, JQ: ".[", Out: &o, Err: &e})
	require.Error(t, err)
}

func TestRender_ColumnsExplicit(t *testing.T) {
	out, _ := render(t, `{"a":1,"b":2,"c":3}`,
		Options{Format: FormatTable, NoColor: true, Columns: []string{"c", "a"}})
	// Only c and a, in that order.
	assert.Less(t, strings.Index(out, "C"), strings.Index(out, "A"))
	assert.NotContains(t, out, "B ")
}

func TestRender_ColumnCap(t *testing.T) {
	_, errb := render(t, `{"a":1,"b":2,"c":3,"d":4}`,
		Options{Format: FormatTable, NoColor: true, MaxCols: 2})
	assert.Contains(t, errb, "note:")
}

func TestRender_EmptyResult(t *testing.T) {
	out, _ := render(t, ``, Options{Format: FormatTable})
	assert.Empty(t, out)
}

func TestFormat_Valid(t *testing.T) {
	assert.True(t, FormatTable.Valid())
	assert.True(t, Format("json").Valid())
	assert.False(t, Format("xml").Valid())
}

func TestTruncCell(t *testing.T) {
	s, cut := truncCell("hello world", 5)
	assert.True(t, cut)
	assert.Equal(t, "hell…", s)
	s, cut = truncCell("hi", 5)
	assert.False(t, cut)
	assert.Equal(t, "hi", s)
}
