NAME := horsebase
VERSION := 0.1
LDFLAGS := -X 'main.version=$(VERSION)

setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports

glide:
ifeq ($(shell command -v glide 2> /dev/null),)
    curl https://glide.sh/get | sh
endif

deps: setup
	glide install

## Update dependencies
update: setup
	glide update

build: deps
	go build -o bin/horsebase

install: build
	go install

test: deps
	go test $$(glide novendor)

## Lint
lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint --set_exit_status $$pkg || exit $$?; \
	done

.PHONY: setup deps build install update test lint
