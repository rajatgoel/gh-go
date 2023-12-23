setup:
        brew install go bufbuild/buf/buf golangci-lint sqlc goreleaser/tap/goreleaser fd
        go install golang.org/x/tools/cmd/goimports@latest

gen:
        buf generate proto
        fd sqlc.yaml . | xargs sqlc generate -f
        goimports -w .

lint: gen
        buf lint proto
        buf breaking proto --against '.git#branch=main,subdir=proto' || true
        golangci-lint run ./...

test: lint
        go vet ./...
        go test ./...
