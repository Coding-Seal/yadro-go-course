.PHONY build:
build:
	@echo Building...
	go mod tidy
	go build -o xkcd ./cmd/xkcd
.PHONY lint:
lint:
	@echo Linting...
	golangci-lint run
.PHONY bench:
bench:
	@echo Running benchmark ...
	go test -bench=. ./...
.PHONY test:
test:
	@echo Running tests ...
	go test -race -coverprofile cover.out ./...
	go tool cover -html=cover.out
.PHONY lint-strict:
lint-strict:
	@echo Linting...
format:
	gofumpt -w .
	wsl -fix ./...
.PHONY sec:
	govulncheck
	trivy fs .

