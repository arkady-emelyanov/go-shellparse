package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringToSlice(t *testing.T) {
	exp := []string{"bash", "-c", "sleep 1"}
	src := `bash -c 'sleep 1'`

	res, err := StringToSlice(src)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestStringToSliceNoVars(t *testing.T) {
	src := `bash -c 'sleep ${SLEEP}'`

	res, err := StringToSlice(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestStringToSliceNoSecondQuote(t *testing.T) {
	src := `bash -c 'sleep`

	res, err := StringToSlice(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestStringToSliceNoRightDelim(t *testing.T) {
	src := `${HELLO`

	res, err := StringToSlice(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestStringToSliceWithVars(t *testing.T) {
	exp := []string{"hello", "hello123", "world", "hello-world", "hello-joe"}
	vars := map[string]string{"USER": "joe"}

	src := `'hello' hello123 "world" \
	"hello-world" hello-${USER}`

	res, err := StringToSliceWithVars(src, vars)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func BenchmarkStringToSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StringToSliceWithVars(`'hello' hello123 "world"`, nil)
	}
}
