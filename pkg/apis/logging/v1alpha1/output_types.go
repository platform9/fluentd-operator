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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OutputSpec defines the desired state of Output
type OutputSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Type   string  `json:"type"`
	Params []Param `json:"params,omitempty"`
}

// Param defines a parameter to be passed along with output, such as credentials
type Param struct {
	Name      string    `json:"name"`
	ValueFrom ValueFrom `json:"valueFrom,omitempty"`
	Value     string    `json:"value"`
}

// ValueFrom defines a reference to credentials specified in a kubernetes secret
type ValueFrom struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}

// OutputStatus defines the observed state of Output
type OutputStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Output is the Schema for the outputs API
// +k8s:openapi-gen=true
type Output struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OutputSpec   `json:"spec,omitempty"`
	Status OutputStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OutputList contains a list of Output
type OutputList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Output `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Output{}, &OutputList{})
}
