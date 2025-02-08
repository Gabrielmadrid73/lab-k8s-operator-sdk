# lab-k8s-operator-sdk
It's a lab using operator-sdk.

## Description

The CRD scaledeployment contains the name of Deployment and number of replicas. The operator will read it and patch the Deployment to contain the quantity of replicas specified in CRD.

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

**Install the CRDs into the cluster:**

```sh
make install
```

**Install the samples into the cluster:**

```sh
kubectl kustomize config/samples/ | kubectl apply -f -
```

**Run the application locally :**

```sh
make install run
```
