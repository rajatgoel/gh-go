package api

import (
	"context"

	frontend_pb "github.com/rajatgoel/gh-go/gen/frontend/v1"
)

type handler struct {
	frontend_pb.UnimplementedFrontendServiceServer
}

func (h *handler) Stub(context.Context, *frontend_pb.StubRequest) (*frontend_pb.StubResponse, error) {
	return &frontend_pb.StubResponse{}, nil
}

func New() frontend_pb.FrontendServiceServer {
	return &handler{}
}
