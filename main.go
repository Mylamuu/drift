package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Mylamuu/drift/internal/storage"
)

func initLogger(logLevel string) {
	var level slog.Level

	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}

func main() {
	config, err := loadConfig(os.Args[1:])
	if err != nil {
		slog.Error("Unable to load config", "error", err)
		os.Exit(1)
	}

	initLogger(config.LogLevel)
	slog.Debug("Configuration has been loaded", "config", fmt.Sprintf("%+v", config))

	if err := storage.Init(config.StoragePath); err != nil {
		slog.Error("Failed to initialize storage.", "error", err)
		os.Exit(1)
	}
}
