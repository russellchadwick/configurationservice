PKG        := ./...
GOCC       ?= $(shell command -v "go")
GOFMT      ?= $(shell command -v "gofmt")
GO         ?= GOPATH=$(GOPATH) $(GOCC)

.PHONY: all
all: proto format build test check install

.PHONE: proto
proto:
	protoc --go_out=plugins=grpc:. proto/configuration.proto

.PHONY: format
format: $(GOCC)
	find . -iname '*.go' | grep -v '\./vendor' | xargs -n1 $(GOFMT) -w -s

.PHONY: build
build:
	$(GO) build $(GOFLAGS) -v

.PHONY: test
test:
	$(GO) test $(GOFLAGS) -i $(PKG)

.PHONY: check
check:
	gometalinter $(PKG) --concurrency=1 --deadline=5m --vendor --exclude pb.go

.PHONE: install
install:
	cd cmd/configurationserver
	go install
	cd cmd/getconfiguration
	go install
