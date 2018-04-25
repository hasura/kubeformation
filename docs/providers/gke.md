# Google Kubernetes Engine (GKE)

For GKE clusters, Kubeformation creates [Google Cloud Deployment
Manager](https://cloud.google.com/deployment-manager/) templates (**GCDM**).
These templates can be used with
[`gcloud`](https://cloud.google.com/sdk/gcloud/) command line tool to create the
Kubernetes cluster. Further edits to the templates can also be applied using
`gcloud` cli. If volumes are present, a corresponding Kubernetes Persistent
Volume and Claim objects are also created, along with underlying [Persistent
Disk](https://cloud.google.com/persistent-disk/). 

The provider string is `gke`.

## Pre-requisites

- [`gcloud`](https://cloud.google.com/sdk/gcloud/) CLI
- [`kubeformation`](../cli.md) CLI (optional)
- A project on Google Cloud (check [this
  link](https://cloud.google.com/resource-manager/docs/creating-managing-projects)
  to create one)
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) CLI

## Step 0 - Write `cluster.yaml`

PS: This step can be skipped if [kubeformation.sh](https://kubeformation.sh) is
used. Jump right ahead to [Step 1](#step-1---generate-gcdm-template).

Here's an example `cluster.yaml`:

```yaml
version: v1
name: my-cluster
provider: gke
k8sVersion: 1.9
nodePools:
- name: pool-1
  type: n1-standard-1
  size: 2
  labels:
    app: my-app
volumes:
- name: my-vol
  size: 10
```

## Step 1 - Generate GCDM template

Download the template from [kubeformation.sh](https://kubeformation.sh)

**or**

Generate the template using [CLI](../cli.md):
```bash
$ mkdir templates
$ kubeformation -f cluster.yaml -o templates
```

This will give us the following files:
- `gke-cluster.yaml`
- `gke-cluster.jinja`
- `k8s-volumes.yaml` (only if `volumes` are present in the spec)

## Step 2 - Add parameters

The following parameters which are provider specific need to be added:
- `ZONE` - GCP zone where the cluster has to be created
- `PROJECT` - GCP project

Open `gke-cluster.yaml` file and add the required GCP zone and project name by
replacing `ZONE` and `PROJECT`:

```yaml
imports:
- path: gke-cluster.jinja

resources:
- name:  my-cluster
  type: gke-cluster.jinja
  properties:
    name: my-cluster
    project: PROJECT
    zone: ZONE
```

## Step 3 - Create the cluster

Create the cluster (and any disks) as defined by `gke-cluster.yaml`:

```bash
$ gcloud deployment-manager deployments create my-cluster --config gke-cluster.yaml
```

That's it! The GKE cluster will be created.

Get `kubectl` context to connect to the cluster:

```bash
$ gcloud container clusters get-credentials my-cluster --zone <zone> --project <project>
```

## Step 4 - Create K8s Persistent Volumes

If the cluster spec also contains volumes, along with underlying disks, the
Kubernetes PV & PVC objects also have to be created, so that the disks can be
used by other k8s deployments etc.

```bash
$ kubectl create -f k8s-volumes.yaml
```

## Tearing down

Delete the deployment to tear down the cluster (and disks):

```bash
$ gcloud deployment-manager deployments delete my-cluster
```
