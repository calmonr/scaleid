package credential_test

import (
	"crypto/tls"
	"testing"

	"github.com/calmonr/scaleid/internal/credential"
	"github.com/stretchr/testify/assert"
)

func TestNewTLSConfig(t *testing.T) {
	t.Parallel()

	t.Run("could not load server key pair", func(t *testing.T) {
		t.Parallel()

		const message = "could not load server key pair"

		cases := []struct {
			name, certPath, keyPath string
		}{
			{
				name:     "inexistent",
				certPath: "testdata/tls/inexistent-cert.pem",
				keyPath:  "testdata/tls/inexistent-key.pem",
			},
			{
				name:     "invalid",
				certPath: "testdata/tls/invalid-server-cert.pem",
				keyPath:  "testdata/tls/invalid-server-key.pem",
			},
		}

		for _, c := range cases {
			c := c

			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				_, err := credential.NewTLSConfig(c.certPath, c.keyPath)
				assert.ErrorContains(t, err, message)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		c, err := credential.NewTLSConfig("testdata/tls/server-cert.pem", "testdata/tls/server-key.pem")
		assert.NoError(t, err)

		assert.Equal(t, uint16(tls.VersionTLS12), c.MinVersion)
		assert.Equal(t, 1, len(c.Certificates))
	})
}

func TestSetTLSConfigClientCA(t *testing.T) {
	t.Parallel()

	t.Run("could not read ca certificate", func(t *testing.T) {
		t.Parallel()

		c := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		err := credential.SetTLSConfigClientCA(c, "testdata/tls/inexistent-ca-cert.pem")
		assert.ErrorContains(t, err, "could not read ca certificate")
	})

	t.Run("could not append client certs", func(t *testing.T) {
		t.Parallel()

		c := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		err := credential.SetTLSConfigClientCA(c, "testdata/tls/invalid-ca-cert.pem")
		assert.ErrorContains(t, err, "could not append client certs")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		c := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		err := credential.SetTLSConfigClientCA(c, "testdata/tls/ca-cert.pem")
		assert.NoError(t, err)

		assert.NotNil(t, c.ClientCAs)
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("could not create tls config", func(t *testing.T) {
		t.Parallel()

		_, err := credential.New(
			"testdata/tls/inexistent-cert.pem",
			"testdata/tls/inexistent-key.pem",
			"",
		)
		assert.ErrorContains(t, err, "could not create tls config")
	})

	t.Run("could not change tls config to mutual", func(t *testing.T) {
		t.Parallel()

		_, err := credential.New(
			"testdata/tls/server-cert.pem",
			"testdata/tls/server-key.pem",
			"testdata/tls/inexistent-ca-cert.pem",
		)
		assert.ErrorContains(t, err, "could not change tls config to mutual")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		_, err := credential.New(
			"testdata/tls/server-cert.pem",
			"testdata/tls/server-key.pem",
			"testdata/tls/ca-cert.pem",
		)
		assert.NoError(t, err)
	})
}
