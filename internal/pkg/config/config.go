package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// LoadConfigs loads a configuration into an object
func LoadConfigs(t string) (Config, error) {
	var f *os.File
	switch t {
	case "env":
		return loadFromENV()
	default:
		wdir, err := os.Getwd()
		configPath := filepath.Join(filepath.Clean(wdir), "../config.yaml")
		f, err := os.Open(configPath)
		if err != nil {
			return Config{}, err
		}
		defer f.Close()
	}
	return loadFromYAML(f)
}

func loadFromENV() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func loadFromYAML(r io.Reader) (Config, error) {
	var cfg Config
	err := yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return Config{}, nil
	}
	return cfg, nil
}

// Config set of configurations needed to run the app
type Config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Db struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     int    `yaml:"port" envconfig:"DB_PORT"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		Dbname   string `yaml:"dbname" envconfig:"DB_NAME"`
		Driver   string `yaml:"driver" envconfig:"DB_DRIVER"`
	} `yaml:"database"`
	Broker struct {
		ConnString string `yaml:"connString" envconfig:"BROKER_CONN_STRING"`
	} `yaml:"broker"`
	App struct {
	} `yaml:"app"`
}
