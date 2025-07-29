package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 可替换的函数变量
var (
	MkdirAll      = os.MkdirAll
	OpenFile      = os.OpenFile
	PathSeparator = os.PathSeparator
)

var NewLogFile = func(path string) io.Writer {
	if !filePathValidator(path) {
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
	f, err := OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log file: %v, output to stdout\n", err)
		return os.Stdout
	}
	return f
}

func filePathValidator(path string) bool {
	if strings.HasSuffix(path, string(PathSeparator)) {
		return false
	}

	path = filepath.Clean(path)
	filename := filepath.Base(path)

	if strings.TrimSpace(path) == "" || strings.HasSuffix(filename, ".") {
		return false
	}

	denyChar := new(regexp.Regexp)

	if PathSeparator == '/' {
		denyChar = regexp.MustCompile(`.*[\]\[!"#$%&'()*+,\\:;<=>?@\^` + "`" + `{|}~].*`)
	} else {
		denyChar = regexp.MustCompile(`.*[\]\[!"#$%&'()*+,/;<=>?@\^` + "`" + `{|}~].*`)
	}

	return !denyChar.MatchString(path)
}
