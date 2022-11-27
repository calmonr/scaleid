package plugin

import (
	"errors"
	"fmt"
	"path/filepath"
	"plugin"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	ErrPluginNotFound      = errors.New("plugin not found")
	ErrPluginSetupFuncCast = errors.New("could not cast setup function")
)

const SetupFuncName = "Setup"

type Storage struct {
	logger  *zap.Logger
	viper   *viper.Viper
	plugins map[string]*Plugin
}

func NewStorage(l *zap.Logger, v *viper.Viper) *Storage {
	return &Storage{
		logger:  l,
		viper:   v,
		plugins: make(map[string]*Plugin),
	}
}

func (s *Storage) Add(p *Plugin) {
	s.plugins[p.ID()] = p
}

func (s *Storage) Get(id string) (*Plugin, error) {
	p, ok := s.plugins[id]
	if !ok {
		return nil, ErrPluginNotFound
	}

	return p, nil
}

func (s *Storage) All() map[string]*Plugin {
	plugins := make(map[string]*Plugin, len(s.plugins))

	for k, v := range s.plugins {
		plugins[k] = v
	}

	return plugins
}

func (s *Storage) LoadShared(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("could not get absolute path: %w", err)
	}

	pattern := filepath.Join(
		abs,
		fmt.Sprintf("*_%s_%s.so", runtime.GOOS, runtime.GOARCH),
	)

	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("could not glob files: %w", err)
	}

	for _, file := range files {
		plugin, err := plugin.Open(file)
		if err != nil {
			return fmt.Errorf("could not open plugin: %w", err)
		}

		symbol, err := plugin.Lookup(SetupFuncName)
		if err != nil {
			return fmt.Errorf("could not lookup symbol: %w", err)
		}

		setup, ok := symbol.(func(*Storage))
		if !ok {
			return ErrPluginSetupFuncCast
		}

		setup(s)
	}

	return nil
}

func (s *Storage) MergeFlagSets(f *pflag.FlagSet) {
	for _, p := range s.plugins {
		if p.FlagSet() != nil {
			f.AddFlagSet(p.FlagSet())
		}
	}
}

func (s *Storage) RegisterGRPCServices(svr *grpc.Server) error {
	for _, p := range s.All() {
		registerer, err := p.Runnable()(s.logger, s.viper)
		if err != nil {
			return fmt.Errorf("could not get plugin registerer: %w", err)
		}

		if err := registerer.Register(svr); err != nil {
			return fmt.Errorf("could not register service: %w", err)
		}
	}

	return nil
}
