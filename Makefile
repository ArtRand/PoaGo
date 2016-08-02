# Use these to override the VERSION and REVISION constants with git tags
VERSION := $(shell git describe --tags)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/ArtRand/PoaGo/VERSION=$(VERSION) -X github.com/ArtRand/PoaGo/REVISION=$(REVISION) 

BINARY_NAME := poago

# Only build your project, and not libraries you vendor using govendor
FILES := $(shell go list ./... | grep -v vendor)

all: test install

test:
	@echo "+$@"
	go test $(FILES) -v -cover

# Assumes GOOS is set in the calling environment or defaults to base arch
build: 
	@echo "+$@"
	go build -v -o '$(BINARY_NAME)_$(REVISION)' -ldflags '$(LDFLAGS)' PoaGo.go

install:
	@echo "+$@"
	go install -v -ldflags '$(LDFLAGS)' $(FILES)
