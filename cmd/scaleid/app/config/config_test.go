package config_test

import (
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/cmd/scaleid/app/bootstrap"
	"github.com/calmonr/scaleid/cmd/scaleid/app/config"
	"github.com/calmonr/scaleid/internal/cli"
	"github.com/calmonr/scaleid/internal/unittest"
	"github.com/stretchr/testify/assert"
)

const (
	valueLogLevel                  = "debug"
	valuePathPlugin                = "/path/to/plugins"
	valueGRPCServerAddress         = ":50052"
	valueGRPCServerNetwork         = "udp"
	valueGRPCServerTLSEnabled      = "true"
	valueGRPCServerTLSCertPath     = "/path/to/cert"
	valueGRPCServerTLSKeyPath      = "/path/to/key"
	valueGRPCServerTLSClientCAPath = "/path/to/clientca"
	valuePluginSharedEnabled       = "true"
)

const (
	expectedLogLevel                  = valueLogLevel
	expectedPathPlugin                = valuePathPlugin
	expectedGRPCServerAddress         = valueGRPCServerAddress
	expectedGRPCServerNetwork         = valueGRPCServerNetwork
	expectedGRPCServerTLSEnabled      = true
	expectedGRPCServerTLSCertPath     = valueGRPCServerTLSCertPath
	expectedGRPCServerTLSKeyPath      = valueGRPCServerTLSKeyPath
	expectedGRPCServerTLSClientCAPath = valueGRPCServerTLSClientCAPath
	expectedPluginSharedEnabled       = true
)

