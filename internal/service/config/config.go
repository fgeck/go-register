package main

// import (
// 	"log"
// 	"os"
// 	"strconv"
// )

// type Config struct {
// 	Host          string
// 	Port          string
// 	DatabaseURL   string
// 	SessionSecret string
// 	Env           string
// 	HTTPS         bool
// 	CSRFKey       string
// 	Argon2Config  Argon2Config
// }

// func loadConfig() *Config {
// 	return &Config{
// 		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
// 		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
// 		SessionSecret: getEnv("SESSION_SECRET", "super-secret-key-32-chars-long"),
// 		Env:           getEnv("ENV", "development"),
// 		HTTPS:         getEnvAsBool("HTTPS", false),
// 		CSRFKey:       getEnv("CSRF_KEY", "csrf-super-secret-key-32-chars"),

// 		Argon2Config: Argon2Config{
// 			Memory:      getEnvAsUint32("ARGON2_MEMORY", 64*1024), // 64MB
// 			Iterations:  getEnvAsUint32("ARGON2_ITERATIONS", 3),
// 			Parallelism: getEnvAsUint8("ARGON2_PARALLELISM", 2),
// 			SaltLength:  getEnvAsUint32("ARGON2_SALT_LENGTH", 16),
// 			KeyLength:   getEnvAsUint32("ARGON2_KEY_LENGTH", 32),
// 		},
// 	}
// }

// func getEnv(key, fallback string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return fallback
// }

// func getEnvAsBool(key string, fallback bool) bool {
// 	if value, exists := os.LookupEnv(key); exists {
// 		b, err := strconv.ParseBool(value)
// 		if err != nil {
// 			log.Printf("Invalid bool value for %s: %v", key, err)
// 			return fallback
// 		}
// 		return b
// 	}
// 	return fallback
// }

// func getEnvAsUint32(key string, fallback uint32) uint32 {
// 	if value, exists := os.LookupEnv(key); exists {
// 		i, err := strconv.ParseUint(value, 10, 32)
// 		if err != nil {
// 			log.Printf("Invalid uint32 value for %s: %v", key, err)
// 			return fallback
// 		}
// 		return uint32(i)
// 	}
// 	return fallback
// }

// func getEnvAsUint8(key string, fallback uint8) uint8 {
// 	if value, exists := os.LookupEnv(key); exists {
// 		i, err := strconv.ParseUint(value, 10, 8)
// 		if err != nil {
// 			log.Printf("Invalid uint8 value for %s: %v", key, err)
// 			return fallback
// 		}
// 		return uint8(i)
// 	}
// 	return fallback
// }
