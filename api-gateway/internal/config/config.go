package config

import "os"

type Config struct {
	AppEnv        string
	Port          string
	GRPCServerURL string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		GRPCServerURL: getEnv("GRPC_SERVER_URL", "localhost:50051"),
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
