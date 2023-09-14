package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	frontendpb "github.com/rajatgoel/gh-go/gen/frontend/v1"
	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

func main() {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()

	ctx := context.Background()
	backend, err := sqlbackend.New(ctx)
	if err != nil {
		log.Fatalf("failed to create backend: %v", err)
	}
	handler := frontend.New(backend)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	frontendpb.RegisterFrontendServiceServer(s, handler)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
