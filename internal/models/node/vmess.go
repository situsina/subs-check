package node

type Vmess struct {
	Info   Info
	Config VmessConfig
}

type VmessConfig struct {
	BaseConfig          `yaml:",inline"`
	Uuid                string       `yaml:"uuid"`
	AlterId             int          `yaml:"alterId"`
	Cipher              string       `yaml:"cipher"`
	PacketEncoding      string       `yaml:"packet-encoding"`
	GlobalPadding       bool         `yaml:"global-padding"`
	AuthenticatedLength bool         `yaml:"authenticated-length"`
	TLS                 bool         `yaml:"tls"`
	Servername          string       `yaml:"servername"`
	Alpn                []string     `yaml:"alpn"`
	Fingerprint         string       `yaml:"fingerprint"`
	ClientFingerprint   string       `yaml:"client-fingerprint"`
	SkipCertVerify      bool         `yaml:"skip-cert-verify"`
	RealityOpts         *RealityOpts `yaml:"reality-opts"`
	Network             string       `yaml:"network"`
	Smux                *SmuxOpts    `yaml:"smux"`
}
