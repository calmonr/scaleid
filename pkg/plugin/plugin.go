package plugin

import (
	"github.com/spf13/pflag"
)

type Plugin struct {
	id, version string
	runnable    Runnable
	flagSet     *pflag.FlagSet
}

func New(id, version string, runnable Runnable, opts ...Option) *Plugin {
	p := &Plugin{
		id:       id,
		version:  version,
		runnable: runnable,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p Plugin) ID() string {
	return p.id
}

func (p Plugin) Version() string {
	return p.version
}

func (p Plugin) Runnable() Runnable {
	return p.runnable
}

func (p Plugin) FlagSet() *pflag.FlagSet {
	return p.flagSet
}
