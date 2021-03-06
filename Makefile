GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./build

ORG := github.com/platform9
REPOPATH ?= $(ORG)/fluentd-operator

DOCKER_IMAGE_NAME = platform9/fluentd-operator
DOCKER_IMAGE_NAME_HELPER = platform9/fluentd-operator-helper

DOCKER_IMAGE_TAG ?= latest

LDFLAGS := -s -w -extldflags '-static'

SRCFILES := $(shell find ./pkg)

test:
	go test ./...

build/bin/fluentd-operator: test build/bin/fluentd-operator-$(GOOS)-$(GOARCH)
	cp build/bin/fluentd-operator-$(GOOS)-$(GOARCH) build/bin/fluentd-operator

build/bin/fluentd-operator-darwin-amd64: $(SRCFILES)
	GOARCH=amd64 CGO_ENABLED=0 GOOS=darwin go build --installsuffix cgo -a -o build/bin/fluentd-operator-darwin-amd64 cmd/manager/main.go

build/bin/fluentd-operator-linux-amd64: $(SRCFILES)
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build --installsuffix cgo -a -o build/bin/fluentd-operator-linux-amd64 cmd/manager/main.go


build/bin/fluentd-operator-helper:  build/bin/fluentd-operator-helper-$(GOOS)-$(GOARCH)
	cp build/bin/fluentd-operator-helper-$(GOOS)-$(GOARCH) build/bin/fluentd-operator-helper

build/bin/fluentd-operator-helper-darwin-amd64: $(SRCFILES)
	GOARCH=amd64 CGO_ENABLED=0 GOOS=darwin go build --installsuffix cgo -a -o build/bin/fluentd-operator-helper-darwin-amd64 cmd/helper/helper.go

build/bin/fluentd-operator-helper-linux-amd64: $(SRCFILES)
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build --installsuffix cgo -a -o build/bin/fluentd-operator-helper-linux-amd64 cmd/helper/helper.go



.PHONY: clean
clean:
	rm -fr build/

.PHONY: binary
binary: build/bin/fluentd-operator

.PHONY: image
image: test build/bin/fluentd-operator-linux-amd64
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: helper
helper: build/bin/fluentd-operator-helper

.PHONY: helper-image
helper-image: build/bin/fluentd-operator-helper-linux-amd64
	docker build -t $(DOCKER_IMAGE_NAME_HELPER) -f hack/helper/Dockerfile .
