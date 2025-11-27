package config

import "github.com/pachirode/docker-demo/internal/docker-demo/options"

type Config struct {
	*options.RunOptions
}

func CreateConfigFromOptions(opts *options.RunOptions) (*Config, error) {
	return &Config{opts}, nil
}
