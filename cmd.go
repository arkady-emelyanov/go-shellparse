package shellparse

import (
	"fmt"
)

// ParseCommand parses string into binary and arguments
// and removes unnecessary escape runes.
func ParseCommand(src string) (string, []string, error) {
	return ParseCommandWithEnv(src, nil)
}

// ParseCommandWithEnv same as ParseCommand, but additionally
// perform replacement of ${VAR} with provided k/v map.
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
