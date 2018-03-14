VERSION = $(shell git describe --always --abbrev=4)
APP     = aptly-mirror

build:
	CGO_ENABLED=1 GOOS=linux go build \
	  -ldflags "-linkmode external -extldflags -static -X main.version=$(VERSION)" \
		-o $(APP) \
		main.go $(LIBS)
	strip $(APP)

install:
	go install .

.PHONY: build install

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
