package main

import (
	"fmt"
	"log"
	"os"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/cmd/scaleid/app/bootstrap"
	"github.com/calmonr/scaleid/cmd/scaleid/app/config"
	"github.com/calmonr/scaleid/cmd/scaleid/app/runnable"
	"github.com/calmonr/scaleid/internal/command"
	"github.com/calmonr/scaleid/pkg/logger"
	"github.com/calmonr/scaleid/pkg/version"
	"github.com/spf13/cobra"
)

func main() {
	cmd, err := rootCommand(app.Name)
	if err != nil {
		log.Fatalf("could not create root command: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalf("could not execute command: %v", err)
	}
}

func rootCommand(name string) (*cobra.Command, error) {
	f, err := bootstrap.FlagSet(name, os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("could not init flags: %w", err)
	}

	v, err := bootstrap.Viper(name, f)
	if err != nil {
		return nil, fmt.Errorf("could not init viper: %w", err)
	}

	c := config.Config{}
	c.Fill(v)

	l, err := logger.New(c.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	defer func() {
		// it should be safe to ignore according to https://github.com/uber-go/zap/issues/328
		_ = l.Sync()
	}()

	s := bootstrap.PluginStorage(l, v)

	if c.Plugin.Shared.Enabled {
		l.Warn("shared plugins are enabled. make sure you trust the plugins being loaded")

		if err := s.LoadShared(c.Path.Plugin); err != nil {
			return nil, fmt.Errorf("could not load shared plugins: %w", err)
		}
	}

	s.MergeFlagSets(f)

	if err := v.BindPFlags(f); err != nil {
		return nil, fmt.Errorf("could not bind flags: %w", err)
	}

	description := "scaleid is a free and open-source distributed unique ID generator."

	cmd := command.NewRoot(name, description, runnable.Root(l, c, s), f)
	cmd.AddCommand(version.Command(cmd.OutOrStderr()))

	return cmd, nil
}
