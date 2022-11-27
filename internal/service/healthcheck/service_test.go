package healthcheck_test

import (
	"context"
	"testing"

	"github.com/calmonr/scaleid/internal/grpctest"
	"github.com/calmonr/scaleid/internal/service/healthcheck"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestServiceRegister(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		p := healthcheck.Plugin()

		r, err := p.Runnable()(zap.NewNop(), viper.New())
		assert.NoError(t, err)

		l := grpctest.Listener()
		s := grpc.NewServer()

		t.Cleanup(func() {
			s.Stop()
			l.Close()
		})

		{
			err := r.Register(s)
			assert.NoError(t, err)
		}

		go func() {
			err := s.Serve(l)
			assert.NoError(t, err)
		}()

		ctx := context.Background()

		c, err := grpctest.Connection(ctx, l)
		assert.NoError(t, err)

		t.Cleanup(func() {
			c.Close()
		})

		client := grpc_health_v1.NewHealthClient(c)

		{
			r, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
			assert.NoError(t, err)

			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, r.Status)
		}
	})
}
