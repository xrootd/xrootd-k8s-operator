module github.com/xrootd/xrootd-k8s-operator

go 1.13

require (
	github.com/RHsyseng/operator-utils v0.0.0-20200304191317-2425bf382482
	github.com/coreos/prometheus-operator v0.41.1 // indirect
	github.com/go-logr/logr v0.2.0
	github.com/msoap/byline v1.1.1
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
	sigs.k8s.io/controller-runtime v0.6.2
)
