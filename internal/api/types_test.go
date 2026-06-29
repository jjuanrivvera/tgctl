package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_UnmarshalVariants(t *testing.T) {
	cases := map[string]struct {
		in   string
		want ID
		err  bool
	}{
		"number":           {`123`, "123", false},
		"string":           {`"456"`, "456", false},
		"large int64":      {`7240817235823`, "7240817235823", false},
		"beyond 2^53":      {`9007199254740993`, "9007199254740993", false},
		"negative chat":    {`-1001234567890`, "-1001234567890", false},
		"null":             {`null`, "", false},
		"float rejected":   {`1.5`, "", true},
		"exp rejected":     {`1e9`, "", true},
		"garbage rejected": {`{}`, "", true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var id ID
			err := json.Unmarshal([]byte(tc.in), &id)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, id)
		})
	}
}

func TestID_MarshalAlwaysString(t *testing.T) {
	b, err := json.Marshal(ID("9007199254740993"))
	require.NoError(t, err)
	assert.Equal(t, `"9007199254740993"`, string(b))
}

func TestInt_UnmarshalVariants(t *testing.T) {
	var n Int
	require.NoError(t, json.Unmarshal([]byte(`42`), &n))
	assert.Equal(t, int64(42), n.Int64())

	require.NoError(t, json.Unmarshal([]byte(`"99"`), &n))
	assert.Equal(t, int64(99), n.Int64())

	require.NoError(t, json.Unmarshal([]byte(`null`), &n))
	assert.Equal(t, int64(0), n.Int64())

	require.Error(t, json.Unmarshal([]byte(`3.14`), &n))
	require.Error(t, json.Unmarshal([]byte(`"abc"`), &n))

	b, _ := json.Marshal(Int(7))
	assert.Equal(t, "7", string(b))
}

func TestBool_UnmarshalVariants(t *testing.T) {
	for _, in := range []string{`true`, `"true"`, `"1"`, `"yes"`} {
		var b Bool
		require.NoError(t, json.Unmarshal([]byte(in), &b))
		assert.True(t, bool(b), in)
	}
	for _, in := range []string{`false`, `"false"`, `"0"`, `"no"`, `null`} {
		var b Bool
		require.NoError(t, json.Unmarshal([]byte(in), &b))
		assert.False(t, bool(b), in)
	}
	var b Bool
	require.Error(t, json.Unmarshal([]byte(`"maybe"`), &b))
}

// FuzzID ensures the flexible decoder never panics on arbitrary input and that any value it
// accepts round-trips back to a quoted string without error.
func FuzzID(f *testing.F) {
	for _, s := range []string{`1`, `"2"`, `null`, `-100`, `9007199254740993`, `1.5`, `"x"`, `[]`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		var id ID
		if err := json.Unmarshal([]byte(in), &id); err != nil {
			return // rejecting bad input is fine; we only require no panic
		}
		if _, err := json.Marshal(id); err != nil {
			t.Fatalf("accepted %q but failed to marshal: %v", in, err)
		}
	})
}

// FuzzInt mirrors FuzzID for the integer decoder.
func FuzzInt(f *testing.F) {
	for _, s := range []string{`0`, `"7"`, `null`, `-9`, `3.14`, `"NaN"`, `1e9`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		var n Int
		_ = json.Unmarshal([]byte(in), &n) // must not panic
	})
}
