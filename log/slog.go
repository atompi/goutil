package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
)

type handler struct {
	level       slog.Level
	writers     map[slog.Level]io.Writer
	logger      *Logger
	slogOptions *slog.HandlerOptions
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	var writer io.Writer
	if !h.logger.MultiFiles {
		writer = h.writers[convertToSlogLevel(h.logger.Level)]
	} else {
		writer = h.writers[r.Level]
	}

	switch h.logger.Format {
	case "json", "JSON":
		return slog.NewJSONHandler(writer, h.slogOptions).Handle(ctx, r)
	default:
		return slog.NewTextHandler(writer, h.slogOptions).Handle(ctx, r)
	}
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h
}

func newSlogHandler(logger *Logger, slogOptions *slog.HandlerOptions) *handler {
	h := &handler{
		level:       convertToSlogLevel(logger.Level),
		writers:     make(map[slog.Level]io.Writer),
		logger:      logger,
		slogOptions: slogOptions,
	}

	if logger.MultiFiles {
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			path := fmt.Sprintf("%s.%s.log", logger.Path, level.String())
			h.writers[level] = NewLogFile(path)
		}
	} else {
		path := fmt.Sprintf("%s.log", logger.Path)
		h.writers[convertToSlogLevel(logger.Level)] = NewLogFile(path)
	}

	return h
}

func NewSlogLogger(logger *Logger) *slog.Logger {
	slogOptions := &slog.HandlerOptions{
		AddSource:   true,
		Level:       convertToSlogLevel(logger.Level),
		ReplaceAttr: newSlogReplaceAttr(),
	}

	slogger := slog.New(newSlogHandler(logger, slogOptions))
	return slogger
}

func convertToSlogLevel(level string) slog.Level {
	switch level {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "warn", "WARN":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func newSlogReplaceAttr() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a
			}
			if source != nil {
				funcPathItems := strings.Split(source.Function, "/")
				sourceFileBase := filepath.Base(source.File)
				source.File = fmt.Sprintf(
					"%s/%s/%s",
					strings.Join(funcPathItems[0:len(funcPathItems)-1], "/"),
					strings.Split(funcPathItems[len(funcPathItems)-1], ".")[0],
					sourceFileBase,
				)
			}
		}
		return a
	}
}
