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
	TMDBAPIKey            string
	TMDBBaseURL           string
	TMDBImageBaseURL      string
	TMDBTimeout           time.Duration
	TMDBCacheTTL          time.Duration
	TMDBMaxRetries        int
	AIProvider            string
	AIModel               string
	AITimeout             time.Duration
	AIMaxInputChars       int
	AIMaxOutputTokens     int
	AIDailyUserBudget     int
	AICacheTTL            time.Duration
	OpusDataAPIKey        string
	OpusDataBaseURL       string
	OpusDataTimeout       time.Duration
	BoxOfficeDemoEnabled  bool
	CookieSecure          bool
	CookieSameSite        string
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
		TMDBAPIKey:         strings.TrimSpace(os.Getenv("TMDB_API_KEY")),
		TMDBBaseURL:        envOrDefault("TMDB_BASE_URL", "https://api.themoviedb.org/3"),
		TMDBImageBaseURL:   envOrDefault("TMDB_IMAGE_BASE_URL", "https://image.tmdb.org/t/p"),
		AIProvider:         strings.ToLower(envOrDefault("AI_PROVIDER", "offline")),
		AIModel:            envOrDefault("AI_MODEL", "gpt-4.1-mini"),
		OpusDataAPIKey:     strings.TrimSpace(os.Getenv("OPUSDATA_API_KEY")),
		OpusDataBaseURL:    strings.TrimRight(strings.TrimSpace(os.Getenv("OPUSDATA_BASE_URL")), "/"),
		CookieSameSite:     strings.ToLower(envOrDefault("COOKIE_SAME_SITE", "lax")),
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
	if cfg.TMDBTimeout, err = time.ParseDuration(envOrDefault("TMDB_TIMEOUT", "8s")); err != nil || cfg.TMDBTimeout <= 0 {
		return Config{}, errors.New("TMDB_TIMEOUT must be a positive Go duration such as 8s")
	}
	if cfg.TMDBCacheTTL, err = time.ParseDuration(envOrDefault("TMDB_CACHE_TTL", "10m")); err != nil || cfg.TMDBCacheTTL <= 0 {
		return Config{}, errors.New("TMDB_CACHE_TTL must be a positive Go duration such as 10m")
	}
	if cfg.TMDBMaxRetries, err = strconv.Atoi(envOrDefault("TMDB_MAX_RETRIES", "3")); err != nil || cfg.TMDBMaxRetries < 0 || cfg.TMDBMaxRetries > 5 {
		return Config{}, errors.New("TMDB_MAX_RETRIES must be an integer between 0 and 5")
	}
	if cfg.AITimeout, err = time.ParseDuration(envOrDefault("AI_TIMEOUT", "8s")); err != nil || cfg.AITimeout <= 0 {
		return Config{}, errors.New("AI_TIMEOUT must be a positive duration")
	}
	if cfg.AICacheTTL, err = time.ParseDuration(envOrDefault("AI_CACHE_TTL", "30m")); err != nil || cfg.AICacheTTL <= 0 {
		return Config{}, errors.New("AI_CACHE_TTL must be a positive duration")
	}
	if cfg.AIMaxInputChars, err = strconv.Atoi(envOrDefault("AI_MAX_INPUT_CHARS", "800")); err != nil || cfg.AIMaxInputChars < 80 || cfg.AIMaxInputChars > 4000 {
		return Config{}, errors.New("AI_MAX_INPUT_CHARS must be between 80 and 4000")
	}
	if cfg.AIMaxOutputTokens, err = strconv.Atoi(envOrDefault("AI_MAX_OUTPUT_TOKENS", "350")); err != nil || cfg.AIMaxOutputTokens < 50 || cfg.AIMaxOutputTokens > 2000 {
		return Config{}, errors.New("AI_MAX_OUTPUT_TOKENS must be between 50 and 2000")
	}
	if cfg.AIDailyUserBudget, err = strconv.Atoi(envOrDefault("AI_DAILY_USER_BUDGET", "20")); err != nil || cfg.AIDailyUserBudget < 1 || cfg.AIDailyUserBudget > 500 {
		return Config{}, errors.New("AI_DAILY_USER_BUDGET must be between 1 and 500")
	}
	if cfg.OpusDataTimeout, err = time.ParseDuration(envOrDefault("OPUSDATA_TIMEOUT", "8s")); err != nil || cfg.OpusDataTimeout <= 0 {
		return Config{}, errors.New("OPUSDATA_TIMEOUT must be a positive duration")
	}
	cfg.BoxOfficeDemoEnabled, err = strconv.ParseBool(envOrDefault("BOX_OFFICE_DEMO_ENABLED", "true"))
	if err != nil {
		return Config{}, errors.New("BOX_OFFICE_DEMO_ENABLED must be true or false")
	}
	cfg.CookieSecure, err = strconv.ParseBool(envOrDefault("COOKIE_SECURE", "false"))
	if err != nil {
		return Config{}, errors.New("COOKIE_SECURE must be true or false")
	}

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
	if cfg.TMDBAPIKey == "" {
		problems = append(problems, "TMDB_API_KEY is required")
	}
	if cfg.OpenAIAPIKey != "" && !strings.Contains(cfg.BasePromptTemplate, "{rankings}") {
		problems = append(problems, "BASE_PROMPT_TEMPLATE must contain {rankings} when OPENAI_API_KEY is configured")
	}
	if cfg.AIProvider != "offline" && cfg.AIProvider != "openai" {
		problems = append(problems, "AI_PROVIDER must be offline or openai")
	}
	if cfg.AIProvider == "openai" && cfg.OpenAIAPIKey == "" {
		problems = append(problems, "OPENAI_API_KEY is required when AI_PROVIDER=openai")
	}
	if (cfg.OpusDataAPIKey == "") != (cfg.OpusDataBaseURL == "") {
		problems = append(problems, "OPUSDATA_API_KEY and OPUSDATA_BASE_URL must be configured together")
	}
	if cfg.CookieSameSite != "lax" && cfg.CookieSameSite != "strict" && cfg.CookieSameSite != "none" {
		problems = append(problems, "COOKIE_SAME_SITE must be lax, strict, or none")
	}
	if cfg.CookieSameSite == "none" && !cfg.CookieSecure {
		problems = append(problems, "COOKIE_SECURE must be true when COOKIE_SAME_SITE=none")
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
