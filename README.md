# Xrootd Operator [![Xrootd operator CI](https://github.com/xrootd/xrootd-k8s-operator/workflows/Xrootd%20operator%20CI/badge.svg)](https://github.com/xrootd/xrootd-k8s-operator/actions?query=workflow%3A"Xrootd+operator+CI") [![Xrootd operator OLM](https://github.com/xrootd/xrootd-k8s-operator/workflows/Xrootd%20operator%20OLM/badge.svg)](https://github.com/xrootd/xrootd-k8s-operator/actions?query=workflow%3A"Xrootd+operator+OLM")

A Kubernetes operator to deploy [Xrootd](https://github.com/xrootd/xrootd) at scale, in order to ease and fully automate deployment and management of XRootD clusters.

### Installation

#### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Access to a Kubernetes cluster:
  - For production, use bare-metal clusters or public cloud platforms.
  - For development, use local K8S Cluster
    - [Kind](https://kind.sigs.k8s.io/) - use a [simple script](https://github.com/k8s-school/kind-travis-ci/blob/master/k8s-create.sh)
    - Or, [K3S](https://k3s.io/)
- [Configure](https://success.docker.com/article/how-to-use-kubectl-to-manage-multiple-kubernetes-clusters) kubectl to use relevant K8S Cluster

#### Using OLM

- TODO

### Development

#### Prerequisites

- Same prerequisites for [Installation](#Installation)
- [Go](https://golang.org/doc/install) v1.13+
- [operator-sdk](https://sdk.operatorframework.io/docs/install-operator-sdk/)

#### Build operator

- Run `make code` to run code format and generate code via operator-sdk.
- Run `make build` to build operator image from scratch and loads it in the k8s cluster. Most of the time, this is sufficient to run if there has been no change in CRs ([pkg/apis](pkg/apis))
- The build command can be configured with the cluster's name and provider to target where the built operator image will be loaded. Set the following environment variables:
  - `CLUSTER_PROVIDER=(kind/k3s)`
  - `CLUSTER_NAME=<cluster name>`

#### Install operator

- Run `make dev-install` to deploy the operator image in the cluster, along with applying the required roles, service accounts etc.

### Usage

- Make sure the xrootd-operator is up and runnning in your K8S cluster (otherwise follow [Installation](#Installation)/[Development](#Development) steps):
  - To check the status, run `kubectl describe pod -l name=xrootd-operator`
- Example manifests to deploy Xrootd instance are at [manifests](manifests) folder.
- To apply any manifest, simply use `kubectl apply`:
  - For example, to apply base sample manifest, run `kubectl apply -k manifests/base`

### Bundle

Xrootd Operator is integrated with OLM.

- To generate OLM CSV manifest, run `./scripts/olm-generate-csv.sh`.
- To build the operator bundle image, run `make bundle`.

### Testing

#### Unit tests

- Run the unit tests with `make tests-unit`.

#### Integration tests

- Run the suite of e2e tests with `make tests-e2e`.

### Troubleshooting

- Check operator logs: `kubectl logs -l name=xrootd-operator`
- [Create issue](https://github.com/xrootd/xrootd-k8s-operator/issues/new/choose) and if needed, provide operator logs too.

### Useful Links

- [Xrootd](https://xrootd.slac.stanford.edu/index.html)
- [Xrootd Config Docs](https://xrootd.slac.stanford.edu/doc/dev50/xrd_config.htm)
- [Xrootd Example Dockerfile](https://github.com/lnielsen/xrootd-docker/blob/master/Dockerfile)
