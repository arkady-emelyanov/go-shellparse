package shellparse

import "fmt"

// StringToSlice parses string into slice.
func StringToSlice(src string) ([]string, error) {
	return StringToSliceWithEnv(src, nil)
}

// StringToSliceWithEnv same as StringToSlice, but additionally
// performs replacement of ${VAR} with provided k/v map.
func StringToSliceWithEnv(input string, env map[string]string) ([]string, error) {
	words, err := splitWordsFsm(input)
	if err != nil {
		return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), input)
	}

	for i := range words {
		// `%{ENV}` to `value`
		replaced, err := replaceVarsFsm(words[i], env)
		if err != nil {
			return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), words[i])
		}

		// 'hello "world\'s"' -> hello "world's"
		unescaped, err := unescapeWordsFsm(replaced)
		if err != nil {
			return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), replaced)
		}

		words[i] = unescaped
	}

	return words, nil
}
