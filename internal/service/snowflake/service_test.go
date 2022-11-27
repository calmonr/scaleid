package snowflake_test

import (
	"context"
	"testing"

	"github.com/calmonr/scaleid/internal/grpctest"
	"github.com/calmonr/scaleid/internal/service/snowflake"
	snowflakev1 "github.com/calmonr/scaleid/pkg/proto/snowflake/v1"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestServiceGenerateID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		p := snowflake.Plugin()

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

		client := snowflakev1.NewSnowflakeServiceClient(c)

		{
			r, err := client.GenerateID(ctx, &snowflakev1.GenerateIDRequest{})
			assert.NoError(t, err)

			assert.NotZero(t, r.Id)
		}
	})
}
