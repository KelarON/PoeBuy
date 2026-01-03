package config

import (
	"fmt"
	"os"
	"poebuy/utils"

	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

const DEFAULT_VISIT_DELAY = 20

var ErrorNoConfigFile = errors.New("Config file not found")

// Config is the main configuration struct
type Config struct {
	General General `yaml:"general"`
	Trade   Trade   `yaml:"trade"`
	errChan chan error
}

type General struct {
	Poesessid string `yaml:"poesessid"`
}

type Trade struct {
	League     string `yaml:"league"`
	VisitDelay int    `yaml:"visit_delay"`
	Links      []Link `yaml:"links"`
}

type Link struct {
	Name    string `yaml:"name"`
	Code    string `yaml:"code"`
	Delay   int64  `yaml:"delay"`
	IsActiv bool   `yaml:"-"`
}

// LoadConfig loads the config from config file
func LoadConfig() (*Config, error) {

	cfg := &Config{}

	if !configFileExists() {
		return cfg, ErrorNoConfigFile
	}

	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	cfg.General.Poesessid, _ = utils.Decrypt(cfg.General.Poesessid)

	// Set default values
	if cfg.Trade.VisitDelay == 0 {
		cfg.Trade.VisitDelay = DEFAULT_VISIT_DELAY
	}

	return cfg, nil
}

// Save saves the config to config file
func (cfg *Config) Save() {
	poe := cfg.General.Poesessid
	encPoe, err := utils.Encrypt(poe)
	if err != nil && cfg.errChan != nil {
		cfg.errChan <- err
	}
	cfg.General.Poesessid = encPoe

	err = utils.WriteStructToYAMLFile("config.yaml", cfg)
	if err != nil && cfg.errChan != nil {
		cfg.errChan <- err
	}
	cfg.General.Poesessid = poe
}

func (cfg *Config) DefineErrorChannel(errChan chan error) {
	cfg.errChan = errChan
}

// configFileExists checks if config file exists
func configFileExists() bool {
	_, err := os.Stat("config.yaml")
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
