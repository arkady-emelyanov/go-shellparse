package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringToMap(t *testing.T) {
	exp := map[string]string{"hello1": "world1", "he/llo2": "world2"}
	src := `'hello1=world1' "he/llo2=world2"`

	res, err := StringToMap(src)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestStringToMapNoKey(t *testing.T) {
	src := `=hello`

	res, err := StringToMap(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestStringToMapNoVal(t *testing.T) {
	src := `hello`

	res, err := StringToMap(src)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestStringToMapWithVars(t *testing.T) {
	exp := map[string]string{"hello1": "world1", "he/llo2": "world2", "hello3": "joe", "hello4": "world4"}
	vars := map[string]string{"USER": "joe"}

	src := `'hello1=world1' "he/llo2=world2" hello3=${USER} \
	"hello4=world4"`

	res, err := StringToMapWithMap(src, vars)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestStringToMap_Multiline(t *testing.T) {
	exp := map[string]string{"FOO": "bar", "BAZ": "foo", "HTTP_PROXY": ""}
	src := `# this is a comment
FOO=bar

# hello
BAZ=foo
HTTP_PROXY=
# another comment`
	res, err := StringToMapWithMap(src, nil)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func BenchmarkStringToMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = StringToMapWithMap(`hello1=world1 hello2=world2`, nil)
	}
}
