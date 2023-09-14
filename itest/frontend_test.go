package itest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	frontend_pb "github.com/rajatgoel/gh-go/gen/frontend/v1"
	"github.com/rajatgoel/gh-go/internal/frontend"
)

func TestStub(t *testing.T) {
	h := frontend.New()
	_, err := h.Stub(context.Background(), &frontend_pb.StubRequest{})
	require.NoError(t, err)
}
