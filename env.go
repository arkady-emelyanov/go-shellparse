package shellparse

import (
	"strings"
	"io/ioutil"
)

func ParseEnvFile(file string) (map[string]string, error) {
	return ParseEnvFileWithEnv(file, nil)
}

// ParseEnvFileWithEnv parse file provided as path into map[string]string
// if file path is prepended with '-' char read file errors will be muted
func ParseEnvFileWithEnv(file string, extraEnv map[string]string) (map[string]string, error) {
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
		tmp, err := StringToMapWithEnv(string(b), extraEnv)
		if err != nil {
			goto Error
		}

		for k, v := range tmp {
			res[k] = v
		}
		return res, nil
	}

Error:
	if muteError == false {
		return nil, err
	} else {
		return res, nil
	}
}

