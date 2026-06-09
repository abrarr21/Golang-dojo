package config

import (
	"log"
	"os"

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

type ImageKitConfig struct {
	ImgPrivateKey string
	ImgPublicKey  string
	ImgURL        string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	ImageKit ImageKitConfig
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".env not found ", err)
	}

	return &Config{
		ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},

		DatabaseConfig{
			MongoDB_URI: requiredEnv("MONGODB_URI"),
			DBName:      getEnv("DB_NAME", "test-dev"),
		},

		ImageKitConfig{
			ImgPrivateKey: optionalEnv("ImageKitPrivateKey"),
			ImgPublicKey:  optionalEnv("ImageKitPublicKey"),
			ImgURL:        optionalEnv("ImageKitURL"),
		},
	}

}

// critically imp (throws fatal error)
func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not defined in .env file", key)
	}
	return value
}

// optional (prints error but code runs)
func optionalEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("%s is not defined in .env file", key)
	}

	return value
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
