package runnable

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/calmonr/scaleid/cmd/scaleid/app/config"
	"github.com/calmonr/scaleid/internal/command"
	"github.com/calmonr/scaleid/internal/credential"
	"github.com/calmonr/scaleid/pkg/graceful"
	"github.com/calmonr/scaleid/pkg/plugin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Root(l *zap.Logger, c config.Config, s *plugin.Storage) command.Runnable {
	return func(cmd *cobra.Command, args []string) error {
		var opts []grpc.ServerOption

		if c.GRPCServer.TLS.Enabled {
			c, err := credential.New(c.GRPCServer.TLS.CertPath, c.GRPCServer.TLS.KeyPath, c.GRPCServer.TLS.ClientCAPath)
			if err != nil {
				return fmt.Errorf("could not create tls credentials: %w", err)
			}

			opts = append(opts, grpc.Creds(c))
		}

		server := grpc.NewServer(opts...)

		{
			err := s.RegisterGRPCServices(server)
			if err != nil {
				return fmt.Errorf("could not register grpc services: %w", err)
			}
		}

		reflection.Register(server)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		g := graceful.New(
			ctx,
			func() error {
				l.Info("gracefully stopping the application")

				server.GracefulStop()

				return nil
			},
			func() error {
				l.Info("forcefully stopping the application")

				server.Stop()

				return nil
			},
			10*time.Minute, // nolint:gomnd
		)

		if err := g.Action(func() error {
			listen, err := net.Listen(c.GRPCServer.Network, c.GRPCServer.Address)
			if err != nil {
				return fmt.Errorf("could not listen on %s: %w", c.GRPCServer.Address, err)
			}

			l.Info("serving grpc server", zap.String("address", c.GRPCServer.Address))

			if err := server.Serve(listen); err != nil {
				return fmt.Errorf("could not serve grpc: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("could not wait for graceful shutdown: %w", err)
		}

		return nil
	}
}
