package controller

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/xrootd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, xrootd.Add)
}
