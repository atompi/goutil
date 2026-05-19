package log

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

// TestFilePathValidator 测试 filePathValidator 函数
func TestFilePathValidator(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		pathSeparator rune
		expected      bool
	}{
		{
			name:          "empty string",
			path:          "",
			pathSeparator: '/',
			expected:      false,
		},
		{
			name:          "whitespace only",
			path:          "   ",
			pathSeparator: '/',
			expected:      false,
		},
		{
			name:          "valid linux path",
			path:          "/tmp/logs/app.log",
			pathSeparator: '/',
			expected:      true,
		},
		{
			name:          "invalid characters linux",
			path:          "/tmp/log*file.log",
			pathSeparator: '/',
			expected:      false,
		},
		{
			name:          "ends with separator linux",
			path:          "/tmp/logs/",
			pathSeparator: '/',
			expected:      false,
		},
		{
			name:          "valid linux relative path",
			path:          "logs/app.log",
			pathSeparator: '/',
			expected:      true,
		},
		{
			name:          "valid windows path",
			path:          "C:\\logs\\app.log",
			pathSeparator: '\\',
			expected:      true,
		},
		{
			name:          "invalid characters windows",
			path:          "C:\\log\\file<>.log",
			pathSeparator: '\\',
			expected:      false,
		},
		{
			name:          "ends with separator windows",
			path:          "C:\\logs\\",
			pathSeparator: '\\',
			expected:      false,
		},
		{
			name:          "valid windows relative path",
			path:          "logs\\app.log",
			pathSeparator: '\\',
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PathSeparator = tt.pathSeparator
			result := ValidPath(tt.path)
			if result != tt.expected {
				t.Errorf("filePathValidator(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}

	// pop original function
	defer func() {
		PathSeparator = os.PathSeparator
	}()
}

// TestNewLogFileInvalidPath 测试 NewLogFile 函数的无效路径情况
func TestNewLogFileInvalidPath(t *testing.T) {
	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Test with invalid path
	writer := NewLogFile("")

	w.Close()
	os.Stderr = oldStderr
	<-done

	// Check that output was written to stderr
	if !strings.Contains(buf.String(), "invalid log path") {
		t.Errorf("Expected error message in stderr, got: %s", buf.String())
	}

	// Check that writer is stdout
	if writer != os.Stdout {
		t.Error("Expected writer to be os.Stdout for invalid path")
	}
}

// TestNewLogFileValidPath 测试 NewLogFile 函数的有效路径情况
func TestNewLogFileValidPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Test with valid path
	writer := NewLogFile(logPath)

	// Check that writer is not stdout
	if writer == os.Stdout {
		t.Error("Expected writer to be a file, not os.Stdout")
	}

	// Try to write to the file
	testMessage := "test log message\n"
	_, err := writer.Write([]byte(testMessage))
	if err != nil {
		t.Errorf("Failed to write to log file: %v", err)
	}

	// Close the file
	if closer, ok := writer.(io.Closer); ok {
		closer.Close()
	}

	// Verify content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}
	if string(content) != testMessage {
		t.Errorf("Expected content %q, got %q", testMessage, string(content))
	}
}

// TestNewLogFileDirectoryCreation 测试 NewLogFile 函数的目录自动创建功能
func TestNewLogFileDirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "subdir1", "subdir2", "test.log")

	// Test that directories are created automatically
	writer := NewLogFile(logPath)

	if writer == os.Stdout {
		t.Error("Expected writer to be a file, not os.Stdout")
	}

	// Check that directories were created
	dir := filepath.Dir(logPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("Expected directories to be created automatically")
	}

	// Clean up
	if closer, ok := writer.(io.Closer); ok {
		closer.Close()
	}
}

// TestNewLogFileWithMockedMkdirAll 测试 NewLogFile 函数，模拟 MkdirAll 失败
func TestNewLogFileWithMockedMkdirAll(t *testing.T) {
	// stash original function
	originalMkdirAll := MkdirAll

	// Mock os.MkdirAll to return an error
	MkdirAll = func(string, os.FileMode) error {
		return os.ErrPermission
	}

	// pop original function
	defer func() {
		MkdirAll = originalMkdirAll
	}()

	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Test with valid path that requires directory creation
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "subdir", "test.log")
	writer := NewLogFile(logPath)

	w.Close()
	os.Stderr = oldStderr
	<-done

	// Check that error message was written to stderr
	if !strings.Contains(buf.String(), "failed to create log directory") {
		t.Errorf("Expected directory creation error message, got: %s", buf.String())
	}

	// Check that writer is stdout when directory creation fails
	if writer != os.Stdout {
		t.Error("Expected writer to be os.Stdout when directory creation fails")
	}
}

// TestNewLogFileWithMockedOpenFile 测试 NewLogFile 函数，模拟 OpenFile 失败
func TestNewLogFileWithMockedOpenFile(t *testing.T) {
	// stash original function
	originalOpenFile := OpenFile

	// Mock OpenFile to return an error
	OpenFile = func(string, int, os.FileMode) (*os.File, error) {
		return nil, os.ErrPermission
	}

	// pop original function
	defer func() {
		OpenFile = originalOpenFile
	}()

	// Capture stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Test with valid existing directory path
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	writer := NewLogFile(logPath)

	w.Close()
	os.Stderr = oldStderr
	<-done

	// Check that error message was written to stderr
	if !strings.Contains(buf.String(), "failed to create log file") {
		t.Errorf("Expected file creation error message, got: %s", buf.String())
	}

	// Check that writer is stdout when file creation fails
	if writer != os.Stdout {
		t.Error("Expected writer to be os.Stdout when file creation fails")
	}
}

// TestFilePathValidatorRegex tests different OS path validation
func TestFilePathValidatorRegex(t *testing.T) {
	if os.PathSeparator == '/' {
		// Test Linux path separator
		// On Unix, backslash is technically allowed in filenames
		result := ValidPath("/tmp/log\\file.log")
		if result {
			t.Log("Note: backslash is allowed in Unix filenames (though bad practice)")
		}

		// Should reject asterisk
		result = ValidPath("/tmp/log*file.log")
		if result {
			t.Error("Linux path should not allow asterisk")
		}
	} else {
		// Test Windows path separator
		// Should reject forward slash
		result := ValidPath("C:\\log/file.log")
		if result {
			t.Error("Windows path should not allow forward slash")
		}

		// Should reject less than sign
		result = ValidPath("C:\\log<file.log")
		if result {
			t.Error("Windows path should not allow less than sign")
		}
	}
}
