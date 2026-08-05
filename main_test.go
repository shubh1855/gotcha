package main

import (
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	os.Exit(m.Run())
}
