package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func initLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: verbose,
		}),
	)
}
