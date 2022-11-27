package bootstrap_test

import (
	"os"
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/cmd/scaleid/app/bootstrap"
	"github.com/calmonr/scaleid/internal/cli"
	"github.com/calmonr/scaleid/internal/service/healthcheck"
	"github.com/calmonr/scaleid/internal/service/snowflake"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestFlagSet(t *testing.T) {
	t.Parallel()

	t.Run("could not parse flags", func(t *testing.T) {
		t.Parallel()

		_, err := bootstrap.FlagSet(app.Name, []string{"---"})
		assert.ErrorContains(t, err, "could not parse flags")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		f, err := bootstrap.FlagSet(app.Name, os.Args[1:])
		assert.NoError(t, err)

		assert.Equal(t, cli.FlagSet().NFlag()+app.FlagSet().NFlag(), f.NFlag())

		cli.FlagSet().VisitAll(func(flag *pflag.Flag) {
			assert.NotNil(t, f.Lookup(flag.Name))
		})

		app.FlagSet().VisitAll(func(flag *pflag.Flag) {
			assert.NotNil(t, f.Lookup(flag.Name))
		})
	})
}

func TestPluginStorage(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cases := []struct{ name string }{
			{name: healthcheck.PluginName},
			{name: snowflake.PluginName},
		}

		s := bootstrap.PluginStorage(nil, nil)

		for _, c := range cases {
			c := c

			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				_, err := s.Get(c.name)
				assert.NoError(t, err)
			})
		}
	})
}
