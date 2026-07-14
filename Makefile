go_bin:=$(shell pwd)/bin

# if you want to implement api updates then use bash `make`
# and extend/update/fix sdk code after

all: clean tools dep lint

tools:
	GOBIN=${go_bin} go install -mod=mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.9.0

clean:
	rm -rf ./bin
	mkdir ./bin

dep:
	go mod tidy

updatedep:
	go list -m -u all

lint:
	./bin/golangci-lint run ./...

test:
	go test --count=1 -v -race ./...

.PHONY: lint test updatedep