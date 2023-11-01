package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
)

type customGinLogWriter struct {
}

var _ io.Writer = customGinLogWriter{}

func (customGinLogWriter) Write(b []byte) (int, error) {
	level, b := parseLevelAndMessage(b)
	slog.Log(context.Background(), level, fmt.Sprintf("[GIN] %s", b))
	return len(b), nil
}

func parseLevelAndMessage(b []byte) (slog.Level, []byte) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	b = bytes.ReplaceAll(b, []byte("\n"), []byte("\n\t"))

	level := slog.LevelInfo
	if trimmed := bytes.TrimPrefix(b, []byte("[GIN-debug] ")); len(trimmed) != len(b) {
		level = slog.LevelDebug
		b = trimmed
	} else if trimmed := bytes.TrimPrefix(b, []byte("[GIN] ")); len(trimmed) != len(b) {
		b = trimmed
	}

	if trimmed := bytes.TrimPrefix(b, []byte("[WARNING] ")); len(trimmed) != len(b) {
		level = slog.LevelWarn
		b = trimmed
	}

	return level, b
}
