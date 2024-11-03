package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

var (
	ErrNoCertOrKey = newInvalidTLSConfigError(
		"cert and key files are required",
	)
)

// NewTLSConfigFromFile creates a tls.Config from cert and key files
// returns an error if cert and key files are not provided.
// If caFile is provided, it will be used to create a RootCAs pool.
func NewTLSConfigFromFile(certFile, keyFile string, caFile string) (*tls.Config, error) {
	if certFile == "" || keyFile == "" {
		return nil, ErrNoCertOrKey
	}

	cfg := &tls.Config{}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		msg := fmt.Sprintf("failed to load cert and key files: %s", err)
		return nil, newInvalidTLSConfigError(msg)
	}

	cfg.Certificates = []tls.Certificate{cert}
	if caFile != "" {
		cp, err := NewCertPool(caFile)
		if err != nil {
			return nil, err
		}
		cfg.RootCAs = cp
	}

	return cfg, nil
}

func NewCertPool(caFile string) (*x509.CertPool, error) {
	pemCerts, err := os.ReadFile(caFile)
	if err != nil {
		return nil, newInvalidTLSConfigError(err.Error())
	}
	return NewCertPoolFromPEM(pemCerts)
}

func NewCertPoolFromPEM(pemCerts []byte) (*x509.CertPool, error) {
	cp := x509.NewCertPool()
	ok := cp.AppendCertsFromPEM(pemCerts)
	if !ok {
		return nil, newInvalidTLSConfigError("failed to parse root certificate")
	}
	return cp, nil
}

func NewSystemCertPool() (*x509.CertPool, error) {
	return x509.SystemCertPool()
}
