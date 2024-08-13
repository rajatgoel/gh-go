package frontend

import (
	"context"

	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
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
