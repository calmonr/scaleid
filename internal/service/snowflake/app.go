package snowflake

import (
	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/spf13/pflag"
)

const (
	PrefixSnowflake = "snowflake."
)

const (
	SuffixDatacenterID = "datacenter-id"
	SuffixWorkerID     = "worker-id"
)

const (
	PrefixPluginInternalSnowflake = app.PrefixPluginInternal + PrefixSnowflake
)

const (
	DefaultDatacenterID uint32 = 1
	DefaultWorkerID     uint32 = 1
)

func FlagSet() *pflag.FlagSet {
	f := new(pflag.FlagSet)

	f.Uint32(PrefixPluginInternalSnowflake+SuffixDatacenterID, DefaultDatacenterID, "Snowflake datacenter id")
	f.Uint32(PrefixPluginInternalSnowflake+SuffixWorkerID, DefaultWorkerID, "Snowflake worker id")

	return f
}
