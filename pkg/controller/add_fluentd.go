package controller

import (
	"github.com/platform9/fluentd-operator/pkg/controller/fluentd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, fluentd.Add)
}
