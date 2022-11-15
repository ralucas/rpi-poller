.PHONY: all
all: build test

.PHONY: build
build:
	mkdir -p bin
	go build -v -o bin ./... 

.PHONY: test
test:
	go test -v -cover -failfast -benchmem -bench=. ./...