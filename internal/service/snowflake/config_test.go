package snowflake_test

import (
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/cmd/scaleid/app/bootstrap"
	"github.com/calmonr/scaleid/internal/cli"
	"github.com/calmonr/scaleid/internal/service/snowflake"
	"github.com/calmonr/scaleid/internal/unittest"
	"github.com/stretchr/testify/assert"
)

const (
	valueDatacenterID = "2"
	valueWorkerID     = "2"
)

const (
	expectedDatacenterID uint32 = 2
	expectedWorkerID     uint32 = 2
)

func TestConfigFill(t *testing.T) {
	t.Parallel()

	t.Run("success: default", func(t *testing.T) {
		t.Parallel()

		v, err := bootstrap.Viper(app.Name, snowflake.FlagSet())
		assert.NoError(t, err)

		c := snowflake.Config{}
		c.Fill(v)

		assert.Equal(t, snowflake.DefaultDatacenterID, c.DatacenterID)
		assert.Equal(t, snowflake.DefaultWorkerID, c.WorkerID)
	})

	// nolint:paralleltest
	t.Run("success: environment variables", func(t *testing.T) {
		unittest.SetEnv(t, app.Name, snowflake.PrefixPluginInternalSnowflake+snowflake.SuffixDatacenterID, valueDatacenterID)
		unittest.SetEnv(t, app.Name, snowflake.PrefixPluginInternalSnowflake+snowflake.SuffixWorkerID, valueWorkerID)

		v, err := bootstrap.Viper(app.Name, snowflake.FlagSet())
		assert.NoError(t, err)

		c := snowflake.Config{}
		c.Fill(v)

		assert.Equal(t, expectedDatacenterID, c.DatacenterID)
		assert.Equal(t, expectedWorkerID, c.WorkerID)
	})

	t.Run("success: flags", func(t *testing.T) {
		t.Parallel()

		f := snowflake.FlagSet()

		err := f.Parse([]string{
			unittest.ComposeFlag(snowflake.PrefixPluginInternalSnowflake+snowflake.SuffixDatacenterID, valueDatacenterID),
			unittest.ComposeFlag(snowflake.PrefixPluginInternalSnowflake+snowflake.SuffixWorkerID, valueWorkerID),
		})
		assert.NoError(t, err)

		v, err := bootstrap.Viper(app.Name, f)
		assert.NoError(t, err)

		c := snowflake.Config{}
		c.Fill(v)

		assert.Equal(t, expectedDatacenterID, c.DatacenterID)
		assert.Equal(t, expectedWorkerID, c.WorkerID)
	})

	t.Run("success: config file", func(t *testing.T) {
		t.Parallel()

		f := snowflake.FlagSet()
		f.AddFlagSet(cli.FlagSet())

		err := f.Parse([]string{
			unittest.ComposeFlag(cli.ConfigFileFlag, "testdata/config.yaml"),
		})
		assert.NoError(t, err)

		v, err := bootstrap.Viper(app.Name, f)
		assert.NoError(t, err)

		c := snowflake.Config{}
		c.Fill(v)

		assert.Equal(t, expectedDatacenterID, c.DatacenterID)
		assert.Equal(t, expectedWorkerID, c.WorkerID)
	})
}
