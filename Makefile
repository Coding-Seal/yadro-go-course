.PHONY build:
build:
	@echo Building...
	go mod tidy
	go build -o myapp ./cmd/xkcd
.PHONY lint:
lint:
	@echo Linting...
	golangci-lint run --enable wsl