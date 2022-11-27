package plugin_test

import (
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/pkg/plugin"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		v := "1.0.0"

		p := plugin.New(app.Name, v, nil)

		assert.Equal(t, app.Name, p.ID())
		assert.Equal(t, v, p.Version())
		assert.Nil(t, p.Runnable())
	})
}

func TestWithFlagSet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		f := pflag.NewFlagSet(app.Name, pflag.ContinueOnError)
		f.String(app.Name, "", "")

		p := plugin.New(app.Name, "1.0.0", nil, plugin.WithFlagSet(f))

		assert.Equal(t, f, p.FlagSet())
	})
}
