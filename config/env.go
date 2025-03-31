package config

import (
	"fmt"
	"log"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func GetConfig() (*Cfg, error) {
	conf, err := initialize()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func initialize() (*Cfg, error) {
	if err := load(); err != nil {
		return nil, fmt.Errorf("can't load config: %w", err)
	}

	cfg := &Cfg{}
	opts := env.Options{
		Prefix:                "APP_",
		UseFieldNameByDefault: true,
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}

	return cfg, nil
}

func load() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
		return err
	}

	return nil
}
