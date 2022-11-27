package grpctest

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/test/bufconn"
)

func BufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		l, err := listener.Dial()
		if err != nil {
			return nil, fmt.Errorf("could not dial bufconn listener: %w", err)
		}

		return l, nil
	}
}
