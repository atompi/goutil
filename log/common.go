package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	MkdirAll      = os.MkdirAll
	OpenFile      = os.OpenFile
	PathSeparator = os.PathSeparator
)

var (
	denyCharUnix    = regexp.MustCompile(`.*[\]\[!"#$%&'()*+,\:;<=>?@\^` + "`" + `{|}~].*`)
	denyCharWindows = regexp.MustCompile(`.*[\]\[!"#$%&'()*+,/;<=>?@\^` + "`" + `{|}~].*`)
)

type Logger struct {
	Level      string
	Format     string
	Path       string
	MultiFiles bool
}

type Config struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Path       string `yaml:"path"`
	MultiFiles bool   `yaml:"multi_files"`
}

func (c *Config) ToOptions() []Options {
	var opts []Options
	if c.Level != "" {
		opts = append(opts, WithLevel(c.Level))
	}
	if c.Format != "" {
		opts = append(opts, WithFormat(c.Format))
	}
	if c.Path != "" {
		opts = append(opts, WithPath(c.Path))
	}
	if c.MultiFiles {
		opts = append(opts, WithMultiFiles(true))
	}
	return opts
}

type Options func(*Logger)

func NewLoggerOptions(opts ...Options) *Logger {
	l := &Logger{
		Level:      "info",
		Format:     "console",
		Path:       "logger",
		MultiFiles: false,
	}
	for _, f := range opts {
		f(l)
	}
	return l
}

func WithLevel(level string) Options {
	return func(l *Logger) {
		l.Level = level
	}
}

func WithFormat(format string) Options {
	return func(l *Logger) {
		l.Format = format
	}
}

func WithPath(path string) Options {
	return func(l *Logger) {
		l.Path = path
	}
}

func WithMultiFiles(multiFiles bool) Options {
	return func(l *Logger) {
		l.MultiFiles = multiFiles
	}
}

var NewLogFile = func(path string) io.Writer {
	if !ValidPath(path) {
		fmt.Fprintf(os.Stderr, "invalid log path: %s, output to stdout\n", path)
		return os.Stdout
	}

	path = filepath.Clean(path)
	dir := filepath.Dir(path)

	err := MkdirAll(dir, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory: %v, output to stdout\n", err)
		return os.Stdout
	}

	f, err := OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log file: %v, output to stdout\n", err)
		return os.Stdout
	}

	return f
}

func ValidPath(path string) bool {
	if strings.HasSuffix(path, string(PathSeparator)) {
		return false
	}

	path = filepath.Clean(path)
	filename := filepath.Base(path)

	if strings.TrimSpace(path) == "" || strings.HasSuffix(filename, ".") {
		return false
	}

	var denyChar *regexp.Regexp
	if PathSeparator == '/' {
		denyChar = denyCharUnix
	} else {
		denyChar = denyCharWindows
	}

	return !denyChar.MatchString(path)
}
