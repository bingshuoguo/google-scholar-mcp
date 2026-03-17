package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
	"github.com/bingshuoguo/google-scholar-mcp/internal/mcpserver"
)

var version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("load config", "error", err)
		os.Exit(1)
	}

	logger := config.NewLogger(cfg.LogLevel)
	server := mcpserver.New(cfg, logger, version)

	if err := server.Run(context.Background()); err != nil {
		logger.Error("server exited", "error", err)
		os.Exit(1)
	}
}
