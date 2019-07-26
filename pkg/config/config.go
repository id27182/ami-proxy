package config

import (
	"fmt"
	"github.com/id27182/ami-proxy/pkg/env"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

const configurationFile = "ami-proxy"
var config *viper.Viper

type Config struct {
	config *viper.Viper
}

func initConfig() error {
	v := viper.New()

	// configure viper to read configuration from environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("AMIPROXY")
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	confDir, err := env.GetExecutableDir()
	if err != nil {
		return fmt.Errorf("unable to determine configuration directory. Original error: %s", err)
	}

	v.SetConfigName(configurationFile)
	v.AddConfigPath(confDir)

	v.ReadInConfig()
	config = v
	return nil
}

func GetConfig() (*Config, error)  {
	if config == nil {
		err := initConfig()
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		config: config,
	}, nil
}

// ProxyConfig describes Service Proxy configuration for ECD
type ProxyConfig struct {
	BindPort string `validate:"required"`

	DestHost     string `validate:"required"`
	DestPort     string `validate:"required"`
	DestProtocol string `validate:"required"`
	DestResource string `validate:"required"`
}

func (c *Config) Proxy() (*ProxyConfig, error)  {
	res := ProxyConfig{
		BindPort: c.config.GetString("proxy.bind.port"),

		DestHost:      c.config.GetString("proxy.dest.host"),
		DestPort:      c.config.GetString("proxy.dest.port"),
		DestProtocol:  c.config.GetString("proxy.dest.protocol"),
		DestResource:  c.config.GetString("proxy.dest.resource"),
	}

	err := validator.New().Struct(res)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive configuration. Original error: %s", FormatValidationError(reflect.TypeOf(res), err, ""))
	}

	return &res, nil
}