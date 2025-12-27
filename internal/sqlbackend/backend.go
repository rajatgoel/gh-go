package sqlbackend

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/dynoinc/gh-go/internal/sqlbackend/sqlgen"
)

type Backend interface {
	Put(ctx context.Context, key int64, value string) error
	Get(ctx context.Context, key int64) (string, error)
	Close(ctx context.Context) error
}

type sqliteBackend struct {
	db *sql.DB
	q  *sqlgen.Queries
}

func New(ctx context.Context) (Backend, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	// Apply migrations to set up the database schema
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &sqliteBackend{
		db: db,
		q:  sqlgen.New(db),
	}, nil
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

func (s *sqliteBackend) Close(context.Context) error {
	return s.db.Close()
}
