test: test-unit

test-unit:
	@go test -mod=readonly ./...
