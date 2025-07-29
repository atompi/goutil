package log

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

// MockWriter 是 io.Writer 的 mock 实现
type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// MockNewLogFile 是 NewLogFile 的 mock 实现
func MockNewLogFile(path string) io.Writer {
	mockWriter := new(MockWriter)
	mockWriter.On("Write", mock.Anything).Return(len("test"), nil)
	return mockWriter
}

func TestNewZapEncoder(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   string // encoder type
	}{
		{
			name:   "json format",
			format: "json",
			want:   "json",
		},
		{
			name:   "JSON format",
			format: "JSON",
			want:   "json",
		},
		{
			name:   "console format",
			format: "console",
			want:   "console",
		},
		{
			name:   "empty format",
			format: "",
			want:   "console",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := newZapEncoder(tt.format)
			assert.NotNil(t, encoder)

			// 通过编码器是否能正常编码判断其可用性
			assert.NotPanics(t, func() {
				encoder.EncodeEntry(zapcore.Entry{}, nil)
			})
		})
	}
}

func TestConvertToZapLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zapcore.Level
	}{
		{
			name:  "debug level",
			level: "debug",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "DEBUG level",
			level: "DEBUG",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "info level",
			level: "info",
			want:  zapcore.InfoLevel,
		},
		{
			name:  "INFO level",
			level: "INFO",
			want:  zapcore.InfoLevel,
		},
		{
			name:  "warn level",
			level: "warn",
			want:  zapcore.WarnLevel,
		},
		{
			name:  "WARN level",
			level: "WARN",
			want:  zapcore.WarnLevel,
		},
		{
			name:  "error level",
			level: "error",
			want:  zapcore.ErrorLevel,
		},
		{
			name:  "ERROR level",
			level: "ERROR",
			want:  zapcore.ErrorLevel,
		},
		{
			name:  "unknown level",
			level: "unknown",
			want:  zapcore.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToZapLevel(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewZapLevels(t *testing.T) {
	allLevels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}

	tests := []struct {
		name  string
		level string
		want  []zapcore.Level
	}{
		{
			name:  "debug level",
			level: "debug",
			want:  allLevels,
		},
		{
			name:  "DEBUG level",
			level: "DEBUG",
			want:  allLevels,
		},
		{
			name:  "info level",
			level: "info",
			want:  allLevels[1:],
		},
		{
			name:  "INFO level",
			level: "INFO",
			want:  allLevels[1:],
		},
		{
			name:  "warn level",
			level: "warn",
			want:  allLevels[2:],
		},
		{
			name:  "WARN level",
			level: "WARN",
			want:  allLevels[2:],
		},
		{
			name:  "error level",
			level: "error",
			want:  allLevels[3:],
		},
		{
			name:  "ERROR level",
			level: "ERROR",
			want:  allLevels[3:],
		},
		{
			name:  "unknown level",
			level: "unknown",
			want:  allLevels,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newZapLevels(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewZapLevelEnablerFunc(t *testing.T) {
	tests := []struct {
		name        string
		level       zapcore.Level
		multiFiles  bool
		testLevels  []zapcore.Level
		wantEnabled []bool
	}{
		{
			name:        "DebugLevel single file",
			level:       zapcore.DebugLevel,
			multiFiles:  false,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{true, true, true, true},
		},
		{
			name:        "DebugLevel multi files",
			level:       zapcore.DebugLevel,
			multiFiles:  true,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{true, false, false, false},
		},
		{
			name:        "InfoLevel single file",
			level:       zapcore.InfoLevel,
			multiFiles:  false,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{false, true, true, true},
		},
		{
			name:        "InfoLevel multi files",
			level:       zapcore.InfoLevel,
			multiFiles:  true,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{false, true, false, false},
		},
		{
			name:        "WarnLevel single file",
			level:       zapcore.WarnLevel,
			multiFiles:  false,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{false, false, true, true},
		},
		{
			name:        "WarnLevel multi files",
			level:       zapcore.WarnLevel,
			multiFiles:  true,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{false, false, true, false},
		},
		{
			name:        "ErrorLevel single file",
			level:       zapcore.ErrorLevel,
			multiFiles:  false,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{false, false, false, true},
		},
		{
			name:        "ErrorLevel multi files",
			level:       zapcore.ErrorLevel,
			multiFiles:  true,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{false, false, false, true},
		},
		{
			name:        "UnknownLevel single file",
			level:       zapcore.InvalidLevel,
			multiFiles:  false,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{true, true, true, true},
		},
		{
			name:        "UnknownLevel multi files",
			level:       zapcore.InvalidLevel,
			multiFiles:  true,
			testLevels:  []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel},
			wantEnabled: []bool{true, false, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enablerFunc := newZapLevelEnablerFunc(tt.level, tt.multiFiles)
			assert.NotNil(t, enablerFunc)

			for i, testLevel := range tt.testLevels {
				got := enablerFunc.Enabled(testLevel)
				assert.Equal(t, tt.wantEnabled[i], got, "level %v", testLevel)
			}
		})
	}
}

// 由于 NewZapLogger 依赖较多外部组件，这里提供一个简化的测试思路
func TestNewZapLogger(t *testing.T) {
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
			logger := NewZapLogger(tt.logger)
			assert.NotNil(t, logger)

			// 测试 logger 是否可以正常工作
			logger.Info("test message")
			logger.Debug("debug message")
			logger.Warn("warn message")
			logger.Error("error message")
		})
	}
}
