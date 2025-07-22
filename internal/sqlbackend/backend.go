package sqlbackend

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"

	"github.com/rajatgoel/gh-go/internal/sqlbackend/sqlgen"
)

type Backend interface {
	Put(ctx context.Context, key int64, value string) error
	Get(ctx context.Context, key int64) (string, error)
}

//go:embed schema.sql
var ddl string

type sqliteBackend struct {
	q *sqlgen.Queries
}

func New(ctx context.Context) (Backend, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return &sqliteBackend{q: sqlgen.New(db)}, nil
}

func (s *sqliteBackend) Put(ctx context.Context, key int64, value string) error {
	_, err := s.q.Put(ctx, sqlgen.PutParams{
		Key:   key,
		Value: value,
	})
	return err
}

func (s *sqliteBackend) Get(ctx context.Context, key int64) (string, error) {
	get, err := s.q.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return get.Value, nil
}
