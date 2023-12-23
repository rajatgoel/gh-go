setup:
        brew install go bufbuild/buf/buf golangci-lint sqlc goreleaser/tap/goreleaser fd

gen:
        buf generate proto
        fd sqlc.yaml . | xargs sqlc generate -f

lint: gen
        buf lint proto
        buf breaking proto --against '.git#branch=main,subdir=proto' || true
        golangci-lint run ./...

test: lint
        go vet ./...
        go test ./...
