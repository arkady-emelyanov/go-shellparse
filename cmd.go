package shellparse

import (
	"fmt"
)

func ParseCommand(src string) (string, []string, error) {
	return ParseCommandWithEnv(src, nil)
}

//
// command and arguments parser
//
// 1) parse string into binary and arguments
// 2) remove unnecessary escape runes
// 3) replace ${VAR} placeholders with appropriate key from env map
//
func ParseCommandWithEnv(src string, env map[string]string) (string, []string, error) {
	parts, err := StringToSliceWithEnv(src, env)
	if err != nil {
		return "", nil, err
	}

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("parse failed for: '%s'", src)
	}

	return parts[0], parts[1:], nil
}
