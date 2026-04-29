package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	callerIDKey  contextKey = "caller_id"
)

type customTextHandler struct {
	out   io.Writer
	level slog.Level
	mu    sync.Mutex
}

var appLogger = slog.New(newCustomTextHandler(os.Stdout, slog.LevelInfo))

// Init configures the global logger based on runtime environment.
func Init(environment string) {
	level := slog.LevelInfo
	if strings.EqualFold(environment, "development") {
		level = slog.LevelDebug
	}
	appLogger = slog.New(newCustomTextHandler(os.Stdout, level))
}

// WithRequestID stores request id in context for downstream logs.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// WithCallerID stores caller id in context for downstream logs.
func WithCallerID(ctx context.Context, callerID string) context.Context {
	return context.WithValue(ctx, callerIDKey, callerID)
}

// Info logs an info-level message with function context.
func Info(ctx context.Context, funcName, msg string, kv ...any) {
	appLogger.Log(ctx, slog.LevelInfo, msg, buildLogAttrs(ctx, funcName, kv...)...)
}

// Debug logs a debug-level message with function context.
func Debug(ctx context.Context, funcName, msg string, kv ...any) {
	appLogger.Log(ctx, slog.LevelDebug, msg, buildLogAttrs(ctx, funcName, kv...)...)
}

// Warn logs a warn-level message with function context.
func Warn(ctx context.Context, funcName, msg string, kv ...any) {
	appLogger.Log(ctx, slog.LevelWarn, msg, buildLogAttrs(ctx, funcName, kv...)...)
}

// Error logs an error-level message with function context.
func Error(ctx context.Context, funcName, msg string, kv ...any) {
	appLogger.Log(ctx, slog.LevelError, msg, buildLogAttrs(ctx, funcName, kv...)...)
}

// newCustomTextHandler creates a log handler with custom output format.
func newCustomTextHandler(out io.Writer, level slog.Level) slog.Handler {
	return &customTextHandler{
		out:   out,
		level: level,
	}
}

// Enabled reports whether this handler should emit the log level.
func (h *customTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle prints log records in [Function] - message key=value format.
func (h *customTextHandler) Handle(_ context.Context, record slog.Record) error {
	funcName := "unknown"
	attrs := make([]string, 0, record.NumAttrs()+1)

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "func" {
			if attr.Value.Kind() == slog.KindString && attr.Value.String() != "" {
				funcName = attr.Value.String()
			}
			return true
		}
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
		return true
	})

	if record.Level != slog.LevelInfo {
		attrs = append([]string{fmt.Sprintf("level=%s", strings.ToUpper(record.Level.String()))}, attrs...)
	}

	line := fmt.Sprintf("[%s] - %s", funcName, record.Message)
	if len(attrs) > 0 {
		line += " " + strings.Join(attrs, " ")
	}
	line += "\n"

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.out, line)
	return err
}

// WithAttrs returns a handler that includes persistent attributes.
func (h *customTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &withAttrsHandler{
		base:  h,
		attrs: attrs,
	}
}

// WithGroup is a no-op because groups are not used in this format.
func (h *customTextHandler) WithGroup(_ string) slog.Handler {
	return h
}

type withAttrsHandler struct {
	base  *customTextHandler
	attrs []slog.Attr
}

// Enabled reports whether this handler should emit the log level.
func (h *withAttrsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

// Handle appends persistent attributes before delegating to base handler.
func (h *withAttrsHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, attr := range h.attrs {
		record.AddAttrs(attr)
	}
	return h.base.Handle(ctx, record)
}

// WithAttrs returns a new handler with merged persistent attributes.
func (h *withAttrsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	combined := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	combined = append(combined, h.attrs...)
	combined = append(combined, attrs...)
	return &withAttrsHandler{
		base:  h.base,
		attrs: combined,
	}
}

// WithGroup is a no-op because groups are not used in this format.
func (h *withAttrsHandler) WithGroup(_ string) slog.Handler {
	return h
}

// buildLogAttrs adds function and context metadata to log attributes.
func buildLogAttrs(ctx context.Context, funcName string, kv ...any) []any {
	attrs := make([]any, 0, len(kv)+6)
	attrs = append(attrs, "func", funcName)

	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}
	if callerID, ok := ctx.Value(callerIDKey).(string); ok && callerID != "" {
		attrs = append(attrs, "caller_id", callerID)
	}

	attrs = append(attrs, kv...)
	return attrs
}
