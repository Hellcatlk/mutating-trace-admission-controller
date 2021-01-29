package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/glog"
	yaml "gopkg.in/yaml.v2"
)

const defaultConfigPath string = "/etc/webhook/config/config.yaml"
const defaultCertPath string = "/etc/webhook/certs/cert.pem"
const defaultKeyPath string = "/etc/webhook/certs/key.pem"
const defaultSpanContextAnnotationKey string = "trace.kubernetes.io/span/context"

var config Config = Config{}

// Get return parsed configuration
func Get() Config {
	return config
}

// Parse reads YAML config into config struct, if path is "",
// use default config path("/etc/webhook/config/config.yaml").
func Parse(path string) (Config, error) {
	if path == "" {
		glog.Warningf("config path is empty, use default:: %v", defaultConfigPath)
		path = defaultConfigPath
	}

	// read config file
	configYaml, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return Config{}, fmt.Errorf("could not read YAML configuration file: %v", err)
	}

	// parse config yaml
	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		return Config{}, fmt.Errorf("could not umarshal YAML configuration file: %v", err)
	}

	// validate config
	if config.Trace.SampleRate < 0 || config.Trace.SampleRate > 1 {
		return Config{}, errors.New("sampling rate must be between 0 and 1 inclusive")
	}
	if config.Trace.SpanContextAnnotationKey == "" {
		glog.Warningf("span context annotationKey is empty, use default: %v", defaultSpanContextAnnotationKey)
		config.Trace.SpanContextAnnotationKey = defaultSpanContextAnnotationKey
	}
	if config.Certificate.CertPath == "" {
		glog.Warningf("cert path is empty, use default: %v", defaultCertPath)
		config.Certificate.CertPath = defaultCertPath
	}
	if config.Certificate.KeyPath == "" {
		glog.Warningf("key path is empty, use default: %v", defaultKeyPath)
		config.Certificate.KeyPath = defaultKeyPath
	}

	return config, nil
}

// LoadX509KeyPair reads and parses a public/private key pair from a pair of files.
// The files must contain PEM encoded data.
// The certificate file may contain intermediate certificates following the leaf certificate to form a certificate chain.
// On successful return, Certificate.Leaf will be nil because the parsed form of the certificate is not retained.
func (c *Config) LoadX509KeyPair() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(c.Certificate.CertPath, c.Certificate.KeyPath)
}
