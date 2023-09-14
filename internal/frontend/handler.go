package frontend

import (
	"context"

	frontend_pb "github.com/rajatgoel/gh-go/gen/frontend/v1"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

type handler struct {
	frontend_pb.UnimplementedFrontendServiceServer

	backend sqlbackend.Backend
}

func (h *handler) Put(
	ctx context.Context,
	req *frontend_pb.PutRequest,
) (*frontend_pb.PutResponse, error) {
	h.backend.Put(ctx, req.Key, req.Value)
	return &frontend_pb.PutResponse{}, nil
}

func (h *handler) Get(
	ctx context.Context,
	req *frontend_pb.GetRequest,
) (*frontend_pb.GetResponse, error) {
	value := h.backend.Get(ctx, req.Key)
	return &frontend_pb.GetResponse{Value: value}, nil
}

func New(backend sqlbackend.Backend) frontend_pb.FrontendServiceServer {
	return &handler{backend: backend}
}
