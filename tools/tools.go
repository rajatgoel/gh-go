//go:build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "golang.org/x/tools/cmd/deadcode"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
