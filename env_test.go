package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnvFile(t *testing.T) {
	res, err := ParseEnvFile("./_testdata/dotenv.txt")
	exp := map[string]string{"FOO": "bar"}

	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestParseEnvFileWithExtraEnv(t *testing.T) {
	env := map[string]string{"BAR": "baz"}

	res, err := ParseEnvFileWithEnv("./_testdata/dotenv_env.txt", env)
	exp := map[string]string{"FOO": "baz"}

	require.NoError(t, err)
	require.Equal(t, exp, res)
}

func TestParseEnvironmentFile_Error(t *testing.T) {
	m, err := ParseEnvFileWithEnv("./testdata/_absent_file_", nil)

	require.Error(t, err)
	require.Empty(t, m)
}

func TestParseEnvironmentFile_Muted(t *testing.T) {
	m, err := ParseEnvFileWithEnv("-./testdata/_absent_file_", nil)

	require.NoError(t, err)
	require.Empty(t, m)
}
