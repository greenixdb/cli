package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ProjectName string `mapstructure:"project_name"`
	Version     string `mapstructure:"version"`
	BuildCode   string `mapstructure:"build_code"`
	OutputDir   string `mapstructure:"output_dir"`
	Verbose     bool   `mapstructure:"verbose"`
}

func Load() (*Config, error) {
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.OutputDir == "" {
		config.OutputDir = "build-output"
	}

	return &config, nil
}

