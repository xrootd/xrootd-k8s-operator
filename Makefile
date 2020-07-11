ROOT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SCRIPTS_DIR := $(ROOT_DIR)/scripts
ENVFILE := $(SCRIPTS_DIR)/env.sh
RELEASE_SUPPORT := $(SCRIPTS_DIR)/release-support.sh

CLUSTER_PROVIDER := kind
OPERATOR_SDK := operator-sdk
IMAGE_LOADER := $(SCRIPTS_DIR)/load-image.sh -p $(CLUSTER_PROVIDER)
ifdef CLUSTER_NAME
	IMAGE_LOADER += -c $(CLUSTER_NAME)
endif

OPERATOR_IMAGE := $(shell . $(ENVFILE) ; echo $${XROOTD_OPERATOR_IMAGE_REPO})
VERSION := $(shell . $(RELEASE_SUPPORT) ; getVersion)

.PHONY: help version build-k8s build-crds build-image build deploy-operator bundle

help: ### Show this help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: build-k8s build-crds build-image ### Build K8S, CRDs and Operator Image

build-k8s:
	@GOROOT=`pwd` $(OPERATOR_SDK) generate k8s

build-crds:
	$(OPERATOR_SDK) generate crds

build-image: ### Built the image and load it in your cluster
	$(OPERATOR_SDK) build $(OPERATOR_IMAGE):$(VERSION)
	@docker tag $(OPERATOR_IMAGE):$(VERSION) $(OPERATOR_IMAGE):latest
	sed "s|REPLACE_IMAGE|$(OPERATOR_IMAGE):$(VERSION)|g" "$(ROOT_DIR)/deploy/operator.yaml.tpl" > "$(ROOT_DIR)/deploy/operator.yaml"
	@echo "Loading operator image in '$(if $(CLUSTER_NAME),$(CLUSTER_NAME),$(CLUSTER_PROVIDER))' cluster"
	@$(IMAGE_LOADER) $(OPERATOR_IMAGE):$(VERSION)

deploy-operator: ### Deploy the operator locally
	@$(ROOT_DIR)/deploy/operator.sh

bundle: ### Bundle the operator in OLM format
	@$(SCRIPTS_DIR)/olm-bundle.sh

version: .release ### Shows the current release tag based on the directory content.
	@. $(RELEASE_SUPPORT); getVersion

.release:
	@echo "release=0.0.0" > .release
	@echo "tag=v0.0.0" >> .release
	@echo 'pre_tag_command=sed -i -e "s/^Version = .*/Version = \"@@RELEASE@@\"/" version/version.go' >> .release
	@echo INFO: .release created
	@cat .release
