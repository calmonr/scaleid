package snowflake

import (
	"context"
	"fmt"

	"github.com/calmonr/scaleid/pkg/generator/snowflake"
	"github.com/calmonr/scaleid/pkg/plugin"
	snowflakev1 "github.com/calmonr/scaleid/pkg/proto/snowflake/v1"
	"google.golang.org/grpc"
)

type Service struct {
	snowflakev1.UnimplementedSnowflakeServiceServer

	generator *snowflake.Generator
}

var _ plugin.GRPCServiceRegisterer = (*Service)(nil)

func NewService(generator *snowflake.Generator) *Service {
	return &Service{generator: generator}
}

func (s Service) Register(svr *grpc.Server) error {
	snowflakev1.RegisterSnowflakeServiceServer(svr, s)

	return nil
}

func (s Service) GenerateID(
	ctx context.Context,
	r *snowflakev1.GenerateIDRequest,
) (*snowflakev1.GenerateIDResponse, error) {
	id, err := s.generator.Generate()
	if err != nil {
		return nil, fmt.Errorf("could not generate snowflake id: %w", err)
	}

	return &snowflakev1.GenerateIDResponse{
		Id: uint64(id),
	}, nil
}
