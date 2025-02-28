package config

import "os"

type Config struct {
	AppEnv string
	Port   string
	DBUrl  string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "50051"),
		DBUrl:  getEnv("DB_URL", "postgres://user_name:secret@localhost:15432/todos?sslmode=disable"),
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
