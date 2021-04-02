package shellparse

import (
	"io/ioutil"
	"os"
	"strings"
)

// ParseVarsFile is helper for parsing dotenv compatible files.
// If file path is prepended with '-' char, file read error will not be raised
func ParseVarsFile(file string) (map[string]string, error) {
	return ParseVarsFileWithMap(file, nil)
}

func ParseVarsFileWithEnv(file string) (map[string]string, error) {
	envs := strings.Join(os.Environ(), " ")
	vars, err := StringToMap(envs)
	if err != nil {
		return nil, err
	}
	return ParseVarsFileWithMap(file, vars)
}

// ParseVarsFileWithMap same as ParseVarsFile, but additionally
// performs replacement of ${VAR} with provided k/v map.
// If file path is prepended with '-' char, file read error will not be raised
func ParseVarsFileWithMap(file string, extraEnv map[string]string) (map[string]string, error) {
	var err error

	res := make(map[string]string)
	muteError := false

	if strings.HasPrefix(file, "-") {
		file = file[1:]
		muteError = true
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		goto Error
	}

	if len(b) > 0 {
		var tmp map[string]string

		tmp, err = StringToMapWithMap(string(b), extraEnv)
		if err != nil {
			goto Error
		}

		for k, v := range tmp {
			res[k] = v
		}
		return res, nil
	}

Error:
	if muteError {
		return res, nil
	}

	return nil, err
}
