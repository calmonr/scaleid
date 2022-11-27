package credential

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

var ErrAppendCertsFromPEM = fmt.Errorf("could not append client certs")

func New(certPath, keyPath, clientCAPath string) (credentials.TransportCredentials, error) {
	config, err := NewTLSConfig(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not create tls config: %w", err)
	}

	if clientCAPath != "" {
		if err := SetTLSConfigClientCA(config, clientCAPath); err != nil {
			return nil, fmt.Errorf("could not change tls config to mutual: %w", err)
		}

		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return credentials.NewTLS(config), nil
}

func NewTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not load server key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func SetTLSConfigClientCA(config *tls.Config, clientCAPath string) error {
	cert, err := os.ReadFile(clientCAPath)
	if err != nil {
		return fmt.Errorf("could not read ca certificate: %w", err)
	}

	pool := x509.NewCertPool()

	if !pool.AppendCertsFromPEM(cert) {
		return ErrAppendCertsFromPEM
	}

	config.ClientCAs = pool

	return nil
}
