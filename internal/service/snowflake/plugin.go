package snowflake

import "github.com/calmonr/scaleid/pkg/plugin"

const PluginName = "snowflake"

func Plugin() *plugin.Plugin {
	return plugin.New(PluginName, "1.0.0", Runnable(), plugin.WithFlagSet(FlagSet()))
}
