package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port         string
	JWTSecret    string
	JWTExpiresIn int64

	// Database
	DatabaseDSN string // 优先使用完整连接字符串
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	// OpenAI
	OpenAIKey string
	BaseURL   string
	Model     string

	// Rate Limit
	RateLimitTTL   int64
	RateLimitLimit int
}

func Load() (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiresIn: getEnvAsInt64("JWT_EXPIRES_IN", 86400), // 24小时

		DatabaseDSN: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "ai_chat"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),

		// 支持.env文件中AI_*格式的变量名
		OpenAIKey: getEnv("AI_API_KEY", getEnv("OPENAI_API_KEY", "")),
		BaseURL:   getEnv("AI_BASE_URL", getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1")),
		Model:     getEnv("AI_MODEL", getEnv("OPENAI_MODEL", "gpt-3.5-turbo")),

		RateLimitTTL:   getEnvAsInt64("RATE_LIMIT_TTL", 60),
		RateLimitLimit: getEnvAsInt("RATE_LIMIT_LIMIT", 60),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(name string, defaultValue int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) DatabaseURL() string {
	// 优先使用完整连接字符串
	if c.DatabaseDSN != "" {
		return c.DatabaseDSN
	}
	return "host=" + c.DBHost +
		" port=" + strconv.Itoa(c.DBPort) +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}
