package plugin

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Runnable func(*zap.Logger, *viper.Viper) (GRPCServiceRegisterer, error)
