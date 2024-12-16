// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlgen

import (
	"context"
)

const get = `-- name: Get :one
SELECT "key", value FROM keyvalue
WHERE key = ? LIMIT 1
`

func (q *Queries) Get(ctx context.Context, key int64) (Keyvalue, error) {
	row := q.db.QueryRowContext(ctx, get, key)
	var i Keyvalue
	err := row.Scan(&i.Key, &i.Value)
	return i, err
}

const put = `-- name: Put :one
INSERT INTO keyvalue (
    key, value
) VALUES (
    ?, ?
)
RETURNING "key", value
`

type PutParams struct {
	Key   int64
	Value string
}

func (q *Queries) Put(ctx context.Context, arg PutParams) (Keyvalue, error) {
	row := q.db.QueryRowContext(ctx, put, arg.Key, arg.Value)
	var i Keyvalue
	err := row.Scan(&i.Key, &i.Value)
	return i, err
}
