package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	APIHost  string `mapstructure:"API_HOST"`
	APIPort  string `mapstructure:"API_PORT"`
	MongoURI string `mapstructure:"API_MONGO_URI" validate:"required"`
	MongoDB  string `mapstructure:"API_MONGO_DB_NAME" validate:"required"`
}

func New() (*Config, error) {
	v := viper.New()

	if err := bindEnvs(v, Config{}); err != nil {
		return nil, fmt.Errorf("failed to bind environment variables: %w", err)
	}

	v.SetDefault("API_HOST", "0.0.0.0")
	v.SetDefault("API_PORT", "8000")

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into map: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("missing required attributes: %w", err)
	}

	return &config, nil
}

// // bindEnv is a workaround for a known issue in viper
// // The issue means that env variables cannot be read unless there is a blank config file or every config value set to have a default
// // Issue ref: https://github.com/spf13/viper/issues/761
func bindEnvs(v *viper.Viper, config Config) error {
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(config, &envKeysMap); err != nil {
		return fmt.Errorf("failed to determine keys %w", err)
	}
	for k := range *envKeysMap {
		if bindErr := v.BindEnv(k); bindErr != nil {
			return fmt.Errorf("failed to bind env variables. %w", bindErr)
		}
	}
	return nil
}
