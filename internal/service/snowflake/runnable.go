package snowflake

import (
	"fmt"

	"github.com/calmonr/scaleid/pkg/clock"
	"github.com/calmonr/scaleid/pkg/generator/snowflake"
	"github.com/calmonr/scaleid/pkg/plugin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Runnable() plugin.Runnable {
	return func(l *zap.Logger, v *viper.Viper) (plugin.GRPCServiceRegisterer, error) {
		c := Config{}
		c.Fill(v)

		l.Info("snowflake plugin is running", zap.Any("config", c))

		g, err := snowflake.NewGenerator(&clock.System{}, c.DatacenterID, c.WorkerID)
		if err != nil {
			return nil, fmt.Errorf("could not create snowflake generator: %w", err)
		}

		return NewService(g), nil
	}
}
