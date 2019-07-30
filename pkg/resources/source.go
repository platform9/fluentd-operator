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

package resources

import (
	"bytes"
	"fmt"

	"github.com/platform9/fluentd-operator/pkg/options"
)

// Source represents implementation of fluentd source configuration.
type Source struct {
	port int
}

// NewSource returns a new source object
func NewSource() *Source {
	return &Source{
		port: *(options.ForwardPort),
	}
}

// Render returns byte array representing fluentd configuration of a source
func (s *Source) Render() ([]byte, error) {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "<source>")
	fmt.Fprintf(&ret, "\n    @type forward")
	fmt.Fprintf(&ret, "\n    port %d", s.port)
	fmt.Fprintf(&ret, "\n    bind 0.0.0.0")
	fmt.Fprintf(&ret, "\n</source>")

	return ret.Bytes(), nil
}
