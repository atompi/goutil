package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHandlerEnabled tests the Enabled method of handler
func TestHandlerEnabled(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		testLevel slog.Level
		expected  bool
	}{
		{
			name:      "Debug level enabled for debug record",
			level:     "debug",
			testLevel: slog.LevelDebug,
			expected:  true,
		},
		{
			name:      "Debug level enabled for info record",
			level:     "debug",
			testLevel: slog.LevelInfo,
			expected:  true,
		},
		{
			name:      "Info level not enabled for debug record",
			level:     "info",
			testLevel: slog.LevelDebug,
			expected:  false,
		},
		{
			name:      "Warn level enabled for error record",
			level:     "warn",
			testLevel: slog.LevelError,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &Logger{
				Level: tt.level,
			}
			h := &handler{
				level:  convertToSlogLevel(tt.level),
				logger: logger,
			}

			result := h.Enabled(context.Background(), tt.testLevel)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// mockWriter is a mock implementation of io.Writer for testing
type mockWriter struct {
	buf *bytes.Buffer
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

// TestHandlerHandle tests the Handle method of handler
func TestHandlerHandle(t *testing.T) {
	ctx := context.Background()
	record := slog.Record{
		Level:   slog.LevelInfo,
		Message: "test message",
		Time:    time.Now(),
	}

	t.Run("Single file mode with text format", func(t *testing.T) {
		buf := &bytes.Buffer{}
		mockW := &mockWriter{buf: buf}

		logger := &Logger{
			Level:      "info",
			MultiFiles: false,
			Format:     "text",
		}

		h := &handler{
			level:   slog.LevelInfo,
			writers: map[slog.Level]io.Writer{slog.LevelInfo: mockW},
			logger:  logger,
		}

		err := h.Handle(ctx, record)
		if err != nil {
			t.Errorf("Handle returned error: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("Expected log output, but buffer is empty")
		}
	})

	t.Run("Multi files mode with json format", func(t *testing.T) {
		buf := &bytes.Buffer{}
		mockW := &mockWriter{buf: buf}

		logger := &Logger{
			Level:      "info",
			MultiFiles: true,
			Format:     "json",
		}

		h := &handler{
			level:   slog.LevelInfo,
			writers: map[slog.Level]io.Writer{slog.LevelInfo: mockW},
			logger:  logger,
		}

		err := h.Handle(ctx, record)
		if err != nil {
			t.Errorf("Handle returned error: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("Expected log output, but buffer is empty")
		}
	})
}

// TestHandlerWithAttrsAndGroup tests WithAttrs and WithGroup methods
func TestHandlerWithAttrsAndGroup(t *testing.T) {
	logger := &Logger{Level: "info"}
	h := &handler{logger: logger}

	resultAttrs := h.WithAttrs([]slog.Attr{})
	if resultAttrs == h {
		t.Error("WithAttrs should return a new handler")
	}

	resultGroup := h.WithGroup("test")
	if resultGroup == h {
		t.Error("WithGroup should return a new handler")
	}
}

// TestNewSlogHandler tests the newSlogHandler function
func TestNewSlogHandler(t *testing.T) {
	// Mock NewLogFile to avoid actual file creation
	originalNewLogFile := NewLogFile
	NewLogFile = func(path string) io.Writer {
		return &mockWriter{buf: &bytes.Buffer{}}
	}
	defer func() {
		NewLogFile = originalNewLogFile
	}()

	t.Run("Single file mode", func(t *testing.T) {
		logger := &Logger{
			Level:      "info",
			Path:       "/tmp/test",
			MultiFiles: false,
		}

		slogOptions := &slog.HandlerOptions{}
		h := newSlogHandler(logger, slogOptions)

		if len(h.writers) != 1 {
			t.Errorf("Expected 1 writer, got %d", len(h.writers))
		}

		expectedLevel := convertToSlogLevel(logger.Level)
		if _, exists := h.writers[expectedLevel]; !exists {
			t.Error("Expected writer for logger level")
		}
	})

	t.Run("Multi files mode", func(t *testing.T) {
		logger := &Logger{
			Level:      "info",
			Path:       "/tmp/test",
			MultiFiles: true,
		}

		slogOptions := &slog.HandlerOptions{}
		h := newSlogHandler(logger, slogOptions)

		expectedLevels := []slog.Level{
			slog.LevelDebug,
			slog.LevelInfo,
			slog.LevelWarn,
			slog.LevelError,
		}

		if len(h.writers) != len(expectedLevels) {
			t.Errorf("Expected %d writers, got %d", len(expectedLevels), len(h.writers))
		}

		for _, level := range expectedLevels {
			if _, exists := h.writers[level]; !exists {
				t.Errorf("Expected writer for level %v", level)
			}
		}
	})
}

// TestConvertToSlogLevel tests the convertToSlogLevel function
func TestConvertToSlogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"unknown", slog.LevelDebug}, // default case
		{"", slog.LevelDebug},        // default case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertToSlogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("convertToSlogLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNewSlogReplaceAttr tests the newSlogReplaceAttr function
func TestNewSlogReplaceAttr(t *testing.T) {
	replaceFunc := newSlogReplaceAttr()

	t.Run("Non-source attribute unchanged", func(t *testing.T) {
		attr := slog.String("key", "value")
		result := replaceFunc([]string{}, attr)
		if result.Key != attr.Key || result.Value.String() != attr.Value.String() {
			t.Error("Non-source attribute should remain unchanged")
		}
	})

	t.Run("Source attribute with valid source", func(t *testing.T) {
		// Get current file path for testing
		_, filename, _, _ := runtimeCaller(0)
		source := &slog.Source{
			File:     filename,
			Function: "github.com/atompi/test/test.Func",
			Line:     100,
		}
		attr := slog.Any(slog.SourceKey, source)

		result := replaceFunc([]string{}, attr)
		resultSource, ok := result.Value.Any().(*slog.Source)
		if !ok {
			t.Fatal("Result value is not *slog.Source")
		}

		// Check that the file path was converted to relative path
		expectedSource := "github.com/atompi/test/test/test.go"
		if resultSource.File != expectedSource {
			t.Errorf("Expected relative path containing %s, got %s", expectedSource, resultSource.File)
		}
	})

	t.Run("Source attribute with nil source", func(t *testing.T) {
		attr := slog.Any(slog.SourceKey, nil)
		result := replaceFunc([]string{}, attr)
		if result.Key != attr.Key {
			t.Error("Attribute key should remain unchanged")
		}
	})
}

// Helper function to mock runtime.Caller for testing
var runtimeCaller = func(_ int) (pc uintptr, file string, line int, ok bool) {
	return 0, "/tmp/src/github.com/atompi/test/test/test.go", 10, true
}

// 由于 NewSlogLogger 依赖较多外部组件，这里提供一个简化的测试思路
func TestNewSlogLogger(t *testing.T) {
	// 保存原始的 NewLogFile 函数
	originalNewLogFile := NewLogFile

	// 使用 mock 替换
	NewLogFile = MockNewLogFile
	defer func() {
		// 恢复原始函数
		NewLogFile = originalNewLogFile
	}()

	tests := []struct {
		name   string
		logger *Logger
	}{
		{
			name: "single file logger",
			logger: &Logger{
				Level:      "info",
				Format:     "json",
				Path:       "/tmp/test",
				MultiFiles: false,
			},
		},
		{
			name: "multi files logger",
			logger: &Logger{
				Level:      "debug",
				Format:     "console",
				Path:       "/tmp/test",
				MultiFiles: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewSlogLogger(tt.logger)
			assert.NotNil(t, logger)

			// 测试 logger 是否可以正常工作
			logger.Info("test message")
			logger.Debug("debug message")
			logger.Warn("warn message")
			logger.Error("error message")
		})
	}
}
