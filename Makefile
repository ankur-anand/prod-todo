# Run all tests.
.PHONY: test
test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -cover -race ./...

