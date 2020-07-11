module github.com/shivanshs9/xrootd-operator

go 1.13

require (
	github.com/RHsyseng/operator-utils v0.0.0-20200619180557-7c49e58877d7
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/redhat-cop/operator-utils v0.3.2
	github.com/shivanshs9/ty v1.1.0
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29 // indirect
	sigs.k8s.io/controller-runtime v0.6.1
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
