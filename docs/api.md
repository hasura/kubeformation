# kubeformation-api

## 1. JSON API

```http
POST /render HTTP/1.1
Content-Type: application/json
{
  "version": "v1",
  "k8sVersion": 1.9,
  "name": "my-cluster",
  "provider": "gke",
  "nodePools": [
    {
      "name": "pool1",
      "size": 1,
      "labels": {
        "app": "my-app"
      },
      "type": "n1-standard-1"
    }
  ],
  "volumes": [
    {
      "size": 10,
      "name": "my-vol"
    }
  ]
}
```

```http
HTTP/1.1 200 OK
Content-Type: application/json
{
  "gke-cluster.yaml": "...",
  "gke-cluster.json": "...",
  "k8s-volumes.yaml": "...",
}
```

## 2. File API

```http
POST /render?download=true HTTP/1.1
Content-Type: application/json
{
  "version": "v1",
  "k8sVersion": 1.9,
  "name": "my-cluster",
  "provider": "gke",
  "nodePools": [
    {
      "name": "pool1",
      "size": 1,
      "labels": {
        "app": "my-app"
      },
      "type": "n1-standard-1"
    }
  ],
  "volumes": [
    {
      "size": 10,
      "name": "my-vol"
    }
  ]
}
```

```http
HTTP/1.1 200 OK
Content-Type: application/zip

<binary-file-data>
```
