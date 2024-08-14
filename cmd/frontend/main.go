package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendv1connect "github.com/rajatgoel/gh-go/proto/frontend/v1/v1connect"
)

func main() {
	port := flag.Int("port", 5051, "The server port")
	flag.Parse()

	ctx := context.Background()
	backend, err := sqlbackend.New(ctx)
	if err != nil {
		log.Fatalf("failed to create backend: %v", err)
	}
	path, handler := frontendv1connect.NewFrontendServiceHandler(frontend.New(backend))

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// health check
	checker := grpchealth.NewStaticChecker(frontendv1connect.FrontendServiceName)
	mux.Handle(grpchealth.NewHandler(checker))

	// reflect
	reflector := grpcreflect.NewStaticReflector(frontendv1connect.FrontendServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", *port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	slog.Info("starting server", "port", *port, "path", path)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start server: %v", err)
	}
}
