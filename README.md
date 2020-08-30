# Xrootd Operator [![Xrootd operator CI](https://github.com/xrootd/xrootd-k8s-operator/workflows/Xrootd%20operator%20CI/badge.svg)](https://github.com/xrootd/xrootd-k8s-operator/actions?query=workflow%3A"Xrootd+operator+CI") [![Xrootd operator OLM](https://github.com/xrootd/xrootd-k8s-operator/workflows/Xrootd%20operator%20OLM/badge.svg)](https://github.com/xrootd/xrootd-k8s-operator/actions?query=workflow%3A"Xrootd+operator+OLM")

[![Go Report Card](https://goreportcard.com/badge/github.com/xrootd/xrootd-k8s-operator)](https://goreportcard.com/report/github.com/xrootd/xrootd-k8s-operator) [![codecov](https://codecov.io/gh/xrootd/xrootd-k8s-operator/branch/master/graph/badge.svg)](https://codecov.io/gh/xrootd/xrootd-k8s-operator)

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

- Run `make manager` to locally build operator binary and `make run` to run it against the configured Kubernetes cluster.
- Run `make build` to build operator image from scratch and loads it in the k8s cluster.
- The build command can be configured with the cluster's name and provider to target where the built operator image will be loaded. Set the following environment variables:
  - `CLUSTER_PROVIDER=(kind/k3s/minishift)`
  - `CLUSTER_NAME=<cluster name>`

### Install operator

- Run `make deploy` to deploy the operator image in the cluster, along with applying the required roles, service accounts etc.
- To uninstall the CRDs, run `make uninstall`. To cleanup everything, including the operator deployment, run `make undeploy`.

### Bundle

Xrootd Operator is integrated with OLM and configured to use [Bundle](https://sdk.operatorframework.io/docs/olm-integration/quickstart-bundle/) format.

- To generate OLM CSV manifests and bundle metadata, run `make bundle`.
- To build the operator bundle image, run `make bundle-build`.

### Testing

- **Unit Tests:** Run the unit tests with `make test`.
- **Integration Tests:** Run the suite of e2e tests with `make test-e2e`.

### OpenShift Cluster

- For local development, it's recommended to use [CodeReady Containers](https://code-ready.github.io/crc/) since it supports Openshift v4+. Minishift is a suitable alternative, however it only supports till OpenShift v3.
- To test operator via scripted approach, `make deploy` works.
- To test operator using OLM, follow [testing guide](https://github.com/operator-framework/community-operators/blob/master/docs/testing-operators.md#testing-operator-deployment-on-kubernetes) for deployment using custom images.
  > TODO: [Testing bundles](https://sdk.operatorframework.io/docs/olm-integration/quickstart-bundle/#testing-bundles) is still not officially supported.

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
