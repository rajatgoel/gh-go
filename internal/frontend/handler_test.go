package frontend

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
)

type stubBackend struct {
	getFunc func(context.Context, int64) (string, error)
}

func (s stubBackend) Put(ctx context.Context, key int64, value string) error {
	return nil
}

func (s stubBackend) Get(ctx context.Context, key int64) (string, error) {
	if s.getFunc != nil {
		return s.getFunc(ctx, key)
	}
	return "", nil
}

func TestHandlerGetSuccess(t *testing.T) {
	t.Parallel()

	const (
		key   int64  = 42
		value string = "hello"
	)

	h := &handler{backend: stubBackend{
		getFunc: func(ctx context.Context, gotKey int64) (string, error) {
			if gotKey != key {
				t.Fatalf("unexpected key: got %d want %d", gotKey, key)
			}
			return value, nil
		},
	}}

	resp, err := h.Get(context.Background(), frontendpb.GetRequest_builder{Key: key}.Build())
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatalf("Get() response is nil, want non-nil")
	}

	if resp.GetValue() != value {
		t.Fatalf("Get() value = %q, want %q", resp.GetValue(), value)
	}
}

func TestHandlerGetNotFound(t *testing.T) {
	t.Parallel()

	h := &handler{backend: stubBackend{
		getFunc: func(ctx context.Context, key int64) (string, error) {
			return "", sql.ErrNoRows
		},
	}}

	resp, err := h.Get(context.Background(), frontendpb.GetRequest_builder{Key: 1}.Build())
	if err == nil {
		t.Fatalf("Get() error = nil, want non-nil")
	}

	if resp != nil {
		t.Fatalf("Get() response = %v, want nil", resp)
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Get() error = %v, want gRPC status error", err)
	}

	if st.Code() != codes.NotFound {
		t.Fatalf("Get() code = %v, want %v", st.Code(), codes.NotFound)
	}

	if st.Message() != sql.ErrNoRows.Error() {
		t.Fatalf("Get() message = %q, want %q", st.Message(), sql.ErrNoRows.Error())
	}
}

func TestHandlerGetUnexpectedError(t *testing.T) {
	t.Parallel()

	const errMsg = "boom"

	h := &handler{backend: stubBackend{
		getFunc: func(ctx context.Context, key int64) (string, error) {
			return "", errors.New(errMsg)
		},
	}}

	resp, err := h.Get(context.Background(), frontendpb.GetRequest_builder{Key: 1}.Build())
	if err == nil {
		t.Fatalf("Get() error = nil, want non-nil")
	}

	if resp != nil {
		t.Fatalf("Get() response = %v, want nil", resp)
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Get() error = %v, want gRPC status error", err)
	}

	if st.Code() != codes.Internal {
		t.Fatalf("Get() code = %v, want %v", st.Code(), codes.Internal)
	}

	if st.Message() != errMsg {
		t.Fatalf("Get() message = %q, want %q", st.Message(), errMsg)
	}
}
