package shellparse

import "fmt"

// StringToSlice parses string into <command, args[]> form
// error code indicates that src string is invalid
func StringToSlice(src string) ([]string, error) {
	return StringToSliceWithEnv(src, nil)
}

// StringToSliceWithEnv parses string delimited by space into slice
//
// 1) parse string into []string
// 2) remove unnecessary escape runes
// 3) replace ${VAR} placeholders
//
func StringToSliceWithEnv(input string, env map[string]string) ([]string, error) {
	words, err := splitWordsFsm(input)
	if err != nil {
		return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), input)
	}

	// words postprocessing
	for i := range words {
		replaced, err := replaceVarsFsm(words[i], env) // `%{ENV}` to `env_value`
		if err != nil {
			return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), words[i])
		}

		unescaped, err := unescapeWordsFsm(replaced) // 'hello "world\'s"' -> hello "world's"
		if err != nil {
			return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), replaced)
		}

		words[i] = unescaped
	}

	return words, nil
}
