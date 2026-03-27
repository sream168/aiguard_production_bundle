package main

import (
	"path/filepath"
	"testing"
)

func TestResolveRuntimeErrorsOnMissingExplicitConfig(t *testing.T) {
	app := NewApp()
	missing := filepath.Join(t.TempDir(), "missing.yaml")

	if _, _, _, _, err := app.resolveRuntime(missing, t.TempDir()); err == nil {
		t.Fatal("expected missing explicit config path to fail")
	}
}
