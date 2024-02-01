run-integration-test:
	go test ./cmd/... -v

format:
	gofmt -w ./cmd ./config ./internal ./routes ./database

vet:
	go vet ./cmd/... ./config/... ./internal/... ./routes/... ./database/...

