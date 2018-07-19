package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSlice(t *testing.T) {
	exp := []string{"bash", "-c", "sleep 1"}
	src := `bash -c 'sleep 1'`

	res, err := StringToSlice(src)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestParseSliceNoEnv(t *testing.T) {
	src := `bash -c 'sleep ${SLEEP}'`

	res, err := StringToSlice(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestParseSliceNoSecondQuote(t *testing.T) {
	src := `bash -c 'sleep`

	res, err := StringToSlice(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestParseSliceWithEnv(t *testing.T) {
	exp := []string{"hello", "hello123", "world", "hello-world", "hello-joe"}
	env := map[string]string{"USER": "joe"}

	src := `'hello' hello123 "world" \
	"hello-world" hello-${USER}`

	res, err := StringToSliceWithEnv(src, env)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func BenchmarkParseSliceString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StringToSliceWithEnv(`'hello' hello123 "world"`, nil)
	}
}
