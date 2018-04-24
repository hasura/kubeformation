# Azure Container Service (AKS)

For AKS clusters, Kubeformation creates [Azure Resource
Manager](https://docs.microsoft.com/en-us/azure/azure-resource-manager/)
templates (**ARM**).

The provider string is `aks`.

## Pre-requisites

- [Azure CLI 2.0](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)
- [`kubeformation`](../cli.md) CLI (optional)
- An Azure Subscription (start for [free
  here](https://azure.microsoft.com/en-us/free/))
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) CLI

## Step 0 - Write `cluster.yaml`

PS: This step can be skipped if [kubeformation.sh](https://kubeformation.sh) is
used. Jump right ahead to [Step 1](#step-1-generate-arm-template).

Here's an example `cluster.yaml`:

```yaml
version: v1
name: my-cluster
provider: aks
k8sVersion: 1.8.1
nodePools:
- name: pool1
  type: Standard_D2_v2
  size: 2
  labels:
    app: my-app
volumes:
- name: my-vol
  size: 10
```

Note:
- nodePool name should only contain `[a-z0-9]` and start with and must start with a
  lowercase letter

## Step 1 - Generate ARM template

Download the template from [kubeformation.sh](https://kubeformation.sh)

**or**

Generate the template using [CLI](../cli.md):
```bash
$ mkdir templates
$ kubeformation -f cluster.yaml -o templates
```

This will give us the following files:
- `aks-deploy.json` (ARM template file, defines all the resources)
- `aks-params.json` (ARM parameters file, defines params used by template)
- `aks-disks.json` (ARM template for creating disks, created only if `volumes` are present)
- `k8s-volumes.yaml` (k8s pv/pvc objects, created only if `volumes` are present)

## Step 2 - Add parameters

The following parameters are required:

- `SSH-PUBLIC-KEY` - SSH public key to be added to each node in the cluster. (can be the user's
  public key from `~/.ssh/id_rsa.pub`)
- `SERVICE-PRINCIPAL-CLIENT-ID`
- `SERVICE-PRINCIPAL-CLIENT-SECRET`

To create service principal, execute the following command:

```bash
$ az ad sp create-for-rbac --name my-cluster-sp
```

The output will be a JSON, of which `appId` is the `SERVICE-PRINCIPAL-CLIENT-ID`
and `password` is `SERVICE-PRINCIPAL-CLIENT-SECRET`.

These parameters should be added to `aks-params.json` file, by replacing the placeholders.

## Step 3 - Create the cluster

- Create an Azure Resource Group
  ```bash
  $ az group create -l westeurope -n my-resource-group
  ```
- Create the cluster
  ```bash
  $ az group deployment create -n my-cluster -g my-resource-group --template-file aks-deploy.json --parameters @aks-params.json 
  ```
- Create disks (if any):
  ```bash
  $ az group deployment create -n my-cluster-disks -g MC_my-resource-group_my-cluster_westeurope --template-file aks-disks.json
  ```
  
  The resource group mentioned here (stating with `MC_`) is automatically
  created by Azure to host all the infrastructure resources required by the
  Kubernetes cluster, including underlying VMs. Hence, the disks should be
  created in this resource group to be able to mount inside the cluster.
  
  Naming convention for this resource group is
  `MC_<original-resource-group-name>_<aks-cluster-name>_<location>`.

## Step 4 - Create K8s Persistent Volumes

If the cluster spec also contains volumes, along with underlying disks, the
Kubernetes PV & PVC objects also have to be created, so that the disks can be
used by other k8s deployments etc.

```bash
$ kubectl create -f k8s-volumes.yaml
```

## Tearing down

Delete the resource group to tear down the cluster and disks:

```bash
$ az group delete -n my-resource-group
```

We need to delete the resource group Azure created automatically to completely
tear down all the resources:

```bash
$ az group delete -n MC_my-resource-group_my-cluster_westeurope 
```
