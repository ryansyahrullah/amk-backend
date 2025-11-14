package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBLogMode string

	JWTAccessSecret  string
	JWTRefreshSecret string

	JWTAccessMinutes int
	JWTRefreshDays   int
}

func LoadConfig() *Config {
	// Coba load .env kalau kamu jalan dari folder auth-amk
	_ = godotenv.Load(".env")

	// Coba juga load .env kalau kamu jalan dari root amk-backend
	_ = godotenv.Load("auth-amk/.env")

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBHost:    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", "auth_db"),
		DBLogMode: getEnv("DB_LOG_MODE", "info"),

		JWTAccessSecret:  getEnv("JWT_AKSES_SECRET", "dev-access-secret"),
		JWTRefreshSecret: getEnv("JWT_PENYEGAR_SECRET", "dev-refresh-secret"),
	}

	cfg.JWTAccessMinutes = getEnvAsInt("JWT_AKSES_DURASI_MENIT", 30)
	cfg.JWTRefreshDays = getEnvAsInt("JWT_PENYEGAR_DURASI_HARI", 7)

	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return def
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("warning: gagal parse env %s='%s' sebagai int, pakai default %d\n", key, valStr, def)
		return def
	}
	return val
}
