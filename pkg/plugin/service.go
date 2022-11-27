package plugin

import (
	"google.golang.org/grpc"
)

type GRPCServiceRegisterer interface {
	Register(*grpc.Server) error
}
