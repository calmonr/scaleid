package grpctest

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func Connection(ctx context.Context, listener *bufconn.Listener) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithContextDialer(BufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	c, err := grpc.DialContext(ctx, "bufnet", opts...)
	if err != nil {
		return nil, fmt.Errorf("could not dial bufconn listener: %w", err)
	}

	return c, nil
}
