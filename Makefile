SHELL := $(shell which bash)
VERBOSE_SHORT_ARG := $(if $(filter $(VERBOSE),true),-v,)
VERBOSE_LONG_ARG := $(if $(filter $(VERBOSE),true),--verbose,)

ROOT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SCRIPTS_DIR := $(ROOT_DIR)/scripts
ENVFILE := $(SCRIPTS_DIR)/env.sh
RELEASE_SUPPORT := $(SCRIPTS_DIR)/release-support.sh

CLUSTER_PROVIDER := kind
OPERATOR_SDK := operator-sdk
IMAGE_LOADER := $(SCRIPTS_DIR)/load-image.sh -p $(CLUSTER_PROVIDER) $(VERBOSE_SHORT_ARG)
ifdef CLUSTER_NAME
	IMAGE_LOADER += -c $(CLUSTER_NAME)
endif

OPERATOR_IMAGE := $(shell . $(ENVFILE) ; echo $${XROOTD_OPERATOR_IMAGE_REPO})
VERSION := $(shell . $(RELEASE_SUPPORT) ; getVersion)

.PHONY: help version bundle olm-generate uninstall code-vet code-fmt code code-gen build-image build dev-install clean

help: ## Display this help
	@echo -e "Usage:\n  make \033[36m<target>\033[0m"
	@awk 'BEGIN {FS = ":.*##"}; \
		/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Application

bundle: olm-generate ## Bundle the operator in OLM format
	@$(SHELL) $(SCRIPTS_DIR)/olm-bundle.sh

uninstall: ## Uninstall the operator
	@echo "....... Uninstalling ......."
	@$(ROOT_DIR)/deploy/operator.sh -u

##@ Development

code-vet: ## Run go vet for this project. More info: https://golang.org/cmd/vet/
	@echo "go vet"
	go vet $$(go list ./... )

code-fmt: ## Run go fmt for this project
	@echo "go fmt"
	go fmt $$(go list ./... )

code: ## Run the default dev commands
	@echo "Running the common required commands for development purposes"
	- make code-fmt
	- make code-vet
	- make code-gen

code-gen: ## Run the operator-sdk commands to generated code (k8s and openapi)
	@echo "Updating the deep copy files with the changes in the API"
	@GOROOT=`pwd` $(OPERATOR_SDK) generate k8s $(VERBOSE_LONG_ARG)
	@echo "Updating the CRD files with the OpenAPI validations"
	$(OPERATOR_SDK) generate crds $(VERBOSE_LONG_ARG)

dev-install: ## Deploy the operator locally
	@echo "....... Installing local build ......."
	@$(ROOT_DIR)/deploy/operator.sh -d $(VERBOSE_SHORT_ARG)

build: build-image ## Build the Operator Image and load it in your cluster
	sed "s|REPLACE_IMAGE|$(OPERATOR_IMAGE):$(VERSION)|g" "$(ROOT_DIR)/deploy/operator.yaml.tpl" > "$(ROOT_DIR)/deploy/operator.yaml"
	@echo "Loading operator image in '$(if $(CLUSTER_NAME),$(CLUSTER_NAME),$(CLUSTER_PROVIDER))' cluster"
	@$(IMAGE_LOADER) $(OPERATOR_IMAGE):$(VERSION)

build-image: ## Build the operator docker image
	$(OPERATOR_SDK) build $(OPERATOR_IMAGE):$(VERSION) $(VERBOSE_LONG_ARG)
	@docker tag $(OPERATOR_IMAGE):$(VERSION) $(OPERATOR_IMAGE):latest

clean: ## Clean build outputs
	rm -r build/_output/

##@ Versioning

version-image: .release ## Shows the current release tag based on the directory content.
	@. $(RELEASE_SUPPORT); getVersion

version: ## Shows the current release version based on version/version.go
	@. $(ENVFILE) ; echo $$XROOTD_OPERATOR_VERSION

olm-generate: ## Generates the required CSV manifests
	@echo "....... Generating CSV ......."
	@$(SHELL) $(SCRIPTS_DIR)/olm-generate-csv.sh

##@ Tests

tests-e2e: ## Run e2e tests
	@echo "....... Running e2e tests ......."
	@$(SCRIPTS_DIR)/run-e2e-tests.sh $(VERBOSE_SHORT_ARG)

tests-unit: ## Run unit tests
	@echo "....... Running unit tests ......."
	@echo "None found"

.release:
	@echo "release=0.0.0" > .release
	@echo "tag=v0.0.0" >> .release
	@echo 'pre_tag_command=sed -i -e "s/^Version = .*/Version = \"@@RELEASE@@\"/" version/version.go' >> .release
	@echo INFO: .release created
	@cat .release
