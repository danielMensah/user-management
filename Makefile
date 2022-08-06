install.oapi:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen

install.go:
	go get -v ./...

install: install.go install.oapi

lint.go.fmt:
	go fmt ./...;

lint.go.golangci:
	golangci-lint run ./...;

lint.go.testfmt:
	test -z $(gofmt -s -l -w .);

lint: lint.go.fmt lint.go.golangci lint.go.testfmt

test.go.test:
	go test -cover -race -coverprofile=c.out ./...;

test.go.coverage: test.go.test
	go tool cover -html=c.out -o coverage.html;

test: test.go.test test.go.coverage

generate.api:
	oapi-codegen --config ./internal/api/config.yaml ./internal/api/api.yaml

local: down
	docker compose up --remove-orphans --build

down.db:
	docker rm -f db

down: down.db
	docker compose down --remove-orphans
