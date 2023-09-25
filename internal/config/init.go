package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	CONFIGPATH = "configs/config.yml"
)

type Server struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

type Postgres struct {
	DBName   string `yaml:"dbname"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
}

type Redis struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	DB       int    `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

func InitConfig() (*Config, error) {
	f, err := os.ReadFile(CONFIGPATH)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	cfg := &Config{}
	yaml.Unmarshal(f, cfg)
	return cfg, nil
}
