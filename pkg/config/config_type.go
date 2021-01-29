package config

// Config represents all config we need to initialize the webhook server
type Config struct {
	Certificate Certificate `yaml:"certificate"`
	Trace       Trace       `yaml:"trace"`
}

// Certificate is the configuration for the certificate
type Certificate struct {
	CertPath string `yaml:"certPath"`
	KeyPath  string `yaml:"keyPath"`
}

// Trace is the configuration for the trace context added to pods
type Trace struct {
	SampleRate               float64 `yaml:"sampleRate"`
	SpanContextAnnotationKey string  `yaml:"spanContextAnnotationKey"`
}
