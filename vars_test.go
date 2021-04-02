package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVarsFile(t *testing.T) {
	res, err := ParseVarsFile("./_testdata/dotenv.txt")
	exp := map[string]string{"FOO": "bar"}

	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestParseVarsFileIncorrect(t *testing.T) {
	res, err := ParseVarsFile("./_testdata/dotenv_incorrect.txt")

	require.Error(t, err)
	require.Nil(t, res)
}

func TestParseVarsFileWithExtraVars(t *testing.T) {
	vars := map[string]string{"FOO": "FOO"}

	res, err := ParseVarsFileWithMap("./_testdata/dotenv_with_vars.txt", vars)
	exp := map[string]string{"FOO": "bar"}

	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestParseVarsFile_Error(t *testing.T) {
	m, err := ParseVarsFileWithMap("./testdata/_absent_file_", nil)

	require.Error(t, err)
	require.Empty(t, m)
}

func TestParseVarsFile_Muted(t *testing.T) {
	m, err := ParseVarsFileWithMap("-./testdata/_absent_file_", nil)

	require.NoError(t, err)
	require.Empty(t, m)
}
