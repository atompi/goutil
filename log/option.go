package log

type Logger struct {
	Level      string
	Format     string
	Path       string
	MultiFiles bool
}

type Options func(*Logger)

func NewLoggerOptions(opts ...Options) *Logger {
	logger := &Logger{
		Level:      "info",
		Format:     "console",
		Path:       "logger",
		MultiFiles: false,
	}
	for _, f := range opts {
		f(logger)
	}
	return logger
}

func WithLevel(level string) Options {
	return func(logger *Logger) {
		logger.Level = level
	}
}

func WithFormat(format string) Options {
	return func(logger *Logger) {
		logger.Format = format
	}
}

func WithPath(path string) Options {
	return func(logger *Logger) {
		logger.Path = path
	}
}

func WithMultiFiles(multiFiles bool) Options {
	return func(logger *Logger) {
		logger.MultiFiles = multiFiles
	}
}
