FROM golang:1.10 as builder

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# setup the working directory
WORKDIR /go/src/github.com/hasura/kubeformation

# set arguments
ARG VERSION

# copy source code
COPY pkg pkg
COPY cmd cmd
# copy Gopkg files
COPY Gopkg.toml Gopkg.lock ./

# install dependencies
RUN dep ensure
RUN echo $VERSION

# build the source
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X github.com/hasura/kubeformation/pkg/cmd.version=${VERSION}" \
    -o /bin/kubeformation-api cmd/api/kubeformation.go

# use a minimal alpine image
FROM alpine:3.7
WORKDIR /bin
# copy the binary from builder
COPY --from=builder /bin/kubeformation-api /bin/kubeformation-api
# run the binary
CMD ["/bin/kubeformation-api"]
