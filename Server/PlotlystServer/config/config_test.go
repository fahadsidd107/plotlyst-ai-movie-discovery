package config

import "testing"

func setValidEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017")
	t.Setenv("DATABASE_NAME", "plotlyst")
	t.Setenv("SECRET_KEY", "12345678901234567890123456789012")
	t.Setenv("SECRET_REFRESH_KEY", "abcdefghijklmnopqrstuvwxyz123456")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
	t.Setenv("RECOMMENDED_MOVIE_LIMIT", "5")
	t.Setenv("OPENAI_API_KEY", "")
}

func TestLoadAcceptsValidEnvironment(t *testing.T) {
	setValidEnvironment(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}
	if cfg.DatabaseName != "plotlyst" || cfg.RecommendedMovieLimit != 5 {
		t.Fatalf("Load() returned unexpected configuration: %+v", cfg)
	}
}

func TestLoadRejectsShortSecrets(t *testing.T) {
	setValidEnvironment(t)
	t.Setenv("SECRET_KEY", "too-short")
	if _, err := Load(); err == nil {
		t.Fatal("Load() accepted an unsafe SECRET_KEY")
	}
}

func TestLoadRejectsInvalidRecommendationLimit(t *testing.T) {
	setValidEnvironment(t)
	t.Setenv("RECOMMENDED_MOVIE_LIMIT", "zero")
	if _, err := Load(); err == nil {
		t.Fatal("Load() accepted an invalid recommendation limit")
	}
}
