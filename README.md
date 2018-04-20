# kubeformation

Bootstrap cloud vendor specific declarative templates from a simple spec for a Kubernetes cluster.

[![GoDoc](https://godoc.org/github.com/hasura/kubeformation?status.svg)](https://godoc.org/github.com/hasura/kubeformation) 

## What is kubeformation?

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

Amazon Elastic Container Service for Kubernetes (EKS) is a work in progress.
