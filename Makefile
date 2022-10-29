.PHONY: all
all: build test

.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v -cover -failfast -benchmem -bench=. ./...