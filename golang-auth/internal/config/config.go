package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	MongoDB_URI string
	DBName      string
}

type JWTConfig struct {
	JWT_SECRET     string
	AccessTokenTTL time.Duration
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		log.Fatal("MongoDB_URI is not defined in .env file")
	}

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT_SECRET is not defined in .env file")
	}

	return &Config{
		ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},

		DatabaseConfig{
			MongoDB_URI: mongodb_uri,
			DBName:      getEnv("DB_NAME", "auth-dev"),
		},

		JWTConfig{
			JWT_SECRET:     jwt_secret,
			AccessTokenTTL: mustParseDuration(getEnv("ACCESS_TOKEN_TTL", "20m")),
		},
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("invalid duration %q: expected foramt 15m, 1h, 7d", s)
	}

	return d
}
