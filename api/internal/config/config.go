package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GoogleClientID  string
	JWTSecret       string
	DatabaseURL     string
}

func Load() *Config {
	
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found (using environment variables)")
	}

	return &Config{
		Port:           mustGet("PORT"),
		GoogleClientID: mustGet("GOOGLE_CLIENT_ID"),
		JWTSecret:      mustGet("JWT_SECRET"),
		DatabaseURL:    mustGet("DATABASE_URL"),
	}
}

func mustGet(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing required env: %s", key)
	}
	return v
}