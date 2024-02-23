package logger

import (
	"context"
	"log/slog"
	"net/http"
)

func NewMock() Logger {
	return &MockLogger{}
}

type MockLogger struct{}

func (m *MockLogger) Trace(msg string, args ...any) {}

func (m *MockLogger) Notice(msg string, args ...any) {}

func (m *MockLogger) Error(msg string, args ...any) {}

func (m *MockLogger) Fatal(msg string, args ...any) {}

func (m *MockLogger) Fatalf(format string, args ...any) {}

func (m *MockLogger) Errorf(format string, args ...any) {}

func (m *MockLogger) Debug(msg string, args ...any) {}

func (m *MockLogger) Debugf(format string, args ...any) {}

func (m *MockLogger) Info(msg string, args ...any) {}

func (m *MockLogger) Infof(format string, args ...any) {}

func (m *MockLogger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {}

func (m *MockLogger) With(args ...any) Logger {
	return m
}

func (m *MockLogger) WithRequestInfo(r *http.Request) Logger {
	return m
}

func (m *MockLogger) WithField(key string, value any) Logger {
	return m
}

func (m *MockLogger) WithFields(fields ...any) Logger {
	return m
}

func (m *MockLogger) WithRequestId(ctx context.Context) Logger {
	return m
}

func (m *MockLogger) InfoContext(ctx context.Context, msg string, args ...any) {}
