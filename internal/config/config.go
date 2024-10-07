package config

import (
	"os"

	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
)

type Config struct {
	Env          string `yaml:"env"`
	DBConnection `yaml:"db_connection"`
	Speller      string `yaml:"speller"`
}

type DBConnection struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-required:"true"`
}

func ParseConfig() (*Config, error) {
	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, err
	}
	replaced := os.ExpandEnv(string(data))
	cfg := &Config{}
	err = yaml.Unmarshal([]byte(replaced), cfg)
	return cfg, err
}
