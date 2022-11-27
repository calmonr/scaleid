package plugin_test

import (
	"context"
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/internal/unittest"
	"github.com/calmonr/scaleid/pkg/plugin"
	testv1 "github.com/calmonr/scaleid/pkg/proto/test/v1"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestStorageAdd(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		p := plugin.New(app.Name, "1.0.0", nil)

		s.Add(p)

		assert.Len(t, s.All(), 1)
	})
}

func TestStorageGet(t *testing.T) {
	t.Parallel()

	t.Run("plugin not found", func(t *testing.T) {
		t.Parallel()

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		_, err := s.Get("test")
		assert.ErrorIs(t, err, plugin.ErrPluginNotFound)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		p := plugin.New(app.Name, "1.0.0", nil)

		s.Add(p)

		p, err := s.Get(app.Name)
		assert.NoError(t, err)

		assert.Equal(t, app.Name, p.ID())
	})
}

func TestStorageMergeFlagSets(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		const (
			t1 = "test1"
			t2 = "test2"
		)

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		f1 := pflag.NewFlagSet(t1, pflag.ContinueOnError)
		f1.String(t1, "", "")

		p1 := plugin.New(t1, "1.0.0", nil, plugin.WithFlagSet(f1))

		f2 := pflag.NewFlagSet(t2, pflag.ContinueOnError)
		f2.String(t2, "", "")

		p2 := plugin.New(t2, "1.0.0", nil, plugin.WithFlagSet(f2))

		s.Add(p1)
		s.Add(p2)

		f3 := pflag.NewFlagSet("", pflag.ContinueOnError)

		s.MergeFlagSets(f3)

		assert.NotNil(t, f3.Lookup(t1))
		assert.NotNil(t, f3.Lookup(t2))
	})
}

// ...

type serviceError struct {
	plugin.GRPCServiceRegisterer
}

var _ plugin.GRPCServiceRegisterer = (*serviceError)(nil)

func (s serviceError) Register(*grpc.Server) error {
	return unittest.ErrDummy
}

// ...

type serviceSuccess struct {
	testv1.UnimplementedTestServiceServer
	plugin.GRPCServiceRegisterer
}

var _ plugin.GRPCServiceRegisterer = (*serviceSuccess)(nil)

func (s serviceSuccess) Register(svr *grpc.Server) error {
	testv1.RegisterTestServiceServer(svr, s)

	return nil
}

func (s serviceSuccess) Test(ctx context.Context, req *testv1.TestRequest) (*testv1.TestResponse, error) {
	return &testv1.TestResponse{}, nil
}

// ...

func TestStorageRegisterGRPCServices(t *testing.T) {
	t.Parallel()

	t.Run("could not get plugin registerer", func(t *testing.T) {
		t.Parallel()

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		p := plugin.New("test", "1.0.0", func(l *zap.Logger, v *viper.Viper) (plugin.GRPCServiceRegisterer, error) {
			return nil, unittest.ErrDummy
		})

		s.Add(p)

		err := s.RegisterGRPCServices(nil)
		assert.ErrorIs(t, err, unittest.ErrDummy)
	})

	t.Run("could not register service", func(t *testing.T) {
		t.Parallel()

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		p := plugin.New("test", "1.0.0", func(l *zap.Logger, v *viper.Viper) (plugin.GRPCServiceRegisterer, error) {
			return &serviceError{}, nil
		})

		s.Add(p)

		err := s.RegisterGRPCServices(nil)
		assert.ErrorIs(t, err, unittest.ErrDummy)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := plugin.NewStorage(zap.NewNop(), viper.New())

		p := plugin.New("test", "1.0.0", func(l *zap.Logger, v *viper.Viper) (plugin.GRPCServiceRegisterer, error) {
			return &serviceSuccess{}, nil
		})

		s.Add(p)

		svr := grpc.NewServer()

		err := s.RegisterGRPCServices(svr)
		assert.NoError(t, err)

		assert.Contains(t, svr.GetServiceInfo(), testv1.TestService_ServiceDesc.ServiceName)
	})
}
