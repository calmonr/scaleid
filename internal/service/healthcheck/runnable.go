package healthcheck

import (
	"github.com/calmonr/scaleid/pkg/plugin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Runnable() plugin.Runnable {
	return func(l *zap.Logger, v *viper.Viper) (plugin.GRPCServiceRegisterer, error) {
		l.Info("healthcheck plugin is running")

		return &Service{}, nil
	}
}
