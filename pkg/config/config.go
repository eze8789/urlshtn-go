package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration support all the configuration to expose server and connect to different Backends
type Configuration struct {
	WebServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"webserver"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"postgres"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

// GenerateConfig map a configuration yaml file to the Configuration structure
func GenerateConfig(filepath string) (*Configuration, error) {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var cfg Configuration
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
