package main

import (
	"log/slog"
	"os"
)

var logger = slog.New(
	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
)
