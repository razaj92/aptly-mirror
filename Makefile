NAME          := aptly-mirror
VERSION       := $(shell git describe --tags --abbrev=1)
LDFLAGS       := -linkmode external -extldflags -static -X 'main.version=$(VERSION)'


.PHONY: setup
setup:
	go get -u -v github.com/golang/dep/cmd/dep

.PHONY: build
build:
	dep ensure -v
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o $(NAME) main.go
	strip $(NAME)
