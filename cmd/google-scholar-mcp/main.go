package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
	"github.com/bingshuoguo/google-scholar-mcp/internal/mcpserver"
)

var version = "dev"
var readBuildInfo = debug.ReadBuildInfo

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	command, err := parseCommand(args)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n\n%s", err, usageText())
		return 2
	}

	switch command {
	case commandHelp:
		fmt.Fprint(stdout, usageText())
		return 0
	case commandVersion:
		fmt.Fprintln(stdout, resolveVersion(version))
		return 0
	case commandStdio:
		if err := runStdio(resolveVersion(version)); err != nil {
			slog.New(slog.NewTextHandler(stderr, nil)).Error("server exited", "error", err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "unsupported command %q\n\n%s", command, usageText())
		return 2
	}
}

func runStdio(serverVersion string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := config.NewLogger(cfg.LogLevel)
	server := mcpserver.New(cfg, logger, serverVersion)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("run stdio server: %w", err)
	}
	return nil
}

const (
	commandHelp    = "help"
	commandStdio   = "stdio"
	commandVersion = "version"
)

func parseCommand(args []string) (string, error) {
	if len(args) == 0 {
		return commandStdio, nil
	}

	switch args[0] {
	case "stdio":
		if len(args) != 1 {
			return "", fmt.Errorf("stdio does not accept additional arguments")
		}
		return commandStdio, nil
	case "--version", "-version", "version":
		if len(args) != 1 {
			return "", fmt.Errorf("version does not accept additional arguments")
		}
		return commandVersion, nil
	case "--help", "-h", "help":
		if len(args) != 1 {
			return "", fmt.Errorf("help does not accept additional arguments")
		}
		return commandHelp, nil
	default:
		return "", fmt.Errorf("unknown command %q", args[0])
	}
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

func usageText() string {
	return `Google Scholar MCP

Usage:
  google-scholar-mcp             Start the MCP server over stdio
  google-scholar-mcp stdio       Start the MCP server over stdio
  google-scholar-mcp version     Print the server version
  google-scholar-mcp --version   Print the server version
  google-scholar-mcp help        Show this help text

Configuration:
  The server reads runtime options from environment variables such as:
  MCP_TRANSPORT, SCHOLAR_TIMEOUT, SCHOLAR_MAX_RESULTS, SCHOLAR_RATE_LIMIT_RPS,
  SCHOLAR_ENABLE_AUTHOR_TOOL, LOG_LEVEL.
`
}
