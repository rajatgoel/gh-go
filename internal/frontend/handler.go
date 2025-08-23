package frontend

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
)

type handler struct {
	frontendpb.UnimplementedFrontendServiceServer

	backend sqlbackend.Backend
}

func New(backend sqlbackend.Backend) frontendpb.FrontendServiceServer {
	return &handler{backend: backend}
}

func (h *handler) Put(
	ctx context.Context,
	req *frontendpb.PutRequest,
) (*frontendpb.PutResponse, error) {
	if err := h.backend.Put(ctx, req.GetKey(), req.GetValue()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &frontendpb.PutResponse{}, nil
}

func (h *handler) Get(
	ctx context.Context,
	req *frontendpb.GetRequest,
) (*frontendpb.GetResponse, error) {
	value, err := h.backend.Get(ctx, req.GetKey())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return frontendpb.GetResponse_builder{Value: value}.Build(), nil
}
