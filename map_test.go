package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMapString(t *testing.T) {
	exp := map[string]string{"hello1": "world1", "he/llo2": "world2", "hello3": "joe", "hello4": "world4"}
	env := map[string]string{"USER": "joe"}

	src := `'hello1=world1' "he/llo2=world2" hello3=${USER} \
	"hello4=world4"`

	res, err := StringToMapWithEnv(src, env)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestParseMapString_Multiline(t *testing.T) {
	exp := map[string]string{"FOO": "bar", "BAZ": "foo", "HTTP_PROXY": ""}
	src := `# this is a comment
FOO=bar

# hello
BAZ=foo
HTTP_PROXY=
# another comment`
	res, err := StringToMapWithEnv(src, nil)
	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func BenchmarkParseMapString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StringToMapWithEnv(`hello1=world1 hello2=world2`, nil)
	}
}
