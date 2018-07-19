package shellparse

import (
	"fmt"
	"strings"
)

func StringToMap(src string) (map[string]string, error) {
	return StringToMapWithEnv(src, nil)
}

//
// map parser
//
// 1) parse string into map[string]string
// 2) remove unnecessary escape runes
// 3) replace ${VAR} placeholders
//
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
