package frontend

import (
	"context"

	frontendpb "github.com/rajatgoel/gh-go/gen/frontend/v1"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

type handler struct {
	frontendpb.UnimplementedFrontendServiceServer

	backend sqlbackend.Backend
}

func (h *handler) Put(
	ctx context.Context,
	req *frontendpb.PutRequest,
) (*frontendpb.PutResponse, error) {
	h.backend.Put(ctx, req.Key, req.Value)
	return &frontendpb.PutResponse{}, nil
}

func (h *handler) Get(
	ctx context.Context,
	req *frontendpb.GetRequest,
) (*frontendpb.GetResponse, error) {
	value := h.backend.Get(ctx, req.Key)
	return &frontendpb.GetResponse{Value: value}, nil
}

func New(backend sqlbackend.Backend) frontendpb.FrontendServiceServer {
	return &handler{backend: backend}
}
