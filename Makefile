test-all: test test-integration

test:
	go test -race ./...

test-integration:
	go test --tags=integration ./integration/

fmt:
	go fmt ./...

lint:
	golangci-lint run