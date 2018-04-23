# Contributing to Kubeformation

Thanks for your interest in Kuberformation. Follow the steps below to
contribute.

## Install dependencies

- [Go 1.10](https://golang.org/doc/install) (not required if using Docker)
- [Docker](https://docs.docker.com/install/) (optional)
- [Dep](https://golang.github.io/dep/docs/installation.html) for managing
  vendor-ed packages
- [GNU Make](https://www.gnu.org/software/make/) for build/test rules

## Fork the repo

Fork the repo and clone it into your `$GOPATH`:

```bash
$ mkdir -p $GOPATH/src/github.com/hasura
$ cd $GOPATH/src/github.com/hasura
$ git clone git@github.com:<username>/kubeformation.git
$ cd kubeformation
```

Add upstream as a remote to pull new changes:

```bash
$ git remote add upstream https://github.com/hasura/kubeformation.git
```

## Code

Kuberformation has two interfaces, a command line tool and a web api. All the
common code, including the conversion on spec to templates are organised in
[`pkg/`](https://github.com/hasura/kubeformation/tree/master/pkg) while the cli
and api code is in
[`cmd/`](https://github.com/hasura/kubeformation/tree/master/cmd). More details 
about code organisation and packages can be found in
[Godoc](https://godoc.org/github.com/hasura/kubeformation). 

### Spec

The spec that denotes the Kubernetes cluster is versioned. Any backwards
compatible changes can be made to the same version spec itself, but a breaking
change should result in a new version. 

Each version is defined and handled by an implementation of
[`VersionedSpecHandler`](https://godoc.org/github.com/hasura/kubeformation/pkg/spec#VersionedSpecHandler)
interface. `v1` spec can be found in 
[`pkg/spec/v1`](https://godoc.org/github.com/hasura/kubeformation/pkg/spec/v1)
as
[`ClusterSpec`](https://godoc.org/github.com/hasura/kubeformation/pkg/spec/v1#ClusterSpec).

The logic of converting the spec into the structure that provider takes is also
implemented in this package. For any breaking change to the spec, a new version
has to be created.

### Provider

Each managed Kubernetes service provider has it's own package, which contains
the templates and the context to render the template. The context is called
`ProviderSpec`.

Each provider should implement the interface called
[`Provider`](https://godoc.org/github.com/hasura/kubeformation/pkg/provider#Provider),
which has a method `MarshalFiles() (map[string][]byte, error)`. It returns a map
of filenames and content, the rendered template.

### CLI

`kubeformation` cli implements a command that takes the cluster spec, renders
the template for a provider and writes the files into a directory.

### API

The API accepts JSON by POST and responds with either a map of files and
contents in JSON or a ZIP file that downloads.

### Tests

Write unit tests as much as you can.

## Build & Test

- `go run cmd/api/kubeformation.go` - test out api
- `go run cmd/cli/kubeformation.go` - test out cli

All build and test rules are written in Makefile.

Without docker:

- `make build-api-local` builds api server using local go installation
- `make build-cli` builds the cli tool
- `make test` runs the unit tests

Using docker:

- `make build-api`
- `make build-cli-in-docker`
- `make test-in-docker`

## Submit PR

If everything looks well, submit a PR ☺
