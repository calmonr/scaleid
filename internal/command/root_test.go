package command_test

import (
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewRoot(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		const (
			flagName           = "flag"
			commandName        = "command"
			commandDescription = "command description"
		)

		runnable := func(*cobra.Command, []string) error {
			return nil
		}

		f := pflag.NewFlagSet(app.Name, pflag.ContinueOnError)
		f.String(flagName, "", "")

		c := command.NewRoot(commandName, commandDescription, runnable, f)

		assert.Equal(t, commandName, c.Use)
		assert.Equal(t, commandDescription, c.Short)
		assert.NotNil(t, c.Flags().Lookup(flagName))
		assert.NoError(t, c.Execute())
	})
}
