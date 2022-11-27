package config_test

import (
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/internal/config"
	"github.com/calmonr/scaleid/internal/unittest"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	configKey   = "key"
	configValue = "value"
)

func TestInit(t *testing.T) {
	t.Parallel()

	// nolint:paralleltest
	t.Run("success: environment variables", func(t *testing.T) {
		unittest.SetEnv(t, app.Name, configKey, configValue)

		f := pflag.NewFlagSet(app.Name, pflag.ContinueOnError)

		c, err := config.Init(app.Name, f)
		assert.NoError(t, err)

		assert.Equal(t, configValue, c.GetString(configKey))
	})

	t.Run("success: flags", func(t *testing.T) {
		t.Parallel()

		f := pflag.NewFlagSet(app.Name, pflag.ContinueOnError)
		f.String(configKey, "", "")

		err := f.Parse([]string{
			unittest.ComposeFlag(configKey, configValue),
		})
		assert.NoError(t, err)

		c, err := config.Init(app.Name, f)
		assert.NoError(t, err)

		assert.Equal(t, configValue, c.GetString(configKey))
	})
}

func TestReadFromFile(t *testing.T) {
	t.Parallel()

	t.Run("could not read in config", func(t *testing.T) {
		t.Parallel()

		v := viper.New()

		err := config.ReadFromFile(v, "testdata/inexistent.yaml")
		assert.ErrorContains(t, err, "could not read in config")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		v := viper.New()

		err := config.ReadFromFile(v, "testdata/config.yaml")
		assert.NoError(t, err)

		assert.Equal(t, configValue, v.GetString(configKey))
	})
}
