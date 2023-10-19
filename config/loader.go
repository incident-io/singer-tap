package config

import (
	"context"
	"os"
)

type Loader interface {
	Load(context.Context) (*Config, error)
}

type LoaderFunc func(context.Context) (*Config, error)

func (l LoaderFunc) Load(ctx context.Context) (*Config, error) {
	return l(ctx)
}

// FileLoader loads config from a filepath
type FileLoader string

func (l FileLoader) Load(context.Context) (*Config, error) {
	data, err := os.ReadFile(string(l))
	if err != nil {
		return nil, err
	}

	return Parse(string(l), data)
}
