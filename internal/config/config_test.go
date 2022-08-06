package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	envVars := map[string]string{
		"API_MONGO_URI":     "mongodb://localhost:27017",
		"API_MONGO_DB_NAME": "test",
		"API_HOST":          "0.0.0.0",
		"API_PORT":          "8000",
	}

	tests := []struct {
		name        string
		envVars     map[string]string
		expected    *Config
		expectedErr string
	}{
		{
			name:    "successfully loads config",
			envVars: envVars,
			expected: &Config{
				MongoURI: "mongodb://localhost:27017",
				MongoDB:  "test",
				APIHost:  "0.0.0.0",
				APIPort:  "8000",
			},
		},
		{
			name:        "Errors when an environment var is missing",
			envVars:     map[string]string{},
			expectedErr: "Field validation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				err := os.Setenv(k, fmt.Sprintf("%v", v))
				require.NoError(t, err)
			}

			got, err := New()

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			os.Clearenv()
		})
	}
}
