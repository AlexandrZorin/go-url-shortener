package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfig(t *testing.T) {
	type want struct {
		serverAddress string
		baseURL       string
	}

	tests := []struct {
		name             string
		envServerAddress string
		envBaseURL       string
		want             want
	}{
		{
			name:             "Environment variables are filled",
			envServerAddress: "localhost:9090",
			envBaseURL:       "http://localhost:9090",
			want: want{
				serverAddress: "localhost:9090",
				baseURL:       "http://localhost:9090",
			},
		},
		{
			name:             "Environment variables are not filled",
			envServerAddress: "",
			envBaseURL:       "",
			want: want{
				serverAddress: "localhost:8080",
				baseURL:       "http://localhost:8080",
			},
		},
		{
			name:             "Only server address environment variable is set",
			envServerAddress: "localhost:9090",
			envBaseURL:       "",
			want: want{
				serverAddress: "localhost:9090",
				baseURL:       "http://localhost:8080",
			},
		},
		{
			name:             "Only base URL environment variable is set",
			envServerAddress: "",
			envBaseURL:       "http://localhost:9090",
			want: want{
				serverAddress: "localhost:8080",
				baseURL:       "http://localhost:9090",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldServerAddress := os.Getenv("SERVER_ADDRESS")
			oldBaseURL := os.Getenv("BASE_URL")

			if tt.envServerAddress != "" {
				os.Setenv("SERVER_ADDRESS", tt.envServerAddress)
			} else {
				os.Unsetenv("SERVER_ADDRESS")
			}

			if tt.envBaseURL != "" {
				os.Setenv("BASE_URL", tt.envBaseURL)
			} else {
				os.Unsetenv("BASE_URL")
			}

			cfg := CreateConfig()

			assert.Equal(t, tt.want.serverAddress, cfg.ServerAddress)
			assert.Equal(t, tt.want.baseURL, cfg.URL)

			if oldServerAddress != "" {
				os.Setenv("SERVER_ADDRESS", oldServerAddress)
			} else {
				os.Unsetenv("SERVER_ADDRESS")
			}

			if oldBaseURL != "" {
				os.Setenv("BASE_URL", oldBaseURL)
			} else {
				os.Unsetenv("BASE_URL")
			}
		})
	}
}
