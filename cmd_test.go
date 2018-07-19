package shellparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommand(t *testing.T) {
	cmds := []string{
		`docker run -it --rm hello-world`,
	}
	expt := [][]string{
		{`docker`, `run`, `-it`, `--rm`, `hello-world`},
	}

	for i, cmd := range cmds {
		bin, res, err := ParseCommand(cmd)
		exp := expt[i]

		require.NoError(t, err)
		require.Equal(t, exp[0], bin)
		require.Equal(t, exp[1:], res)
	}
}

func TestParseCommandWithEnv(t *testing.T) {
	cmds := []string{
		`touch ${HOME}/.hushlogin`,
		`touch \${HOME}/.hushlogin`,
	}
	expt := [][]string{
		{`touch`, `/home/john/.hushlogin`},
		{`touch`, `${HOME}/.hushlogin`},
	}

	env := map[string]string{"HOME": "/home/john"}
	for i, cmd := range cmds {
		bin, res, err := ParseCommandWithEnv(cmd, env)
		exp := expt[i]

		require.NoError(t, err)
		require.Equal(t, exp[0], bin)
		require.Equal(t, exp[1:], res)
	}
}

func TestParseQuotes(t *testing.T) {
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
		bin, res, err := ParseCommand(cmd)
		exp := expt[i]

		require.NoError(t, err)
		require.Equal(t, exp[0], bin)
		require.Equal(t, exp[1:], res)
	}
}

func BenchmarkParseCommand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ParseCommand(`bash -c 'echo "look\'s like i\'m failing" && sleep 3 && exit 1'`)
	}
}
