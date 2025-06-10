# ----------------------------------
# Installation
install:
	@go install ./...

# ----------------------------------
# Linting
LINT_IMAGE=golangci/golangci-lint:v2.1.6
lint:
	@echo "Running golangci-lint..." && \
	docker run -t --rm -v $(CURDIR):/app -w /app $(LINT_IMAGE) golangci-lint run

# ----------------------------------
# Tests
test: test-unit

test-unit:
	@go test -mod=readonly ./...

test-unit-cover:
	@go test -mod=readonly -coverprofile=coverage.out ./...
