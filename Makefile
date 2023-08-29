# ----------------------------------
# Installation
install:
	@go install ./...

# ----------------------------------
# Linting
lint:
	@echo "Running golangci-lint..." && \
	golangci-lint run && \
	echo " > Done."

# ----------------------------------
# Tests
test: test-unit

test-unit:
	@go test -mod=readonly ./...

test-unit-cover:
	@go test -mod=readonly -coverprofile=coverage.out ./...
