# set variables
IMAGE_REGISTRY ?= "hasura"
API_IMAGE_NAME ?= "kubeformation-api"

VERSION := $(shell build/get-version.sh)
PWD := $(shell pwd)

# install dependencies
deps:
	$(shell curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)

# build api server locally
build-api-local:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
	-ldflags "-X github.com/hasura/kubeformation/pkg/cmd.version=${VERSION}" \
	-o _output/$(VERSION)/kubeformation-api cmd/api/kubeformation.go

# build api server docker image
build-api:
	docker build -t $(IMAGE_REGISTRY)/$(API_IMAGE_NAME):$(VERSION) \
	--build-arg VERSION=$(VERSION) \
	-f build/api.dockerfile .

# push api server docker image
push-api:
	docker push $(IMAGE_REGISTRY)/$(API_IMAGE_NAME):$(VERSION)

# build and push api server docker image
api: build-api push-api

# build cli locally, for all given platform/arch
build-cli:
	go get github.com/mitchellh/gox
	gox -ldflags "-X github.com/hasura/kubeformation/pkg/cmd.version=$(VERSION)" \
	-os="linux darwin windows" \
	-arch="amd64" \
	-output="_output/$(VERSION)/{{.OS}}-{{.Arch}}/kubeformation" \
	./cmd/cli/

# build cli inside a docker container
build-cli-in-docker:
	docker build -t kubeformation-cli-builder -f build/cli-builder.dockerfile build
	docker run --rm -it \
	-v $(PWD):/go/src/github.com/hasura/kubeformation \
	kubeformation-cli-builder \
	dep ensure && make build-cli

# run tests
test:
	go test -v ./...

# run tests in docker
test-in-docker:
	docker build -t kubeformation-cli-builder -f build/cli-builder.dockerfile build
	docker run --rm -it \
	-v $(PWD):/go/src/github.com/hasura/kubeformation \
	kubeformation-cli-builder \
	dep ensure && go test -v ./...
