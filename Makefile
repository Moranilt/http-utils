test:
	@go test ./...

test-cover:
	@go test -coverprofile=coverage.out ./...

cover-html:
	@go tool cover -html=coverage.out