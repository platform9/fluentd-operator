package v1alpha1

import (
	"reflect"
	"testing"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	"github.com/platform9/fluentd-operator/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var input = &loggingv1alpha1.Output{
	ObjectMeta: metav1.ObjectMeta{
		Name: "sample-objstore",
	},
	Spec: loggingv1alpha1.OutputSpec{
		Type:   "elasticsearch",
		Params: []loggingv1alpha1.Param{},
	},
}

func TestCreate(t *testing.T) {
	var test = struct {
		description string
		expected    string
		obj         runtime.Object
	}{
		"Create Output Custom Resource", "v1alpha1.Output", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	result, err := client.LoggingV1alpha1().Outputs().Create(input)
	if err != nil {
		t.Errorf("Error creating Output object. Error is: %s", err)
	}

	// Checking type of the result
	t.Run("check result type", func(t *testing.T) {
		if test.expected != reflect.TypeOf(*result).String() {
			t.Errorf("Got Unexpected result. Result: %v", result)
		}
	})
}

func TestList(t *testing.T) {
	var test = struct {
		description string
		expected    string
		obj         runtime.Object
	}{
		"List Output Custom Resource Objects", "v1alpha1.OutputList", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	object, err := client.LoggingV1alpha1().Outputs().Create(input)
	if err != nil {
		t.Errorf("Error creating Output object. Error is: %s", err)
	}

	result, err := client.LoggingV1alpha1().Outputs().List(metav1.ListOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Checking type of the result
	t.Run("check result type", func(t *testing.T) {
		if test.expected != reflect.TypeOf(*result).String() {
			t.Errorf("Got Unexpected result. Result: %v ", result)
		}
	})

	// Checking if object is present in the OutputList
	t.Run("check output object", func(t *testing.T) {
		status := false
		for _, element := range result.Items {
			if object.ObjectMeta.Name == element.ObjectMeta.Name {
				status = true
				break
			}
		}
		if status == false {
			t.Errorf("Unexpected result. Didn't get object %s in the Output List", object.ObjectMeta.Name)
		}
	})
}

func TestGet(t *testing.T) {
	var test = struct {
		description string
		name        string
		expected    string
		obj         runtime.Object
	}{
		"Get Output Custom Resource object", "sample-objstore", "v1alpha1.Output", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	object, err := client.LoggingV1alpha1().Outputs().Create(input)
	if err != nil {
		t.Errorf("Error creating Output object. Error is: %s", err)
	}

	result, err := client.LoggingV1alpha1().Outputs().Get(test.name, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Checking type of the result
	t.Run("check result type", func(t *testing.T) {
		if test.expected != reflect.TypeOf(*result).String() {
			t.Errorf("Got Unexpected result. Result: %v", result)
		}
	})

	// Checking if object is returned in the result
	t.Run("check output object", func(t *testing.T) {
		if result == object {
			t.Errorf("Unexpected result. Didn't get object %s in the Output List", object.ObjectMeta.Name)
		}
	})
}

func TestUpdate(t *testing.T) {
	var test = struct {
		description string
		expected    string
		obj         runtime.Object
	}{
		"Update Output Custom Resource Object", "v1alpha1.Output", &loggingv1alpha1.Output{},
	}

	newInput := &loggingv1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample-objstore",
		},
		Spec: loggingv1alpha1.OutputSpec{
			Type: "elasticsearch",
			Params: []loggingv1alpha1.Param{
				{
					Name:  "url",
					Value: "http://elasticsearch.default.svc.cluster.local:9200 ",
				},
			},
		},
	}

	client := fake.NewSimpleClientset(test.obj)
	object, err := client.LoggingV1alpha1().Outputs().Create(input)
	if err != nil {
		t.Errorf("Error creating Output object. Error is: %s", err)
	}

	result, err := client.LoggingV1alpha1().Outputs().Update(newInput)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Checking type of result
	t.Run("check result type", func(t *testing.T) {
		if test.expected != reflect.TypeOf(*result).String() {
			t.Errorf("Got Unexpected result. Result: %v", result)
		}
	})

	// Checking if object is returned in the result
	t.Run("check output object", func(t *testing.T) {
		if result == object {
			t.Errorf("Unexpected result. Didn't get object %s in the Output List", object.ObjectMeta.Name)
		}
	})
}

func TestDelete(t *testing.T) {
	var test = struct {
		description string
		name        string
		expected    string
		obj         runtime.Object
	}{
		"Delete Output Custom Resource Object", "sample-objstore", "v1alpha1.Output", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	_, err := client.LoggingV1alpha1().Outputs().Create(input)
	if err != nil {
		t.Errorf("Error creating Output object. Error is: %s", err)
	}

	err = client.LoggingV1alpha1().Outputs().Delete(test.name, &metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
