run-integration-test:
	go test -coverpkg=./... -coverprofile=coverage.out ./... -v

coverage-report:
	go tool cover -html=coverage.out

format:
	gofmt -w ./cmd ./config ./internal ./routes ./database

vet:
	go vet ./cmd/... ./config/... ./internal/... ./routes/... ./database/...

