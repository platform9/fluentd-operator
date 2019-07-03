package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/platform9/fluentd-operator/pkg/options"
)

// GetCfgMapData returns a map of filename==>contents
func GetCfgMapData(subdir string) (map[string]string, error) {
	subDir := fmt.Sprintf("%s/conf/%s", *(options.CfgDir), subdir)
	data := map[string]string{}
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

		data[file.Name()] = string(content)
	}

	return data, nil
}
