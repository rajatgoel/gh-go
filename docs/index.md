---
hide:
  - navigation
  - toc
---

# gh-go

A Github template to kick start a new Go module.

## Integrations

- [x] [buf](https://buf.build/)
    * [ ] Check if protobuf code needs to be generated before check-in.
- [x] [sqlc](https://sqlc.dev/) for SQL backend.
- [x] [golangci-lint](https://golangci-lint.run/) for linters.
- [x] [goreleaser](https://goreleaser.com/) for releases.

## Project layout

```yaml
    # mkdocs
    mkdocs.yml    # The configuration file.
    docs/
        blog/     # Development log
        index.md  # The documentation homepage.
        ...       # Other markdown pages, images and other files.

    # Go
    cmd/          # Binaries
        frontend/
        tools/
        ...
    internal/
        frontend/ # Internal libraries
        ...
    itest/        # Integration tests
        ...

    # Protobuf
    proto/
    gen/

    justfile      # Helper commands 
```
