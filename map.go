package shellparse

import (
	"fmt"
	"strings"
)

// StringToMap parses key value pairs into map. Handles multiline
// strings and comments.
func StringToMap(src string) (map[string]string, error) {
	return StringToMapWithMap(src, nil)
}

// StringToMapWithMap the same as StringToMap, but additionally
// performs replacement of ${VAR} with provided k/v map.
func StringToMapWithMap(src string, vars map[string]string) (map[string]string, error) {
	words, err := StringToSliceWithMap(src, vars)
	if err != nil {
		return nil, err
	}

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
