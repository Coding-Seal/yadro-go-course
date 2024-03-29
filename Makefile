.PHONY build:
build:
	@echo Building...
	go build -o myapp
.PHONY lint:
lint:
	@echo Linting...
	golangci-lint run --enable wsl