package shellparse

import (
	"fmt"
	"os"
	"strings"
)

// Command parses string into binary and arguments
// and removes unnecessary escape runes.
func Command(src string) (string, []string, error) {
	return CommandWithMap(src, nil)
}

func CommandWithEnv(src string) (string, []string, error) {
	envs := strings.Join(os.Environ(), " ")
	vars, err := StringToMap(envs)
	if err != nil {
		return "", nil, err
	}
	return CommandWithMap(src, vars)
}

// CommandWithMap same as Command, but additionally
// performs replacement of ${VAR} with provided k/v map.
func CommandWithMap(src string, vars map[string]string) (string, []string, error) {
	parts, err := StringToSliceWithMap(src, vars)
	if err != nil {
		return "", nil, err
	}

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("parse failed for: '%s'", src)
	}

	return parts[0], parts[1:], nil
}
