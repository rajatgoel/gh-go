gen:
        go tool buf generate proto
        find . -name sqlc.yaml | xargs go tool sqlc generate -f
        go mod tidy

lint: gen
        go tool buf lint proto
        go tool buf breaking proto --against '.git#branch=main,subdir=proto'
        go tool golangci-lint run ./...

test: lint
        go vet ./...
        go test ./...
