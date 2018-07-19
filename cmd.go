package shellparse

import (
	"fmt"
)

// Command parses string into binary and arguments
// and removes unnecessary escape runes.
func Command(src string) (string, []string, error) {
	return CommandWithEnv(src, nil)
}

// CommandWithEnv same as Command, but additionally
// performs replacement of ${VAR} with provided k/v map.
func CommandWithEnv(src string, env map[string]string) (string, []string, error) {
	parts, err := StringToSliceWithEnv(src, env)
	if err != nil {
		return "", nil, err
	}

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("parse failed for: '%s'", src)
	}

	return parts[0], parts[1:], nil
}
