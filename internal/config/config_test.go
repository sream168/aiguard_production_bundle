package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDoesNotValidateLLMByDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("runtime:\n  concurrency: 2\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected load to succeed, got error: %v", err)
	}
	if cfg.Runtime.Concurrency != 2 {
		t.Fatalf("unexpected concurrency: %d", cfg.Runtime.Concurrency)
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validate to fail when llm config is missing")
	}
}

func TestEnvOverridesProxySettings(t *testing.T) {
	t.Setenv("AIGUARD_OPENAI_PROXY_ENABLED", "true")
	t.Setenv("AIGUARD_OPENAI_PROXY_HTTP", "http://127.0.0.1:7890")
	t.Setenv("AIGUARD_OPENAI_PROXY_HTTPS", "http://127.0.0.1:7890")
	t.Setenv("AIGUARD_OPENAI_PROXY_NO_PROXY", "127.0.0.1,localhost")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.OpenAI.Proxy.Enabled {
		t.Fatal("expected proxy to be enabled from env")
	}
	if cfg.OpenAI.Proxy.HTTP != "http://127.0.0.1:7890" {
		t.Fatalf("unexpected proxy http: %s", cfg.OpenAI.Proxy.HTTP)
	}
	if cfg.OpenAI.Proxy.HTTPS != "http://127.0.0.1:7890" {
		t.Fatalf("unexpected proxy https: %s", cfg.OpenAI.Proxy.HTTPS)
	}
	if cfg.OpenAI.Proxy.NoProxy != "127.0.0.1,localhost" {
		t.Fatalf("unexpected proxy no_proxy: %s", cfg.OpenAI.Proxy.NoProxy)
	}
}
