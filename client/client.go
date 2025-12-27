// Package client provides a type-safe interface to the frontend service.
//
// By default the client connects using insecure transport credentials which is
// convenient for local development and tests. Applications that require
// transport security should provide their own credentials using
// WithTransportCredentials. For explicit insecure mode, WithInsecure is still
// available.
package client

import (
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	frontendpb "github.com/dynoinc/gh-go/proto/frontend/v1"
)

// Client provides a type-safe interface to the frontend service
type Client struct {
	conn   *grpc.ClientConn
	client frontendpb.FrontendServiceClient
}

type Option func(*clientConfig)

type clientConfig struct {
	target string
	dialer func(context.Context, string) (net.Conn, error)
	creds  credentials.TransportCredentials
}

// WithTarget sets the gRPC target
func WithTarget(target string) Option {
	return func(c *clientConfig) {
		c.target = target
	}
}

// WithDialer sets a custom dialer function
func WithDialer(dialer func(context.Context, string) (net.Conn, error)) Option {
	return func(c *clientConfig) {
		c.dialer = dialer
	}
}

// WithTransportCredentials sets the transport credentials used for the gRPC
// connection.
func WithTransportCredentials(creds credentials.TransportCredentials) Option {
	return func(c *clientConfig) {
		c.creds = creds
	}
}

// WithInsecure disables transport security (useful for testing)
func WithInsecure() Option {
	return func(c *clientConfig) {
		c.creds = insecure.NewCredentials()
	}
}

// New creates a new client with the given options
func New(opts ...Option) (*Client, error) {
	config := &clientConfig{
		target: "localhost:5051", // default target
	}

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	if config.target == "" {
		return nil, fmt.Errorf("target is required")
	}

	var dialOpts []grpc.DialOption

	// Add OpenTelemetry stats handler
	dialOpts = append(dialOpts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

	// Add custom dialer if provided
	if config.dialer != nil {
		dialOpts = append(dialOpts, grpc.WithContextDialer(config.dialer))
	}

	// Set transport credentials based on configuration
	// Use provided credentials or default to insecure for backwards
	// compatibility.
	creds := config.creds
	if creds == nil {
		creds = insecure.NewCredentials()
	}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))

	// Create gRPC connection
	conn, err := grpc.NewClient(config.target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to frontend service: %w", err)
	}

	client := frontendpb.NewFrontendServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the underlying gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Put stores a key-value pair
func (c *Client) Put(ctx context.Context, key int64, value string) error {
	req := frontendpb.PutRequest_builder{
		Key:   key,
		Value: value,
	}.Build()

	_, err := c.client.Put(ctx, req)
	return err
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key int64) (string, error) {
	req := frontendpb.GetRequest_builder{
		Key: key,
	}.Build()

	resp, err := c.client.Get(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetValue(), nil
}
