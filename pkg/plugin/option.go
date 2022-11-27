package plugin

import "github.com/spf13/pflag"

type Option func(*Plugin)

func WithFlagSet(flags *pflag.FlagSet) Option {
	return func(p *Plugin) {
		p.flagSet = flags
	}
}
