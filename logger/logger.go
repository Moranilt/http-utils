package logger

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync/atomic"

	"log/slog"
)

type LoggerType string

const (
	TYPE_JSON    LoggerType = "json"
	TYPE_DEFAULT LoggerType = "default"
)

type ContextKey string

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(New(os.Stdout, TYPE_JSON))
}

func SetDefault(l Logger) {
	defaultLogger.Store(l)
}

func Default() Logger {
	return defaultLogger.Load().(Logger)
}

const (
	CtxRequestId ContextKey = "request_id"
)

const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
	LevelError  = slog.Level(4)
	LevelDebug  = slog.Level(1)
	LevelInfo   = slog.Level(0)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:  "TRACE",
	LevelNotice: "NOTICE",
	LevelFatal:  "FATAL",
	LevelError:  "ERROR",
	LevelDebug:  "DEBUG",
	LevelInfo:   "INFO",
}

type SLogger struct {
	root *slog.Logger
}

type Logger interface {
	Trace(msg string, args ...any)

	Notice(msg string, args ...any)

	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)

	Debug(msg string, args ...any)
	Debugf(format string, args ...any)

	InfoContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	Infof(format string, args ...any)

	Log(ctx context.Context, level slog.Level, msg string, args ...any)

	With(args ...any) Logger
	WithRequestInfo(r *http.Request) Logger
	WithField(key string, value any) Logger
	WithFields(fields ...any) Logger
	WithRequestId(ctx context.Context) Logger
}

func New(output io.Writer, t LoggerType) Logger {
	if output == nil {
		output = os.Stdout
	}

	var l *slog.Logger
	if t == TYPE_JSON {
		l = slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level:       LevelTrace,
			ReplaceAttr: renameLevel,
		}))
	} else {
		l = slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
			Level:       LevelTrace,
			ReplaceAttr: renameLevel,
		}))
	}

	logger := &SLogger{
		l,
	}
	return logger
}

func (s *SLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	s.root.Log(ctx, LevelInfo, msg, args...)
}

func (s *SLogger) Error(msg string, args ...any) {
	s.root.Log(context.Background(), LevelError, msg, args...)
}

func (s *SLogger) Debug(msg string, args ...any) {
	s.root.Log(context.Background(), LevelDebug, msg, args...)
}

func (s *SLogger) Trace(msg string, args ...any) {
	s.root.Log(context.Background(), LevelTrace, msg, args...)
}

func (s *SLogger) Notice(msg string, args ...any) {
	s.root.Log(context.Background(), LevelNotice, msg, args...)
}

func (s *SLogger) Fatal(msg string, args ...any) {
	s.root.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

func (s *SLogger) Fatalf(format string, args ...any) {
	s.root.Log(context.Background(), LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (s *SLogger) Errorf(format string, args ...any) {
	s.root.Log(context.Background(), LevelError, fmt.Sprintf(format, args...))
}

func (s *SLogger) Debugf(format string, args ...any) {
	s.root.Log(context.Background(), LevelDebug, fmt.Sprintf(format, args...))
}

func (s *SLogger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	s.root.Log(ctx, level, msg, args...)
}

func (s *SLogger) With(args ...any) Logger {
	return &SLogger{
		root: s.root.With(args...),
	}
}

func (s *SLogger) WithField(key string, value any) Logger {
	return &SLogger{
		root: s.root.With(key, value),
	}
}

func (s *SLogger) WithFields(fields ...any) Logger {
	return &SLogger{
		root: s.root.With(fields...),
	}
}

func (l *SLogger) Infof(format string, args ...any) {
	l.root.Log(context.Background(), LevelInfo, fmt.Sprintf(format, args...))
}

func (l *SLogger) Info(msg string, args ...any) {
	l.root.Log(context.Background(), LevelInfo, msg, args...)
}

func (l *SLogger) WithRequestInfo(r *http.Request) Logger {
	log := l.WithRequestId(r.Context())
	var clientIP string

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		clientIP = ip
	}

	return log.With(
		"path", r.URL.Path,
		"method", r.Method,
		"ip", clientIP,
	)
}
func (l *SLogger) WithRequestId(ctx context.Context) Logger {
	requestId := ctx.Value(CtxRequestId)
	if requestId != "" {
		return &SLogger{
			root: l.root.With("request_id", requestId),
		}
	}
	return l
}

func renameLevel(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelLabel, exists := LevelNames[level]
		if !exists {
			levelLabel = level.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}

	return a
}
