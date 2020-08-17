SHELL := $(shell which bash)

VERBOSE_SHORT_ARG := $(if $(filter $(VERBOSE),true),-v,)
VERBOSE_LONG_ARG := $(if $(filter $(VERBOSE),true),--verbose,)

ROOT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SCRIPTS_DIR := $(ROOT_DIR)/hack/scripts
ENVFILE := $(SCRIPTS_DIR)/env.sh
RELEASE_SUPPORT := $(SCRIPTS_DIR)/release-support.sh

CLUSTER_PROVIDER := $(shell . $(ENVFILE) ; echo $${KUBE_CLUSTER_PROVIDER})
OPERATOR_SDK := operator-sdk
IMAGE_LOADER := $(SCRIPTS_DIR)/load-image.sh -p $(CLUSTER_PROVIDER) $(VERBOSE_SHORT_ARG)
ifdef CLUSTER_NAME
	IMAGE_LOADER += -c $(CLUSTER_NAME)
endif

# Current Operator version
VERSION ?= $(shell . $(RELEASE_SUPPORT) ; getVersion)
# Current Bundle Version
BUNDLE_VERSION ?= $(shell . $(ENVFILE) ; echo $${XROOTD_OPERATOR_VERSION})
# Default bundle image tag
BUNDLE_IMG ?= $(shell . $(ENVFILE) ; echo $${XROOTD_OPERATOR_BUNDLE_IMAGE})
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= $(shell . $(ENVFILE) ; echo $${XROOTD_OPERATOR_IMAGE})
# Produce CRDs that work with apiextensions.k8s.io/v1
CRD_OPTIONS ?= "crd:crdVersions=v1"

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

##@ Controller
manager: generate fmt vet ## Build manager binary
	go build -o bin/manager main.go

run: generate fmt vet manifests ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go

install: manifests kustomize ## Install CRDs into a cluster
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

build: docker-build ## Build the Operator Image and load it in your cluster
	@echo "Loading operator image in '$(if $(CLUSTER_NAME),$(CLUSTER_NAME),$(CLUSTER_PROVIDER))' cluster"
	@$(IMAGE_LOADER) ${IMG}:${VERSION}

deploy: manifests kustomize ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}:${VERSION}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: manifests kustomize ## Uninstalls the controller and CRDs in the configured Kubernetes cluster in ~/.kube/config
	$(KUSTOMIZE) build config/default | kubectl delete -f -


##@ Tests
test: generate fmt vet manifests ## Run tests
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/master/hack/setup-envtest.sh
	@{ \
		set -e; \
		source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; \
		fetch_envtest_tools $(ENVTEST_ASSETS_DIR); \
		setup_envtest_env $(ENVTEST_ASSETS_DIR); \
		echo "....... Running tests ......."; \
		go test ./... -coverprofile cover.out; \
	}

test-e2e: ## Run e2e tests
	@echo "....... Running e2e tests ......."
	@$(SCRIPTS_DIR)/run-e2e-tests.sh $(VERBOSE_SHORT_ARG)


##@ Code
fmt: ## Run go fmt against code
	go fmt ./...

vet: ## Run go vet against code
	go vet ./...

generate: controller-gen ## Generate code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

manifests: controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

docker-build: ## Build the docker image
	docker build . -t ${IMG}:${VERSION}
	docker tag ${IMG}:${VERSION} ${IMG}:latest

docker-push: ## Push the docker image
	docker push ${IMG}


##@ OLM
.PHONY: bundle
bundle: manifests ## Generate bundle manifests and metadata, then validate generated files.
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(BUNDLE_VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .


##@ Misc.
help: ## Display this help
	@echo -e "Usage:\n  make \033[36m<target>\033[0m"
	@awk 'BEGIN {FS = ":.*##"}; \
		/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: version
version: ## Shows the current release version based on version/version.go
	@. $(ENVFILE) ; echo $$XROOTD_OPERATOR_VERSION

.PHONY: version-image
version-image: .release ## Shows the current release tag based on the directory content.
	@. $(RELEASE_SUPPORT); getVersion

sample: kustomize ## Install sample manifests
	$(KUSTOMIZE) build manifests/base | kubectl apply -f -


# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# find or download kustomize if necessary
kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

.release:
	@echo "release=0.0.0" > .release
	@echo "tag=v0.0.0" >> .release
	@echo 'pre_tag_command=sed -i -e "s/^Version = .*/Version = \"@@RELEASE@@\"/" version/version.go' >> .release
	@echo INFO: .release created
	@cat .release
