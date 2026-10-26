package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI            string
	DatabaseName          string
	AllowedOrigins        []string
	SecretKey             string
	SecretRefreshKey      string
	OpenAIAPIKey          string
	BasePromptTemplate    string
	RecommendedMovieLimit int64
	Port                  string
	ShutdownTimeout       time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load(".env")
	cfg := Config{
		MongoDBURI:         strings.TrimSpace(os.Getenv("MONGODB_URI")),
		DatabaseName:       strings.TrimSpace(os.Getenv("DATABASE_NAME")),
		SecretKey:          strings.TrimSpace(os.Getenv("SECRET_KEY")),
		SecretRefreshKey:   strings.TrimSpace(os.Getenv("SECRET_REFRESH_KEY")),
		OpenAIAPIKey:       strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
		BasePromptTemplate: os.Getenv("BASE_PROMPT_TEMPLATE"),
		Port:               envOrDefault("PORT", "8080"),
		ShutdownTimeout:    10 * time.Second,
	}
	for _, origin := range strings.Split(envOrDefault("ALLOWED_ORIGINS", "http://localhost:5173"), ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
		}
	}
	limit, err := strconv.ParseInt(envOrDefault("RECOMMENDED_MOVIE_LIMIT", "5"), 10, 64)
	if err != nil || limit < 1 || limit > 100 {
		return Config{}, errors.New("RECOMMENDED_MOVIE_LIMIT must be an integer between 1 and 100")
	}
	cfg.RecommendedMovieLimit = limit

	var problems []string
	if cfg.MongoDBURI == "" {
		problems = append(problems, "MONGODB_URI is required")
	}
	if cfg.DatabaseName == "" {
		problems = append(problems, "DATABASE_NAME is required")
	}
	if len(cfg.SecretKey) < 32 {
		problems = append(problems, "SECRET_KEY must contain at least 32 characters")
	}
	if len(cfg.SecretRefreshKey) < 32 {
		problems = append(problems, "SECRET_REFRESH_KEY must contain at least 32 characters")
	}
	if cfg.SecretKey != "" && cfg.SecretKey == cfg.SecretRefreshKey {
		problems = append(problems, "SECRET_KEY and SECRET_REFRESH_KEY must be different")
	}
	if len(cfg.AllowedOrigins) == 0 {
		problems = append(problems, "ALLOWED_ORIGINS must contain at least one origin")
	}
	if cfg.OpenAIAPIKey != "" && !strings.Contains(cfg.BasePromptTemplate, "{rankings}") {
		problems = append(problems, "BASE_PROMPT_TEMPLATE must contain {rankings} when OPENAI_API_KEY is configured")
	}
	if len(problems) > 0 {
		return Config{}, fmt.Errorf("invalid environment configuration: %s", strings.Join(problems, "; "))
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
