package command

import (
	"github.com/calmonr/scaleid/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRoot(name, description string, r Runnable, f *pflag.FlagSet) *cobra.Command {
	cmd := &cobra.Command{
		Use:           name,
		Short:         description,
		Version:       version.Get().String(),
		RunE:          r,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().AddFlagSet(f)

	return cmd
}
