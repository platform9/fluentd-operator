/*
Copyright 2019 Platform9 Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
