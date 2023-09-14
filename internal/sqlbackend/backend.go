package sqlbackend

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"

	"github.com/rajatgoel/gh-go/internal/sqlbackend/sqlgen"
)

type Backend interface {
	Put(ctx context.Context, key int64, value string)
	Get(ctx context.Context, key int64) string
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

func (s *sqliteBackend) Put(ctx context.Context, key int64, value string) {
	_, _ = s.q.Put(ctx, sqlgen.PutParams{
		Key:   key,
		Value: value,
	})
}

func (s *sqliteBackend) Get(ctx context.Context, key int64) string {
	get, _ := s.q.Get(ctx, key)
	return get.Value
}
