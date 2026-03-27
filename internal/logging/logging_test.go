package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadTailReturnsSuffix(t *testing.T) {
	path := filepath.Join(t.TempDir(), "aiguard.log")
	content := "line1\nline2\nline3\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	got, err := ReadTail(path, len("line3\n"))
	if err != nil {
		t.Fatalf("read tail: %v", err)
	}
	if got != "line3\n" {
		t.Fatalf("unexpected tail: %q", got)
	}
}

func TestAppendLineRotatesOversizedLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aiguard.log")

	if err := os.WriteFile(path, []byte(strings.Repeat("x", 2*1024*1024)), 0o644); err != nil {
		t.Fatalf("write oversized log: %v", err)
	}

	Infof(path, "new line")

	backupPath := path + ".1"
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("expected rotated backup: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read current log: %v", err)
	}
	if !strings.Contains(string(data), "new line") {
		t.Fatalf("expected current log to contain new line, got %q", string(data))
	}
}
