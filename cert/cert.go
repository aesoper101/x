package cert

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/go-rootcerts"
	"net"
	"strings"
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
	return ConfigureTLS(
		&TLSConfig{
			CertFile: certFile,
			KeyFile:  keyFile,
			CAFile:   caFile,
		},
	)
}

// ConfigureTLS  configures a tls.Config from the given TLSConfig
func ConfigureTLS(tlsConfig *TLSConfig) (*tls.Config, error) {
	tlsClientConfig := createBaseTLSConfig(tlsConfig)

	if tlsConfig.Address != "" {
		if err := setServerName(tlsClientConfig, tlsConfig.Address); err != nil {
			return nil, err
		}
	}

	if err := configureCertificates(tlsClientConfig, tlsConfig); err != nil {
		return nil, err
	}

	if err := configureRootCertificates(tlsClientConfig, tlsConfig); err != nil {
		return nil, err
	}

	return tlsClientConfig, nil
}

// createBaseTLSConfig 创建基础的 tls.Config 对象
func createBaseTLSConfig(tlsConfig *TLSConfig) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
	}
}

// setServerName 设置 tls.Config 的 ServerName
func setServerName(tlsClientConfig *tls.Config, address string) error {
	server := address
	hasPort := strings.LastIndex(server, ":") > strings.LastIndex(server, "]")
	if hasPort {
		var err error
		server, _, err = net.SplitHostPort(server)
		if err != nil {
			return err
		}
	}
	tlsClientConfig.ServerName = server
	return nil
}

// configureCertificates 配置 tls.Config 的 Certificates
func configureCertificates(tlsClientConfig *tls.Config, tlsConfig *TLSConfig) error {
	if len(tlsConfig.CertPEM) != 0 && len(tlsConfig.KeyPEM) != 0 {
		tlsCert, err := tls.X509KeyPair(tlsConfig.CertPEM, tlsConfig.KeyPEM)
		if err != nil {
			return err
		}
		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	} else if len(tlsConfig.CertPEM) != 0 || len(tlsConfig.KeyPEM) != 0 {
		return fmt.Errorf("both client cert and client key must be provided")
	}

	if tlsConfig.CertFile != "" && tlsConfig.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
		if err != nil {
			return err
		}
		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	} else if tlsConfig.CertFile != "" || tlsConfig.KeyFile != "" {
		return fmt.Errorf("both client cert and client key must be provided")
	}

	return nil
}

// configureRootCertificates 配置 tls.Config 的 RootCAs
func configureRootCertificates(tlsClientConfig *tls.Config, tlsConfig *TLSConfig) error {
	if tlsConfig.CAFile != "" || tlsConfig.CAPath != "" || len(tlsConfig.CAPem) != 0 {
		rootConfig := &rootcerts.Config{
			CAFile:        tlsConfig.CAFile,
			CAPath:        tlsConfig.CAPath,
			CACertificate: tlsConfig.CAPem,
		}
		if err := rootcerts.ConfigureTLS(tlsClientConfig, rootConfig); err != nil {
			return err
		}
	}
	return nil
}
