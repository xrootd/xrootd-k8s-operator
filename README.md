# Xrootd Operator [![Xrootd operator CI](https://github.com/xrootd/xrootd-k8s-operator/workflows/Xrootd%20operator%20CI/badge.svg)](https://github.com/xrootd/xrootd-k8s-operator/actions?query=workflow%3A"Xrootd+operator+CI") [![Xrootd operator OLM](https://github.com/xrootd/xrootd-k8s-operator/workflows/Xrootd%20operator%20OLM/badge.svg)](https://github.com/xrootd/xrootd-k8s-operator/actions?query=workflow%3A"Xrootd+operator+OLM")

A Kubernetes operator to deploy [Xrootd](https://github.com/xrootd/xrootd) at scale, in order to ease and fully automate deployment and management of XRootD clusters.

## Installation

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Access to a Kubernetes cluster:
  - For production, use bare-metal clusters or public cloud platforms.
  - For development, use local K8S Cluster
    - [Kind](https://kind.sigs.k8s.io/) - use a [simple script](https://github.com/k8s-school/kind-travis-ci/blob/master/k8s-create.sh)
    - Or, [K3S](https://k3s.io/)
- [Configure](https://success.docker.com/article/how-to-use-kubectl-to-manage-multiple-kubernetes-clusters) kubectl to use relevant K8S Cluster

### Using OLM

- TODO

---

## Development

### Prerequisites

- Same prerequisites for [Installation](#Installation)
- [Go](https://golang.org/doc/install) v1.13+
- [operator-sdk](https://sdk.operatorframework.io/docs/install-operator-sdk/)

### Build operator

- Run `make code` to run code format and generate code via operator-sdk.
- Run `make build` to build operator image from scratch and loads it in the k8s cluster. Most of the time, this is sufficient to run if there has been no change in CRs ([pkg/apis](pkg/apis))
- The build command can be configured with the cluster's name and provider to target where the built operator image will be loaded. Set the following environment variables:
  - `CLUSTER_PROVIDER=(kind/k3s/minishift)`
  - `CLUSTER_NAME=<cluster name>`

### Install operator

- Run `make dev-install` to deploy the operator image in the cluster, along with applying the required roles, service accounts etc.

### Bundle

Xrootd Operator is integrated with OLM.

- Install [operator-courier](https://pypi.org/project/operator-courier) python library for linting of Operator Metadata and to push the application information to quay.io appregistry.
- To generate OLM CSV manifest, run `make olm-generate`.
- To push the operator bundle information to quay.io, run `make push-bundle`.

### Testing

- **Unit Tests:** Run the unit tests with `make tests-unit`.
- **Integration Tests:** Run the suite of e2e tests with `make tests-e2e`.

### OpenShift Cluster

- For local development, it's recommended to use [CodeReady Containers](https://code-ready.github.io/crc/) since it supports Openshift v4+. Minishift is a suitable alternative, however it only supports till OpenShift v3.
- To test operator via scripted approach, `make dev-install` works.
- To test operator using OLM, follow [testing guide](https://github.com/operator-framework/community-operators/blob/master/docs/testing-operators.md#testing-operator-deployment-on-kubernetes) for deployment using custom images.
- Simply run `make push-bundle` to generate the CSV manifests and push to quay.io appregistry.

> **NOTE:**
> Minishift uses Kubernetes v1.11.x, so it only supports till OLM v0.14.x (because later OLM versions uses apiextensions.k8s.io/v1 for CRD manifests)

---

## Usage

- Make sure the xrootd-operator is up and runnning in your K8S cluster (otherwise follow [Installation](#Installation)/[Development](#Development) steps):
  - To check the status, run `kubectl describe pod -l name=xrootd-operator`
- Example manifests to deploy Xrootd instance are at [manifests](manifests) folder.
- To apply any manifest, simply use `kubectl apply`:
  - For example, to apply base sample manifest, run `kubectl apply -k manifests/base`

---

## Troubleshooting

- Check operator logs: `kubectl logs -l name=xrootd-operator`
- [Create issue](https://github.com/xrootd/xrootd-k8s-operator/issues/new/choose) and if needed, provide operator logs too.

## Useful Links

- [Xrootd](https://xrootd.slac.stanford.edu/index.html)
- [Xrootd Config Docs](https://xrootd.slac.stanford.edu/doc/dev50/xrd_config.htm)
- [Xrootd Example Dockerfile](https://github.com/lnielsen/xrootd-docker/blob/master/Dockerfile)
