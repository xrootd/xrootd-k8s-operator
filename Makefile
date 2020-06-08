OPERATOR_SDK := operator-sdk
OPERATOR_IMAGE := "xrootd/xrootd-operator"

.PHONY: build-k8s build-crds build-image build

build: build-k8s build-crds build-image

build-k8s:
	@GOROOT=`pwd` $(OPERATOR_SDK) generate k8s

build-crds:
	$(OPERATOR_SDK) generate crds

build-image:
	$(OPERATOR_SDK) build $(OPERATOR_IMAGE)
