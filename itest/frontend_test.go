package itest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/rajatgoel/gh-go/api"
)

func TestStub(t *testing.T) {
	h := api.New()
	_, err := h.Stub(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
}
