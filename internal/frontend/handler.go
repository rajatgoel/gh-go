package frontend

import (
	"context"

	"connectrpc.com/connect"

	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
	frontendv1connect "github.com/rajatgoel/gh-go/proto/frontend/v1/v1connect"
)

type handler struct {
	frontendv1connect.UnimplementedFrontendServiceHandler

	backend sqlbackend.Backend
}

func (h *handler) Put(
	ctx context.Context,
	req *connect.Request[frontendpb.PutRequest],
) (*connect.Response[frontendpb.PutResponse], error) {
	h.backend.Put(ctx, req.Msg.Key, req.Msg.Value)
	return connect.NewResponse(&frontendpb.PutResponse{}), nil
}

func (h *handler) Get(
	ctx context.Context,
	req *connect.Request[frontendpb.GetRequest],
) (*connect.Response[frontendpb.GetResponse], error) {
	value := h.backend.Get(ctx, req.Msg.Key)
	return connect.NewResponse(&frontendpb.GetResponse{Value: value}), nil
}

func New(backend sqlbackend.Backend) frontendv1connect.FrontendServiceHandler {
	return &handler{backend: backend}
}
