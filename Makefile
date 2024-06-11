BINARY_NAME=xkcd
PORT=8090

.PHONY: build
build:
	@echo Building...
	@go mod tidy
	go build -o ${BINARY_NAME} ./cmd/xkcd
.PHONY: format
format:
	@gofumpt -w .
	@wsl -fix ./...
.PHONY: sec
sec:
	@govulncheck
	@trivy fs .
.PHONY: lint
lint:
	@echo Linting...
	@golangci-lint run
.PHONY: bench
bench:
	@echo Running benchmark ...
	@go test -bench=. ./...
.PHONY: test
test:
	@echo Running tests ...
	@go test -race -coverprofile test/out/cover.out ./...
	@go tool cover -html=test/out/cover.out
.PHONY: e2e
e2e: build
	@echo Running e2e tests...
	@./${BINARY_NAME} -p=${PORT} 2>&1 1>/dev/null &
	@sleep 3s;
	@python3 test/e2e/update.py ${PORT};
	@python3 test/e2e/pics.py ${PORT};
	@kill $$(lsof -t -i:${PORT})
.PHONY: web
web:
	@echo Building...
	@go mod tidy
	go build -o web-server ./cmd/web
