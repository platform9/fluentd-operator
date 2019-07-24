package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/platform9/fluentd-operator/pkg/options"
)

// GetCfgMapData returns a map of filename==>contents
func GetCfgMapData(subdir string) (map[string][]byte, error) {
	subDir := fmt.Sprintf("%s/%s", *(options.CfgDir), subdir)
	data := map[string][]byte{}
	files, err := ioutil.ReadDir(subDir)
	if err != nil {
		return data, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", subDir, file.Name()))
		if err != nil {
			return data, err
		}

		data[file.Name()] = content
	}

	return data, nil
}
