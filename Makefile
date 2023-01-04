.PHONY: all
all: build test

.PHONY: build
build:
	mkdir -p bin
	go build -v -o bin ./... 

.PHONY: test
test:
	go test --tags=unit -v -cover -failfast -benchmem -bench=. ./...

.PHONY: test-integration
test-integration:
	go test --tags=integration -v -cover -failfast -benchmem -bench=. ./...