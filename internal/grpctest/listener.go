package grpctest

import "google.golang.org/grpc/test/bufconn"

const BufSize = 1024 * 1024

func Listener() *bufconn.Listener {
	return bufconn.Listen(BufSize)
}
