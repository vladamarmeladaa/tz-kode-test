package logger

import (
	"context"
	"fmt"
	"log/slog"
)

func NewMockLogger() *slog.Logger {
	return slog.New(NewMockHandler())
}

type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (m *MockHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (m *MockHandler) Handle(context.Context, slog.Record) error {
	return fmt.Errorf("")
}

func (m *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m
}

func (m *MockHandler) WithGroup(name string) slog.Handler {
	return m
}
