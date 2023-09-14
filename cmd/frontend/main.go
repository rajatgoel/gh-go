package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	frontend_pb "github.com/rajatgoel/gh-go/gen/frontend/v1"
	"github.com/rajatgoel/gh-go/internal/frontend"
)

func main() {
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	frontend_pb.RegisterFrontendServiceServer(s, frontend.New())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
