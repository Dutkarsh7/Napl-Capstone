package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	UploadDir            string
	OCRServiceURL        string
	ValidationServiceURL string
	AllowedOrigins       []string
	WorkerPoolSize       int
	WorkerQueueCap       int
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smart_pharmacy?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "smart-pharmacy-super-secret-capstone-jwt-key-2026")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")
	ocrURL := getEnv("OCR_SERVICE_URL", "http://localhost:8000")
	valURL := getEnv("VALIDATION_SERVICE_URL", "http://localhost:5005")

	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000,http://127.0.0.1:5173,http://localhost:8080")
	var origins []string
	for _, s := range strings.Split(originsStr, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			origins = append(origins, s)
		}
	}

	poolSize := getEnvInt("WORKER_POOL_SIZE", 4)
	queueCap := getEnvInt("WORKER_QUEUE_CAP", 100)

	return &Config{
		Port:                 port,
		DatabaseURL:          dbURL,
		JWTSecret:            jwtSecret,
		UploadDir:            uploadDir,
		OCRServiceURL:        ocrURL,
		ValidationServiceURL: valURL,
		AllowedOrigins:       origins,
		WorkerPoolSize:       poolSize,
		WorkerQueueCap:       queueCap,
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && strings.TrimSpace(val) != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		if n, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return n
		}
	}
	return defaultVal
}
