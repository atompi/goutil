package log

import (
	"reflect"
	"testing"
)

func TestNewLoggerOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Options
		expected *Logger
	}{
		{
			name:     "No options",
			opts:     nil,
			expected: &Logger{Level: "info", Format: "console", Path: "logger", MultiFiles: false},
		},
		{
			name:     "WithLevel",
			opts:     []Options{WithLevel("debug")},
			expected: &Logger{Level: "debug", Format: "console", Path: "logger", MultiFiles: false},
		},
		{
			name:     "WithFormat",
			opts:     []Options{WithFormat("json")},
			expected: &Logger{Level: "info", Format: "json", Path: "logger", MultiFiles: false},
		},
		{
			name:     "WithPath",
			opts:     []Options{WithPath("logs/app.log")},
			expected: &Logger{Level: "info", Format: "console", Path: "logs/app.log", MultiFiles: false},
		},
		{
			name:     "WithMultiFiles",
			opts:     []Options{WithMultiFiles(true)},
			expected: &Logger{Level: "info", Format: "console", Path: "logger", MultiFiles: true},
		},
		{
			name: "Multiple options",
			opts: []Options{
				WithLevel("warn"),
				WithFormat("json"),
				WithPath("logs/app.log"),
				WithMultiFiles(true),
			},
			expected: &Logger{
				Level:      "warn",
				Format:     "json",
				Path:       "logs/app.log",
				MultiFiles: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLoggerOptions(tt.opts...)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("NewLoggerOptions() = %v, want %v", got, tt.expected)
			}
		})
	}
}
