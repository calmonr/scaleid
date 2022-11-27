package healthcheck

import (
	"github.com/calmonr/scaleid/pkg/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct{}

var _ plugin.GRPCServiceRegisterer = (*Service)(nil)

func (s Service) Register(svr *grpc.Server) error {
	grpc_health_v1.RegisterHealthServer(svr, health.NewServer())

	return nil
}
