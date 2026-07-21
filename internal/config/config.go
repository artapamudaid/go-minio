package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	ServerSecretKey string
	Port            string
}

func LoadConfig() *Config {
	// Membaca file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca langsung dari OS environment")
	}

	return &Config{
		Endpoint:        getEnv("NEO_ENDPOINT", "127.0.0.1:9000"),
		AccessKeyID:     getEnv("NEO_ACCESS_KEY", "ROOTUSER"),
		SecretAccessKey: getEnv("NEO_SECRET_KEY", "CHANGEME123"),
		UseSSL:          getEnv("NEO_USE_SSL", "false") == "true",
		BucketName:      getEnv("NEO_BUCKET", "my-bucket"),
		ServerSecretKey: getEnv("SERVER_SECRET_KEY", "my-secret-key"),
		Port:            getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
