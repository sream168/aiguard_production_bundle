package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var fileLocks sync.Map

func ResolvePath(logsDir string) string {
	return filepath.Join(logsDir, "aiguard.log")
}

func EnsureFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	return file.Close()
}

func Infof(path, format string, args ...any) {
	_ = appendLine(path, "INFO", fmt.Sprintf(format, args...))
}

func Errorf(path, format string, args ...any) {
	_ = appendLine(path, "ERROR", fmt.Sprintf(format, args...))
}

func ReadTail(path string, maxBytes int) (string, error) {
	if err := EnsureFile(path); err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if maxBytes > 0 && len(data) > maxBytes {
		data = data[len(data)-maxBytes:]
		if idx := strings.IndexByte(string(data), '\n'); idx >= 0 && idx < len(data)-1 {
			data = data[idx+1:]
		}
	}
	return string(data), nil
}

func appendLine(path, level, message string) error {
	if err := EnsureFile(path); err != nil {
		return err
	}
	mu := fileLock(path)
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf("%s [%s] %s\n", time.Now().Format(time.RFC3339), strings.ToUpper(strings.TrimSpace(level)), strings.TrimSpace(message))
	_, err = file.WriteString(line)
	return err
}

func fileLock(path string) *sync.Mutex {
	lock, _ := fileLocks.LoadOrStore(path, &sync.Mutex{})
	return lock.(*sync.Mutex)
}
