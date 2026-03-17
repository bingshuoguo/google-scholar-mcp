package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTransport      = "stdio"
	defaultBaseURL        = "https://scholar.google.com"
	defaultTimeout        = 15 * time.Second
	defaultMaxResults     = 10
	defaultRateLimitRPS   = 0.5
	defaultEnableAuthor   = true
	defaultLogLevel       = "info"
	defaultAcceptLanguage = "en-US,en;q=0.9"
	defaultUserAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
)

type Config struct {
	Transport      string
	BaseURL        string
	Timeout        time.Duration
	MaxResults     int
	RateLimitRPS   float64
	UserAgent      string
	AcceptLanguage string
	EnableAuthor   bool
	LogLevel       slog.Level
}

func Load() (Config, error) {
	cfg := Config{
		Transport:      getEnv("MCP_TRANSPORT", defaultTransport),
		BaseURL:        strings.TrimRight(getEnv("SCHOLAR_BASE_URL", defaultBaseURL), "/"),
		Timeout:        defaultTimeout,
		MaxResults:     defaultMaxResults,
		RateLimitRPS:   defaultRateLimitRPS,
		UserAgent:      getEnv("SCHOLAR_USER_AGENT", defaultUserAgent),
		AcceptLanguage: getEnv("SCHOLAR_ACCEPT_LANGUAGE", defaultAcceptLanguage),
		EnableAuthor:   defaultEnableAuthor,
		LogLevel:       slog.LevelInfo,
	}

	if timeout := os.Getenv("SCHOLAR_TIMEOUT"); timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err != nil {
			return Config{}, fmt.Errorf("parse SCHOLAR_TIMEOUT: %w", err)
		}
		cfg.Timeout = parsed
	}

	if maxResults := os.Getenv("SCHOLAR_MAX_RESULTS"); maxResults != "" {
		parsed, err := strconv.Atoi(maxResults)
		if err != nil {
			return Config{}, fmt.Errorf("parse SCHOLAR_MAX_RESULTS: %w", err)
		}
		cfg.MaxResults = parsed
	}

	if rps := os.Getenv("SCHOLAR_RATE_LIMIT_RPS"); rps != "" {
		parsed, err := strconv.ParseFloat(rps, 64)
		if err != nil {
			return Config{}, fmt.Errorf("parse SCHOLAR_RATE_LIMIT_RPS: %w", err)
		}
		cfg.RateLimitRPS = parsed
	}

	if enabled := os.Getenv("SCHOLAR_ENABLE_AUTHOR_TOOL"); enabled != "" {
		parsed, err := strconv.ParseBool(enabled)
		if err != nil {
			return Config{}, fmt.Errorf("parse SCHOLAR_ENABLE_AUTHOR_TOOL: %w", err)
		}
		cfg.EnableAuthor = parsed
	}

	level, err := parseLogLevel(getEnv("LOG_LEVEL", defaultLogLevel))
	if err != nil {
		return Config{}, err
	}
	cfg.LogLevel = level

	if cfg.Transport != "stdio" {
		return Config{}, fmt.Errorf("unsupported MCP_TRANSPORT %q: only stdio is implemented", cfg.Transport)
	}
	if cfg.MaxResults < 1 {
		return Config{}, fmt.Errorf("SCHOLAR_MAX_RESULTS must be >= 1")
	}
	if cfg.RateLimitRPS <= 0 {
		return Config{}, fmt.Errorf("SCHOLAR_RATE_LIMIT_RPS must be > 0")
	}

	return cfg, nil
}

func NewLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported LOG_LEVEL %q", value)
	}
}
