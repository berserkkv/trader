package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var (
	logger      *slog.Logger
	logLevel    string
	environment string
)

func Init(level, env string) {
	var handler slog.Handler
	logLevel = level
	environment = env

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	switch level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	if env == "local" {
		handler = NewColorTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)

	logger.Debug("Logger initialized", "level", level)
}

func Env() string {
	return environment
}

func Level() string {
	return logLevel
}

func Get() *slog.Logger {
	return logger
}

// ColorTextHandler custom handler
type ColorTextHandler struct {
	out    io.Writer
	opts   *slog.HandlerOptions
	levels map[slog.Level]string
}

func NewColorTextHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &ColorTextHandler{
		out:  w,
		opts: opts,
		levels: map[slog.Level]string{
			slog.LevelDebug: "\033[36mDEBUG\033[0m", // Cyan
			slog.LevelInfo:  "\033[32mINFO\033[0m",  // Green
			slog.LevelWarn:  "\033[33mWARN\033[0m",  // Yellow
			slog.LevelError: "\033[31mERROR\033[0m", // Red
		},
	}
}

func (h *ColorTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorTextHandler) Handle(_ context.Context, record slog.Record) error {
	levelColor := h.levels[record.Level]
	ts := time.Now().Format(time.RFC3339)

	// Colors
	messageColor := "\033[1;34m" // bright white
	keyColor := "\033[37m"       // cyan
	valueColor := "\033[35m"     // yellow
	reset := "\033[0m"

	// Build attributes string
	attrs := ""
	record.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf("%s%s=%s%v%s ", keyColor, a.Key, valueColor, a.Value, reset)
		return true
	})

	// Get the file and line number using runtime.Caller
	_, file, line, _ := runtime.Caller(3) // Skip two levels to get the caller's file/line

	// Build full log line with source info
	logLine := fmt.Sprintf("%s [%s] %s%s%s %s %s:%d\n",
		ts, levelColor,
		messageColor, record.Message, reset,
		attrs,
		file,
		line,
	)

	_, err := h.out.Write([]byte(logLine))
	return err
}

func (h *ColorTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Ignore structured fields for simplicity (can be added later if needed)
	return h
}

func (h *ColorTextHandler) WithGroup(name string) slog.Handler {
	return h
}
