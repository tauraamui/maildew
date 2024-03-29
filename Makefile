.DEFAULT_GOAL := default

.PHONY: default
default: build

.PHONY: ci-run
ci-run: test lint build

.PHONY: run
run:
	go run cmd/dew/main.go

.PHONY: test
test:
	gotestsum ./...

.PHONY: test-verbose
test-verbose:
	gotestsum --format standard-verbose ./...

.PHONY: benchmark
benchmark:
	go test -run="none" -bench=. -benchmem ./...

.PHONY: coverage
coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -func coverage.out && rm coverage.out

.PHONY: coverage-persist
coverage-persist:
	go test -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -func coverage.out

.PHONY: install-gotestsum
install-gotestsum:
	go get github.com/gotestyourself/gotestsum

.PHONY: install-linter
install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.41.1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	mkdir -p builds && go build -o ./builds/dew ./cmd/dew/

.PHONY: build-with-mat-profile
build-with-mat-profile:
	mkdir -p builds && go build -tags matprofile -o ./builds/dragond ./cmd/dew/

.PHONY: run-build
run-build: build
	./builds/dew

