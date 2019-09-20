package v1alpha1

import (
	"reflect"
	"testing"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	"github.com/platform9/fluentd-operator/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestList(t *testing.T) {
	var test = struct {
		description string
		expected    string
		actual      string
		obj         runtime.Object
	}{
		"List output objects", "v1alpha1.OutputList", "", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	result, err := client.LoggingV1alpha1().Outputs().List(metav1.ListOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	test.actual = reflect.TypeOf(*result).String()
	if test.actual != test.expected {
		t.Errorf("Unexpected result. ( Got: %s, want: %s )", test.actual, test.expected)
	}
}

func TestCreate(t *testing.T) {
	var test = struct {
		description string
		expected    string
		actual      string
		obj         runtime.Object
	}{
		"Create output object", "v1alpha1.Output", "", &loggingv1alpha1.Output{},
	}

	input := &loggingv1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample-objstore",
		},
		Spec: loggingv1alpha1.OutputSpec{
			Type:   "elasticsearch",
			Params: []loggingv1alpha1.Param{},
		},
	}

	client := fake.NewSimpleClientset(test.obj)
	result, err := client.LoggingV1alpha1().Outputs().Create(input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	test.actual = reflect.TypeOf(*result).String()
	if test.actual != test.expected {
		t.Errorf("Unexpected result. ( Got: %s, want: %s )", test.actual, test.expected)
	}
}

func TestGet(t *testing.T) {
	var test = struct {
		description string
		input       string
		expected    string
		actual      string
		obj         runtime.Object
	}{
		"List output objects", "sample-objstore", "v1alpha1.Output", "", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	result, err := client.LoggingV1alpha1().Outputs().Get(test.input, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	test.actual = reflect.TypeOf(*result).String()
	if test.actual != test.expected {
		t.Errorf("Unexpected result. ( Got: %s, want: %s )", test.actual, test.expected)
	}
}

func TestUpdate(t *testing.T) {
	var test = struct {
		description string
		expected    string
		actual      string
		obj         runtime.Object
	}{
		"Create output object", "v1alpha1.Output", "", &loggingv1alpha1.Output{},
	}

	input := &loggingv1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample-objstore",
		},
		Spec: loggingv1alpha1.OutputSpec{
			Type:   "elasticsearch",
			Params: []loggingv1alpha1.Param{},
		},
	}

	client := fake.NewSimpleClientset(test.obj)
	result, err := client.LoggingV1alpha1().Outputs().Update(input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	test.actual = reflect.TypeOf(*result).String()
	if test.actual != test.expected {
		t.Errorf("Unexpected result. ( Got: %s, want: %s )", test.actual, test.expected)
	}
}

func TestDelete(t *testing.T) {
	var test = struct {
		description string
		input       string
		expected    string
		actual      string
		obj         runtime.Object
	}{
		"Delete output object", "sample-objstore", "v1alpha1.Output", "", &loggingv1alpha1.Output{},
	}

	client := fake.NewSimpleClientset(test.obj)
	err := client.LoggingV1alpha1().Outputs().Delete(test.input, &metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
