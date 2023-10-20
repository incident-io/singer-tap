package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

func LoadAndParse[T any](path string, obj T) (*T, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return &obj, errors.Wrapf(err, "unable to read file at path %v", path)
	}

	return ParseContents(b, obj)
}

func ParseContents[T any](content []byte, obj T) (*T, error) {
	if err := json.Unmarshal(content, &obj); err != nil {
		return &obj, err
	}
	return &obj, nil
}
