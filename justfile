gen:
        go run github.com/bufbuild/buf/cmd/buf generate proto
        find . -name sqlc.yaml | xargs go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f

lint: gen
        go run github.com/bufbuild/buf/cmd/buf lint proto
        go run github.com/bufbuild/buf/cmd/buf breaking proto --against '.git#branch=main,subdir=proto'
        go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

test: lint
        go vet ./...
        go test ./...
