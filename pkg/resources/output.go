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
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Output implements the Resource interface for type "output"
type Output struct {
	client     client.Client
	obj        *v1alpha1.Output
	paramCache map[string]string
}

// NewOutput returns a new output resource
func NewOutput(c client.Client, in *v1alpha1.Output) *Output {
	return &Output{
		client:     c,
		obj:        in,
		paramCache: map[string]string{},
	}
}

// Render returns byte array representing fluentd configuration for an output object
func (o *Output) Render() ([]byte, error) {
	validTypes := map[string]bool{
		"stdout":        true,
		"elasticsearch": true,
		"s3":            true,
		"ender":         true,
	}
	outputType := strings.ToLower(o.obj.Spec.Type)

	if _, ok := validTypes[outputType]; !ok {
		// TODO: Build error handling
		return []byte{}, fmt.Errorf("Invalid type: %s", o.obj.Spec.Type)
	}

	params := map[string]string{}
	var err error
	switch outputType {
	case "elasticsearch":
		if params, err = o.getEsParams(); err != nil {
			return []byte{}, err
		}
	case "s3":
		if params, err = o.getS3Params(); err != nil {
			return []byte{}, err
		}
	}

	var ret bytes.Buffer
	fmt.Fprintf(&ret, "<match kube.**>")
	for k, v := range params {
		fmt.Fprintf(&ret, "\n    %s    %s", k, v)
	}
	fmt.Fprintf(&ret, "\n</match>")

	// Always append null match in the end
	fmt.Fprintf(&ret, "\n<match **>")
	fmt.Fprintf(&ret, "\n    @type null")
	fmt.Fprintf(&ret, "\n</match>")
	return ret.Bytes(), nil
}

func (o *Output) getEsParams() (map[string]string, error) {
	indexName := fmt.Sprintf("fluentd-%s", o.obj.Name)
	params := map[string]string{}

	params["@type"] = "elasticsearch"

	var err error

	for _, p := range o.obj.Spec.Params {
		name := strings.ToLower(p.Name)
		v := p.Value
		if len(v) == 0 {
			if v, err = o.getValueFrom(&p.ValueFrom); err != nil {
				return map[string]string{}, err
			}
		}

		params[name] = v
	}

	if _, ok := params["index_name"]; !ok {
		params["index_name"] = indexName
	}

	if v, ok := params["url"]; ok {
		u, err := url.Parse(v)
		if err != nil {
			return map[string]string{}, err
		}
		if u.Port() != "" {
			params["port"] = u.Port()
		}
		if u.Hostname() != "" {
			params["host"] = u.Hostname()
		}
		if u.Scheme != "" {
			params["scheme"] = u.Scheme
		}
		delete(params, "url")
	} else {
		params["host"] = "elasticsearch"
		params["port"] = "9200"
		params["scheme"] = "http"
	}
	return params, nil
}

func (o *Output) getS3Params() (map[string]string, error) {
	var params = make(map[string]string, 1)
	params["@type"] = "s3"

	for _, p := range o.obj.Spec.Params {
		name := strings.ToLower(p.Name)
		v := p.Value
		var err error
		if len(v) == 0 {
			if v, err = o.getValueFrom(&p.ValueFrom); err != nil {
				return map[string]string{}, err
			}
		}

		params[name] = v
	}

	mandatoryParams := []string{"aws_key_id", "aws_sec_key", "s3_bucket", "s3_region"}

	for _, mp := range mandatoryParams {
		if _, ok := params[mp]; !ok {
			return map[string]string{}, fmt.Errorf("Mandatory S3 parameter %s is missing", mp)
		}
	}

	return params, nil
}

func (o *Output) getValueFrom(vf *v1alpha1.ValueFrom) (string, error) {
	key := fmt.Sprintf("%s.%s.%s", vf.Namespace, vf.Name, vf.Key)
	if v, ok := o.paramCache[key]; ok {
		return v, nil
	}

	secret := corev1.Secret{}
	secretName := types.NamespacedName{Name: vf.Name, Namespace: vf.Namespace}

	if err := o.client.Get(context.TODO(), secretName, &secret); err != nil {
		log.Printf("Secret: %s, Error: %s", key, err)
		return "", err
	}

	for k, v := range secret.Data {
		key = fmt.Sprintf("%s.%s.%s", vf.Namespace, vf.Name, k)
		o.paramCache[key] = string(v)
	}

	if v, ok := o.paramCache[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("Key %s was not found in secret %s", vf.Key, vf.Name)
}
