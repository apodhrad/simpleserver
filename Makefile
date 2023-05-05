TARGET = $(CURDIR)/target

build:
	@mkdir -p $(TARGET)
	@go build -o $(TARGET)

test:
	@go clean -testcache
	@go test ./...

test-coverage:
	@mkdir -p $(TARGET)
	@go test ./... -coverprofile=$(TARGET)/coverage.out
	go tool cover -html=$(TARGET)/coverage.out
