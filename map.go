package shellparse

import (
	"fmt"
	"strings"
)

// StringToMap parses key value pairs into map. Handles multiline
// strings and comments.
func StringToMap(src string) (map[string]string, error) {
	return StringToMapWithEnv(src, nil)
}

// StringToMapWithEnv the same as StringToMap, but additionally
// performs replacement of ${VAR} with provided k/v map.
func StringToMapWithEnv(src string, env map[string]string) (map[string]string, error) {
	words, err := StringToSliceWithEnv(src, env)
	if err != nil {
		return nil, err
	}

	// split by '=' rune
	res := make(map[string]string)
	for i := range words {
		j := strings.IndexRune(words[i], '=')
		if j < 0 {
			return nil, fmt.Errorf("no `=` found in string `%s`", words[i])
		}
		if j == 0 {
			return nil, fmt.Errorf("empty key in string `%s`", words[i])
		}

		key := words[i][:j]
		val := words[i][j+1:]
		res[key] = val
	}

	return res, nil
}
