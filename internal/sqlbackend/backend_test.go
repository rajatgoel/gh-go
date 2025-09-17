package sqlbackend

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLiteBackendPutOverwritesExistingValue(t *testing.T) {
	ctx := context.Background()

	backend, err := New(ctx)
	require.NoError(t, err)

	const key int64 = 42

	require.NoError(t, backend.Put(ctx, key, "first"))
	require.NoError(t, backend.Put(ctx, key, "second"))

	got, err := backend.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, "second", got)
}
