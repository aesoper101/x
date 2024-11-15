package cert

// TLSConfig is used to generate a TLSClientConfig that's useful for talking to
// Consul using TLS.
type TLSConfig struct {
	// Address is the optional address of the Consul server. The port, if any
	// will be removed from here and this will be set to the ServerName of the
	// resulting config.
	Address string `json:"address" mapstructure:"address"`

	// CAFile is the optional path to the CA certificate used for Consul
	// communication, defaults to the system bundle if not specified.
	CAFile string `json:"ca_file" mapstructure:"ca_file"`

	// CAPath is the optional path to a directory of CA certificates to use for
	// Consul communication, defaults to the system bundle if not specified.
	CAPath string `json:"ca_path" mapstructure:"ca_path"`

	// CAPem is the optional PEM-encoded CA certificate used for Consul
	// communication, defaults to the system bundle if not specified.
	CAPem []byte `json:"ca_pem" mapstructure:"ca_pem"`

	// CertFile is the optional path to the certificate for Consul
	// communication. If this is set then you need to also set KeyFile.
	CertFile string `json:"cert_file" mapstructure:"cert_file"`

	// CertPEM is the optional PEM-encoded certificate for Consul
	// communication. If this is set then you need to also set KeyPEM.
	CertPEM []byte `json:"cert_pem" mapstructure:"cert_pem"`

	// KeyFile is the optional path to the private key for Consul communication.
	// If this is set then you need to also set CertFile.
	KeyFile string `json:"key_file" mapstructure:"key_file"`

	// KeyPEM is the optional PEM-encoded private key for Consul communication.
	// If this is set then you need to also set CertPEM.
	KeyPEM []byte `json:"key_pem" mapstructure:"key_pem"`

	// InsecureSkipVerify if set to true will disable TLS host verification.
	InsecureSkipVerify bool `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify"`
}
