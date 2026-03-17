package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
	"github.com/bingshuoguo/google-scholar-mcp/internal/mcpserver"
)

var version = "dev"
var readBuildInfo = debug.ReadBuildInfo

func main() {
	if wantsVersion(os.Args[1:]) {
		fmt.Println(resolveVersion(version))
		return
	}

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

func wantsVersion(args []string) bool {
	return len(args) == 1 && (args[0] == "--version" || args[0] == "-version" || args[0] == "version")
}

func resolveVersion(binaryVersion string) string {
	if binaryVersion != "" && binaryVersion != "dev" {
		return binaryVersion
	}

	info, ok := readBuildInfo()
	if !ok {
		return "dev"
	}

	switch info.Main.Version {
	case "", "(devel)":
		return "dev"
	default:
		return info.Main.Version
	}
}
