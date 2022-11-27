package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const Separator = "_"

func Replacer() viper.StringReplacer {
	return strings.NewReplacer("-", Separator, ".", Separator)
}

func Init(prefix string, flags *pflag.FlagSet) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.EnvKeyReplacer(Replacer()))
	v.AutomaticEnv()

	v.SetEnvPrefix(prefix)

	if err := v.BindPFlags(flags); err != nil {
		return nil, fmt.Errorf("could not bind flags: %w", err)
	}

	return v, nil
}

func ReadFromFile(v *viper.Viper, file string) error {
	v.SetConfigFile(file)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("could not read in config: %w", err)
	}

	return nil
}
