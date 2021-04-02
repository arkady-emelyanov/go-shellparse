package shellparse

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	cmds := []string{
		`docker run -it --rm hello-world`,
	}
	expt := [][]string{
		{`docker`, `run`, `-it`, `--rm`, `hello-world`},
	}

	for i, cmd := range cmds {
		bin, args, err := Command(cmd)
		exp := expt[i]

		require.NoError(t, err)
		require.Equal(t, exp[0], bin)
		require.Equal(t, exp[1:], args)
	}
}

func TestCommandMultiline(t *testing.T) {
	cmd := `docker run -it \
				--rm \
				-v /tmp:/tmp:rw`
	bin, args, err := Command(cmd)

	exp := []string{"docker", "run", "-it", "--rm", "-v", "/tmp:/tmp:rw"}
	require.NoError(t, err)
	require.Equal(t, exp[0], bin)
	require.Equal(t, exp[1:], args)
}

func TestCommandEmpty(t *testing.T) {
	bin, args, err := Command("")

	require.Error(t, err)
	require.Equal(t, "", bin)
	require.Nil(t, args)
}

func TestCommandWithEnv(t *testing.T) {

	_ = os.Setenv("SLEEP", "1")
	bin, args, err := CommandWithEnv(`bash -c 'sleep ${SLEEP}'`)
	_ = os.Unsetenv("SLEEP")

	exp := []string{"bash", "-c", "sleep 1"}
	require.NoError(t, err)
	require.Equal(t, exp[0], bin)
	require.Equal(t, exp[1:], args)
}

func TestCommandWithVars(t *testing.T) {
	cmds := []string{
		`touch ${HOME}/.hushlogin`,
		`touch \${HOME}/.hushlogin`,
	}
	expt := [][]string{
		{`touch`, `/home/john/.hushlogin`},
		{`touch`, `${HOME}/.hushlogin`},
	}

	vars := map[string]string{"HOME": "/home/john"}
	for i, cmd := range cmds {
		bin, args, err := CommandWithMap(cmd, vars)
		exp := expt[i]

		require.NoError(t, err)
		require.Equal(t, exp[0], bin)
		require.Equal(t, exp[1:], args)
	}
}

func TestCommandWithQuotes(t *testing.T) {
	multilineCmd := `bash -c '
echo "ok" &&
echo "ok" &&
echo "ok"
'`
	multilineExp := []string{`bash`, `-c`, `
echo "ok" &&
echo "ok" &&
echo "ok"
`}

	cmds := []string{
		`touch ~/.config/config.txt`,
		`echo "it's ok"`,
		`it\"s ok`,
		`it\'s ok`,
		`\\'it was escaped\\'`,
		`bash -c 'echo "look\'s like i\'m failing" && sleep 3 && exit 1'`,
		multilineCmd,
	}

	expt := [][]string{
		{`touch`, `~/.config/config.txt`},
		{`echo`, `it's ok`},
		{`it"s`, `ok`},
		{`it's`, `ok`},
		{`it was escaped`},
		{`bash`, `-c`, `echo "look's like i'm failing" && sleep 3 && exit 1`},
		multilineExp,
	}

	for i, cmd := range cmds {
		bin, res, err := Command(cmd)
		exp := expt[i]

		require.NoError(t, err)
		require.Equal(t, exp[0], bin)
		require.Equal(t, exp[1:], res)
	}
}

func BenchmarkParseCommand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Command(`bash -c 'echo "look\'s like i\'m failing" && sleep 3 && exit 1'`)
	}
}
