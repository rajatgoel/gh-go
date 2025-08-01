gen:
        go tool buf generate proto
        find . -name sqlc.yaml | xargs go tool sqlc generate -f
        go tool goimports -local github.com/rajatgoel/gh-go -w .
        go mod tidy

lint: gen
        go tool buf lint proto
        go tool buf breaking proto --against '.git#branch=main,subdir=proto'
        go fix ./...
        go vet ./...
        go tool staticcheck ./...
        go tool govulncheck ./...
        go tool golangci-lint run ./...

test: lint
        go test ./...

docker: 
        docker build -t localhost/frontend .
