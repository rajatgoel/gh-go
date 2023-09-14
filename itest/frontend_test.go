package itest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	frontendpb "github.com/rajatgoel/gh-go/gen/frontend/v1"
	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

func TestStub(t *testing.T) {
	b, err := sqlbackend.New(context.Background())
	require.NoError(t, err)
	h := frontend.New(b)

	key, value := int64(1), "value"
	_, err = h.Put(context.Background(), &frontendpb.PutRequest{
		Key:   key,
		Value: value,
	})
	require.NoError(t, err)

	resp, err := h.Get(context.Background(), &frontendpb.GetRequest{
		Key: key,
	})
	require.NoError(t, err)
	require.Equal(t, value, resp.Value)
}
