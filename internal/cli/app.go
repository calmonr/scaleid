package cli

import (
	"github.com/spf13/pflag"
)

const (
	HelpFlag       = "help"
	VersionFlag    = "version"
	LogLevelFlag   = "log-level"
	ConfigFileFlag = "config-file"
)

const (
	DefaultLogLevel = "info"
)

func FlagSet() *pflag.FlagSet {
	f := new(pflag.FlagSet)

	f.BoolP(HelpFlag, "h", false, "Show help")
	f.BoolP(VersionFlag, "v", false, "Show version")
	f.StringP(LogLevelFlag, "l", DefaultLogLevel, "Minimal allowed log level")
	f.StringP(ConfigFileFlag, "c", "", "Path to config file")

	return f
}
