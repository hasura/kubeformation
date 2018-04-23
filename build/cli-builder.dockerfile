FROM golang:1.10

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# install gox
RUN go get github.com/mitchellh/gox

# setup the working directory
WORKDIR /go/src/github.com/hasura/kubeformation
