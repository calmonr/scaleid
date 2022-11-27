package config

import (
	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/internal/cli"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string
	Path     struct {
		Plugin string
	}
	GRPCServer struct {
		Address, Network string
		TLS              struct {
			Enabled                         bool
			CertPath, KeyPath, ClientCAPath string
		}
	}
	Plugin struct {
		Shared struct {
			Enabled bool
		}
	}
}

func (c *Config) Fill(v *viper.Viper) {
	// general
	c.LogLevel = v.GetString(cli.LogLevelFlag)

	// path
	c.Path.Plugin = v.GetString(app.PrefixPath + app.SuffixPlugin)

	// grpc server
	c.GRPCServer.Address = v.GetString(app.PrefixGRPCServer + app.SuffixAddress)
	c.GRPCServer.Network = v.GetString(app.PrefixGRPCServer + app.SuffixNetwork)
	c.GRPCServer.TLS.Enabled = v.GetBool(app.PrefixGRPCServerTLS + app.SuffixEnabled)
	c.GRPCServer.TLS.CertPath = v.GetString(app.PrefixGRPCServerTLS + app.SuffixCertPath)
	c.GRPCServer.TLS.KeyPath = v.GetString(app.PrefixGRPCServerTLS + app.SuffixKeyPath)
	c.GRPCServer.TLS.ClientCAPath = v.GetString(app.PrefixGRPCServerTLS + app.SuffixClientCAPath)

	// plugin
	c.Plugin.Shared.Enabled = v.GetBool(app.PrefixPluginShared + app.SuffixEnabled)
}
