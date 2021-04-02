package shellparse

import (
	"fmt"
	"os"
	"strings"
)

// StringToSlice parses string into slice.
func StringToSlice(src string) ([]string, error) {
	return StringToSliceWithMap(src, nil)
}

// StringToSliceWithEnv same as StringToSlice, but expand current environment
// variables in string.
func StringToSliceWithEnv(src string) ([]string, error) {
	envs := strings.Join(os.Environ(), " ")
	vars, err := StringToMap(envs)
	if err != nil {
		return nil, err
	}
	return StringToSliceWithMap(src, vars)
}

// StringToSliceWithMap same as StringToSlice, but additionally
// performs replacement of ${VAR} with provided k/v map.
func StringToSliceWithMap(input string, vars map[string]string) ([]string, error) {
	words, err := splitWordsFsm(input)
	if err != nil {
		return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), input)
	}

	for i := range words {
		// `${ENV}` to `value` if vars is not nil
		if vars != nil {
			words[i], err = replaceVarsFsm(words[i], vars)
			if err != nil {
				return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), words[i])
			}
		}

		// 'hello "world\'s"' -> hello "world's"
		words[i], err = unescapeWordsFsm(words[i])
		if err != nil {
			return nil, fmt.Errorf("`%s` in value `%s`", err.Error(), words[i])
		}
	}

	return words, nil
}
