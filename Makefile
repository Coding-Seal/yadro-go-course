.PHONY build:
build:
	@echo Building...
	go build
.PHONY lint:
lint:
	@echo Linting...
	golangci-lint run