package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

func ImportConfig[T any](path string, defaultParsers func() *T, useEnv bool) (Config[T], error) {
	c := DefaultConfig[T](defaultParsers)

	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return Config[T]{}, fmt.Errorf("open config file: %w", err)
		}

		defer f.Close()

		err = yaml.NewDecoder(f).Decode(&c)
		if err != nil {
			return Config[T]{}, fmt.Errorf("decode yaml: %w", err)
		}
	}

	if useEnv {
		// Важно: не понятно баг или фича, но nil значения инициализируются даже если нет окружения.
		err := envconfig.Process("APP", &c)
		if err != nil {
			return Config[T]{}, fmt.Errorf("decode env: %w", err)
		}
	}

	return c, nil
}
