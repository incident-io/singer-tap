package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	APIKey   string `json:"api_key,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.APIKey, validation.Required.
			Error("must provide an api_key to authenticate against the incident.io API.")),
	)
}
