.PHONY: dep
dep:
	go mod tidy
	go mod download

.PHONY: gen
gen:
	protoc --go_out=./ protocol.proto
	protoc -I /usr/local/include -I . --gotag_out=:./ protocol.proto
	go generate ./...
	go run ./cmd/gen/gen.go

.PHONY: lint
lint:
	golangci-lint run --tests

.PHONY: test
test:
	GOMAXPROCS=4 go test ./... -p 4 -parallel 4 -count=1

.PHONY: build
build:
	go build -o api ./cmd/api/api.go

.PHONY: migrate
migrate:
	go run ./cmd/migrate/migrate.go up
