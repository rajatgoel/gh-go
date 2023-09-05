package itest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rajatgoel/gh-go/api"
	frontend_pb "github.com/rajatgoel/gh-go/gen/frontend/v1"
)

func TestStub(t *testing.T) {
	h := api.New()
	_, err := h.Stub(context.Background(), &frontend_pb.StubRequest{})
	require.NoError(t, err)
}
