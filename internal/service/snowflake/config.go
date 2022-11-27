package snowflake

import (
	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/spf13/viper"
)

type Config struct {
	DatacenterID, WorkerID uint32
}

func (c *Config) Fill(v *viper.Viper) {
	c.DatacenterID = v.GetUint32(app.PrefixPluginInternal + PrefixSnowflake + SuffixDatacenterID)
	c.WorkerID = v.GetUint32(app.PrefixPluginInternal + PrefixSnowflake + SuffixWorkerID)
}
