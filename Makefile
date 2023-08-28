# ----------------------------------
# Installation
install:
	@go install ./...


# ----------------------------------
# Tests
test: test-unit

test-unit:
	@go test -mod=readonly ./...
