# Version number
VERSION=$(shell ./tools/image-tag | cut -d, -f 1)

GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO_OPT= -mod vendor -ldflags "-X main.Branch=$(GIT_BRANCH) -X main.Revision=$(GIT_REVISION) -X main.Version=$(VERSION)"

.PHONY: cli build-go
cli: build-go

build-go:
	goreleaser build --single-target --snapshot --clean

.PHONY: gen-pkl
gen-pkl:
	PKL_EXEC=${PWD}/pkl pkl-gen-go pkg/config/pkl/AppCopnfig.pkl --generator-settings pkg/config/pkl/generator-settings.pkl

.PHONY: docker
docker:
	 goreleaser --pl