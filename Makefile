BASE_REGISTRY ?= docker.io/

.PHONY: lint test build

lint:
	@echo ">>running lint"
	@docker run --rm -t				\
	-v $(PWD):/app -w /app 			\
	$(BASE_REGISTRY)golangci/golangci-lint:v1.56.2 golangci-lint -v run --config=./.golangci.yml

test:
	@go test -v ./...

build:
	@go build -v ./...