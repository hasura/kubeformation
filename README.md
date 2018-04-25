# kubeformation

Create declarative specifications for your managed Kubernetes cloud vendor (GKE, AKS).

[![GoDoc](https://godoc.org/github.com/hasura/kubeformation?status.svg)](https://godoc.org/github.com/hasura/kubeformation) 

## Motivation

With Kubernetes, it becomes possible to start making everything about your application declarative. As cloud vendors start providing managed Kubernetes services, provisioning a Kubernetes cluster via the vendor’s API becomes declarative as well.

Kubeformation is a simple web UI and CLI that helps you create “Google Deployment manager” or “Azure Resoure Manager” templates which are a _little_ painful to create by hand.

Once you have this file, you can run your cloud vendor CLI on it to provision your cluster. You can edit this file to add vendor specific configuration too.

## Usage

Head to https://kubeformation.sh, build a cluster, download the templates

**or** 

Write cluster.yaml yourself, use the `kubeformation` CLI

## Example

Let's look at a spec that defines a Kubernetes cluster: `cluster.yaml`

```yaml
version: v1
name: cluster-name
provider: gke
k8sVersion: "1.9"
nodePools:
- name: db-pool
  type: n1-standard-1
  size: 1
  labels:
    app: postgres
- name: backend-pool
  type: n1-standard-2
  size: 2
  labels:
    app: backend
volumes:
- name: postgres
  size: 10
```

`kubeformation` can read this file and generate [Google Cloud Deployment
Manager](https://cloud.google.com/deployment-manager/) template, which can then
be used with `gcloud` command to create the GKE cluster. This is a declarative
template that can be used to further do create or modify the cluster.

`kubeformation` is exclusively meant for managed Kubernetes providers. The
following providers are supported:

1. Google Kubernetes Engine (GKE)
2. Azure Container Service (AKS)

Read complete docs [here](docs/README.md).
