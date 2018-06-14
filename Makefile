#!/usr/bin/env bash

## setup version variables
COMMIT := $(shell git rev-parse HEAD)
BUILDTIME := $(shell date +"%F-%T-%Z")
LDFLAGS := -X github.com/modeneis/uminer/src/version.version=v0.1

OPENCL_HEADERS=/usr/local/cuda-8.0/targets/x86_64-linux/include LIBOPENCL=/usr/local/cuda-8.0/targets/x86_64-linux/lib


.DEFAULT_GOAL := help
.PHONY: test test-race lint check cover install-linters format run build log help

PACKAGES = $(shell find ./src -type d -not -path '\./src')

test: ## Run tests
	go test ./src/... -timeout=1m -cover

test-race: ## Run tests with -race. Note: expected to fail, but look for "DATA RACE" failures specifically
	go test ./src/... -timeout=2m -race

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	gometalinter --deadline=3m -j 2 --disable-all --tests --vendor \
		-E deadcode \
		-E errcheck \
		-E gas \
		-E goconst \
		-E gofmt \
		-E goimports \
		-E golint \
		-E ineffassign \
		-E interfacer \
		-E maligned \
		-E megacheck \
		-E misspell \
		-E nakedret \
		-E structcheck \
		-E unconvert \
		-E unparam \
		-E varcheck \
		-E vet \
		./...

check: lint ## Run tests and linters

cover: ## Runs tests on ./src/ with HTML code coverage
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/alecthomas/gometalinter
	gometalinter --vendored-linters --install

format:  # Formats the code. Must have goimports installed (use make install-linters).
	# This sorts imports by [stdlib, 3rdpart]
	goimports -w -local github.com/modeneis/uminer ./src
	goimports -w -local github.com/modeneis/uminer ./main.go

	# This performs code simplifications
	gofmt -s -w ./src
	gofmt -s -w ./main.go

build:
	echo "will compile app";
	OPENCL_HEADERS=$OPENCL_HEADERS go build -o uminer -ldflags "$(LDFLAGS)" main.go;

start: build
	echo "will start mining";
	nohup ./uminer --cointype=SIA --url=stratum+tcp://fcn-xmr.pool.minergate.com:45590 --username=modeneis --password=test  &
	make log

log:
	tail -f ./nohup.out


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'



.PHONY: all test
