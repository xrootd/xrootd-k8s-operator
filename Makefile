OPERATOR_SDK := operator-sdk
OPERATOR_IMAGE := "xrootd/xrootd-operator"

ROOT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
RELEASE_SUPPORT := $(ROOT_DIR)/release-support.sh
VERSION := $(shell . $(RELEASE_SUPPORT) ; getVersion)

.PHONY: help version build-k8s build-crds build-image build deploy-operator

help: ### Show this help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: build-k8s build-crds build-image ### Build K8S, CRDs and Operator Image

build-k8s:
	@GOROOT=`pwd` $(OPERATOR_SDK) generate k8s

build-crds:
	$(OPERATOR_SDK) generate crds

build-image:
	$(OPERATOR_SDK) build $(OPERATOR_IMAGE):$(VERSION)
	@docker tag $(OPERATOR_IMAGE):$(VERSION) $(OPERATOR_IMAGE):latest
	sed "s|REPLACE_IMAGE|$(OPERATOR_IMAGE):$(VERSION)|g" "$(ROOT_DIR)/deploy/operator.yaml.tpl" > "$(ROOT_DIR)/deploy/operator.yaml"

deploy-operator: ### Deploy the operator locally
	@sh $(ROOT_DIR)/deploy/operator.sh

version: .release ### Shows the current release tag based on the directory content.
	@. $(RELEASE_SUPPORT); getVersion

.release:
	@echo "release=0.0.0" > .release
	@echo "tag=v0.0.0" >> .release
	@echo 'pre_tag_command=sed -i -e "s/^Version = .*/Version = \"@@RELEASE@@\"/" version/version.go' >> .release
	@echo INFO: .release created
	@cat .release
