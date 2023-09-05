package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	frontend_pb "github.com/rajatgoel/gh-go/gen/frontend/v1"
)

type handler struct {
	frontend_pb.UnimplementedFrontendServiceServer
}

func (h *handler) Stub(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func New() frontend_pb.FrontendServiceServer {
	return &handler{}
}