func TestConfigFill(t *testing.T) {
	t.Parallel()

	t.Run("success: default", func(t *testing.T) {
		t.Parallel()

		f, err := bootstrap.FlagSet(app.Name, []string{})
		assert.NoError(t, err)

		v, err := bootstrap.Viper(app.Name, f)
		assert.NoError(t, err)

		c := config.Config{}
		c.Fill(v)

		assert.Equal(t, cli.DefaultLogLevel, c.LogLevel)
		assert.Equal(t, "", c.Path.Plugin)
		assert.Equal(t, app.DefaultGRPCServerAddress, c.GRPCServer.Address)
		assert.Equal(t, app.DefaultGRPCServerNetwork, c.GRPCServer.Network)
		assert.Equal(t, app.DefaultGRPCServerTLSEnabled, c.GRPCServer.TLS.Enabled)
		assert.Equal(t, "", c.GRPCServer.TLS.CertPath)
		assert.Equal(t, "", c.GRPCServer.TLS.KeyPath)
		assert.Equal(t, "", c.GRPCServer.TLS.ClientCAPath)
		assert.Equal(t, app.DefaultPluginSharedEnabled, c.Plugin.Shared.Enabled)
	})

	// nolint:paralleltest
	t.Run("success: environment variables", func(t *testing.T) {
		unittest.SetEnv(t, app.Name, cli.LogLevelFlag, valueLogLevel)
		unittest.SetEnv(t, app.Name, app.PrefixPath+app.SuffixPlugin, valuePathPlugin)
		unittest.SetEnv(t, app.Name, app.PrefixGRPCServer+app.SuffixAddress, valueGRPCServerAddress)
		unittest.SetEnv(t, app.Name, app.PrefixGRPCServer+app.SuffixNetwork, valueGRPCServerNetwork)
		unittest.SetEnv(t, app.Name, app.PrefixGRPCServerTLS+app.SuffixEnabled, valueGRPCServerTLSEnabled)
		unittest.SetEnv(t, app.Name, app.PrefixGRPCServerTLS+app.SuffixCertPath, valueGRPCServerTLSCertPath)
		unittest.SetEnv(t, app.Name, app.PrefixGRPCServerTLS+app.SuffixKeyPath, valueGRPCServerTLSKeyPath)
		unittest.SetEnv(t, app.Name, app.PrefixGRPCServerTLS+app.SuffixClientCAPath, valueGRPCServerTLSClientCAPath)
		unittest.SetEnv(t, app.Name, app.PrefixPluginShared+app.SuffixEnabled, valuePluginSharedEnabled)

		f, err := bootstrap.FlagSet(app.Name, []string{})
		assert.NoError(t, err)

		v, err := bootstrap.Viper(app.Name, f)
		assert.NoError(t, err)

		c := config.Config{}
		c.Fill(v)

		assert.Equal(t, expectedLogLevel, c.LogLevel)
		assert.Equal(t, expectedPathPlugin, c.Path.Plugin)
		assert.Equal(t, expectedGRPCServerAddress, c.GRPCServer.Address)
		assert.Equal(t, expectedGRPCServerNetwork, c.GRPCServer.Network)
		assert.Equal(t, expectedGRPCServerTLSEnabled, c.GRPCServer.TLS.Enabled)
		assert.Equal(t, expectedGRPCServerTLSCertPath, c.GRPCServer.TLS.CertPath)
		assert.Equal(t, expectedGRPCServerTLSKeyPath, c.GRPCServer.TLS.KeyPath)
		assert.Equal(t, expectedGRPCServerTLSClientCAPath, c.GRPCServer.TLS.ClientCAPath)
		assert.Equal(t, expectedPluginSharedEnabled, c.Plugin.Shared.Enabled)
	})

	t.Run("success: flags", func(t *testing.T) {
		t.Parallel()

		f, err := bootstrap.FlagSet(app.Name, []string{
			unittest.ComposeFlag(cli.LogLevelFlag, valueLogLevel),
			unittest.ComposeFlag(app.PrefixPath+app.SuffixPlugin, valuePathPlugin),
			unittest.ComposeFlag(app.PrefixGRPCServer+app.SuffixAddress, ":50052"),
			unittest.ComposeFlag(app.PrefixGRPCServer+app.SuffixNetwork, "udp"),
			unittest.ComposeFlag(app.PrefixGRPCServerTLS+app.SuffixEnabled, "true"),
			unittest.ComposeFlag(app.PrefixGRPCServerTLS+app.SuffixCertPath, "/path/to/cert"),
			unittest.ComposeFlag(app.PrefixGRPCServerTLS+app.SuffixKeyPath, "/path/to/key"),
			unittest.ComposeFlag(app.PrefixGRPCServerTLS+app.SuffixClientCAPath, "/path/to/clientca"),
			unittest.ComposeFlag(app.PrefixPluginShared+app.SuffixEnabled, "true"),
		})
		assert.NoError(t, err)

		v, err := bootstrap.Viper(app.Name, f)
		assert.NoError(t, err)

		c := config.Config{}
		c.Fill(v)

		assert.Equal(t, expectedLogLevel, c.LogLevel)
		assert.Equal(t, expectedPathPlugin, c.Path.Plugin)
		assert.Equal(t, expectedGRPCServerAddress, c.GRPCServer.Address)
		assert.Equal(t, expectedGRPCServerNetwork, c.GRPCServer.Network)
		assert.Equal(t, expectedGRPCServerTLSEnabled, c.GRPCServer.TLS.Enabled)
		assert.Equal(t, expectedGRPCServerTLSCertPath, c.GRPCServer.TLS.CertPath)
		assert.Equal(t, expectedGRPCServerTLSKeyPath, c.GRPCServer.TLS.KeyPath)
		assert.Equal(t, expectedGRPCServerTLSClientCAPath, c.GRPCServer.TLS.ClientCAPath)
		assert.Equal(t, expectedPluginSharedEnabled, c.Plugin.Shared.Enabled)
	})

	t.Run("success: config file", func(t *testing.T) {
		t.Parallel()

		f, err := bootstrap.FlagSet(app.Name, []string{
			unittest.ComposeFlag(cli.ConfigFileFlag, "testdata/config.yaml"),
		})
		assert.NoError(t, err)

		v, err := bootstrap.Viper(app.Name, f)
		assert.NoError(t, err)

		c := config.Config{}
		c.Fill(v)

		assert.Equal(t, expectedLogLevel, c.LogLevel)
		assert.Equal(t, expectedPathPlugin, c.Path.Plugin)
		assert.Equal(t, expectedGRPCServerAddress, c.GRPCServer.Address)
		assert.Equal(t, expectedGRPCServerNetwork, c.GRPCServer.Network)
		assert.Equal(t, expectedGRPCServerTLSEnabled, c.GRPCServer.TLS.Enabled)
		assert.Equal(t, expectedGRPCServerTLSCertPath, c.GRPCServer.TLS.CertPath)
		assert.Equal(t, expectedGRPCServerTLSKeyPath, c.GRPCServer.TLS.KeyPath)
		assert.Equal(t, expectedGRPCServerTLSClientCAPath, c.GRPCServer.TLS.ClientCAPath)
		assert.Equal(t, expectedPluginSharedEnabled, c.Plugin.Shared.Enabled)
	})
}
