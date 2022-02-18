GO ?= go
GOPATH := $(shell go env GOPATH)
INSTALL_FLAGS=-s -w
IMAGE_NAME=eu.gcr.io/test-234353/devices
IMAGE_VERSION = $(shell echo $${CI_BUILD_REF_NAME:=master}|sed 's/\//./')
GO_PACKAGES=$(shell go list ./... | grep -v /vendor/)

all: build

build:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 $(GO) build -ldflags="$(INSTALL_FLAGS)" -o out/linux/devices ./cmd/serve

test:
	$(GO) test $(GO_PACKAGES)

docker-build: build
	@docker build -t $(IMAGE_NAME):$(IMAGE_VERSION) .

push-push:
	@docker push $(IMAGE_NAME):$(IMAGE_VERSION)
