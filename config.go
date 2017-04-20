package exploregithub

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	PoolSize int `yaml:"pool_size"`
}

type OAuthConfig struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type ServerConfig struct {
	Bind string `yaml:"bind"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	OAuth    OAuthConfig    `yaml:"oauth"`
}

func NewConfig(path string) (*Config, error) {
	var config Config

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("failed to read configuration file")
	}
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, errors.New("failed to unmarshal file")
	}

	return &config, nil
}
