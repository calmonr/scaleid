package bootstrap

import (
	"fmt"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/internal/cli"
	"github.com/calmonr/scaleid/internal/config"
	"github.com/calmonr/scaleid/internal/service/healthcheck"
	"github.com/calmonr/scaleid/internal/service/snowflake"
	"github.com/calmonr/scaleid/pkg/plugin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func FlagSet(name string, args []string) (*pflag.FlagSet, error) {
	f := &pflag.FlagSet{
		SortFlags: true,
		ParseErrorsWhitelist: pflag.ParseErrorsWhitelist{
			UnknownFlags: true,
		},
	}

	f.Init(name, pflag.ContinueOnError)
	f.SetInterspersed(true)
	f.AddFlagSet(cli.FlagSet())
	f.AddFlagSet(app.FlagSet())

	if err := f.Parse(args); err != nil {
		return nil, fmt.Errorf("could not parse flags: %w", err)
	}

	return f, nil
}

func Viper(prefix string, flags *pflag.FlagSet) (*viper.Viper, error) {
	viper, err := config.Init(prefix, flags)
	if err != nil {
		return nil, fmt.Errorf("could not init config: %w", err)
	}

	if file := viper.GetString(cli.ConfigFileFlag); file != "" {
		if err := config.ReadFromFile(viper, file); err != nil {
			return nil, fmt.Errorf("could not read config file: %w", err)
		}
	}

	return viper, nil
}

func PluginStorage(logger *zap.Logger, viper *viper.Viper) *plugin.Storage {
	s := plugin.NewStorage(logger, viper)

	// internal plugins
	s.Add(healthcheck.Plugin())
	s.Add(snowflake.Plugin())

	return s
}
