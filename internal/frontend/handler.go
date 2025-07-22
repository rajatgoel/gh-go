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
	if err := h.backend.Put(ctx, req.Msg.GetKey(), req.Msg.GetValue()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&frontendpb.PutResponse{}), nil
}

func (h *handler) Get(
	ctx context.Context,
	req *connect.Request[frontendpb.GetRequest],
) (*connect.Response[frontendpb.GetResponse], error) {
	value, err := h.backend.Get(ctx, req.Msg.GetKey())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(frontendpb.GetResponse_builder{Value: value}.Build()), nil
}

func New(backend sqlbackend.Backend) frontendv1connect.FrontendServiceHandler {
	return &handler{backend: backend}
}
