package objects

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

const SuiteName = "K8SObjects"

func TestK8SObjects(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		SuiteName,
		[]Reporter{printer.NewlineReporter{}})
}
