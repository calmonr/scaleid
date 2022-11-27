package app

import (
	"github.com/spf13/pflag"
)

const Name = "scaleid"

const (
	PrefixPath       = "path."
	PrefixGRPCServer = "grpc-server."
	PrefixTLS        = "tls."
	PrefixPlugin     = "plugin."
	PrefixInternal   = "internal."
	PrefixShared     = "shared."
)

const (
	PrefixGRPCServerTLS  = PrefixGRPCServer + PrefixTLS
	PrefixPluginInternal = PrefixPlugin + PrefixInternal
	PrefixPluginShared   = PrefixPlugin + PrefixShared
)

const (
	SuffixPlugin       = "plugin"
	SuffixAddress      = "address"
	SuffixNetwork      = "network"
	SuffixCertPath     = "cert-path"
	SuffixKeyPath      = "key-path"
	SuffixClientCAPath = "client-ca-path"
	SuffixEnabled      = "enabled"
)

const (
	DefaultGRPCServerAddress    = ":50051"
	DefaultGRPCServerNetwork    = "tcp"
	DefaultGRPCServerTLSEnabled = false
	DefaultPluginSharedEnabled  = false
)

func FlagSet() *pflag.FlagSet {
	f := new(pflag.FlagSet)

	// path
	f.String(PrefixPath+SuffixPlugin, "", "Path to shared plugins")

	// grpc server
	f.String(PrefixGRPCServer+SuffixAddress, DefaultGRPCServerAddress, "gRPC server address")
	f.String(PrefixGRPCServer+SuffixNetwork, DefaultGRPCServerNetwork, "gRPC server network")
	f.Bool(PrefixGRPCServerTLS+SuffixEnabled, DefaultGRPCServerTLSEnabled, "Enable TLS for gRPC server")
	f.String(PrefixGRPCServerTLS+SuffixCertPath, "", "gRPC server TLS certificate path")
	f.String(PrefixGRPCServerTLS+SuffixKeyPath, "", "gRPC server TLS private key path")
	f.String(PrefixGRPCServerTLS+SuffixClientCAPath, "", "gRPC server TLS client CA certificate path")

	// plugin
	f.Bool(PrefixPluginShared+SuffixEnabled, DefaultPluginSharedEnabled, "Enable shared plugins")

	return f
}
