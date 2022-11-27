package healthcheck

import "github.com/calmonr/scaleid/pkg/plugin"

const PluginName = "healthcheck"

func Plugin() *plugin.Plugin {
	return plugin.New(PluginName, "1.0.0", Runnable())
}
