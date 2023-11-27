package ginx

import "time"

type (
	// A HttpConfig is a http config.
	HttpConfig struct {
		Host         string        `json:"host" yaml:"host"`
		Port         int           `json:"port" yaml:"port"`
		Mode         string        `json:"mode" yaml:"mode"`
		ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
		WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	}
	// A PrivateKeyConf is a private key config.
	PrivateKeyConf struct {
		Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
		KeyFile     string `json:"key_file" yaml:"key_file"`
	}
	// A SignatureConf is a signature config.
	SignatureConf struct {
		Strict      bool             `json:"strict"`
		Expiry      time.Duration    `json:"expiry"`
		PrivateKeys []PrivateKeyConf `json:"private_keys"`
	}

	setConfigOpt func(*HttpConfig)
)

func DefaultHttpConfig() *HttpConfig {
	conf := &HttpConfig{
		Host:         "127.0.0.1",
		Port:         10101,
		Mode:         "release",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	return conf
}

func WithHost(host string) setConfigOpt {
	return func(conf *HttpConfig) {
		conf.Host = host
	}
}

func WithPort(port int) setConfigOpt {
	return func(conf *HttpConfig) {
		conf.Port = port
	}
}

func WithMode(mode string) setConfigOpt {
	return func(conf *HttpConfig) {
		conf.Mode = mode
	}
}

func WithReadTimeout(timeout time.Duration) setConfigOpt {
	return func(conf *HttpConfig) {
		conf.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) setConfigOpt {
	return func(conf *HttpConfig) {
		conf.WriteTimeout = timeout
	}
}
