package shellparse

import (
	"fmt"
)

// Command parses string into binary and arguments
// and removes unnecessary escape runes.
func Command(src string) (string, []string, error) {
	return CommandWithVars(src, nil)
}

// CommandWithVars same as Command, but additionally
// performs replacement of ${VAR} with provided k/v map.
func CommandWithVars(src string, vars map[string]string) (string, []string, error) {
	parts, err := StringToSliceWithVars(src, vars)
	if err != nil {
		return "", nil, err
	}

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("parse failed for: '%s'", src)
	}

	return parts[0], parts[1:], nil
}
