setup:
        brew install bufbuild/buf/buf 

gen:
        buf generate proto

lint:
        buf lint proto

test:
        go test ./...
