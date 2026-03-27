package review

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"aiguard/internal/model"
	"aiguard/internal/uiapi"
)

func TestNewOrchestrator(t *testing.T) {
	orch := NewOrchestrator()
	if orch == nil {
		t.Error("expected non-nil orchestrator")
	}
	if orch.git == nil {
		t.Error("expected non-nil git manager")
	}
	if orch.locker == nil {
		t.Error("expected non-nil locker")
	}
}

func TestEmitProgress(t *testing.T) {
	called := false
	emit := func(name string, payload any) {
		called = true
	}
	emitProgress(emit, "task1", "测试", 50, "测试消息", model.Summary{})
	if !called {
		t.Error("emit function not called")
	}
}

func TestRunHonorsReviewConfigFlags(t *testing.T) {
	repoDir := initReviewRepo(t)
	workspaceDir := t.TempDir()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	configContent := "review:\n" +
		"  workspace_dir: " + quoteYAML(workspaceDir) + "\n" +
		"  export_formats: [\"json\"]\n" +
		"  enable_project_brief: false\n" +
		"  enable_prescan: false\n" +
		"  redact_secrets_before_llm: true\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	done, err := NewOrchestrator().Run(context.Background(), "task-flags", uiapi.StartReviewRequest{
		LocalRepoPath: repoDir,
		SourceBranch:  "feature/redaction",
		TargetBranch:  "main",
		ConfigPath:    configPath,
		WorkspaceDir:  workspaceDir,
	}, func(string, any) {})
	if err != nil {
		t.Fatalf("run review: %v", err)
	}

	if !reflect.DeepEqual(done.Report.ProjectBrief, model.ProjectBrief{}) {
		t.Fatal("expected project brief to be skipped")
	}
	if len(done.Report.Findings) != 0 {
		t.Fatalf("expected prescan findings to be skipped, got %d", len(done.Report.Findings))
	}
	if _, err := os.Stat(done.HTMLPath); err == nil {
		t.Fatal("expected html export to be skipped")
	}
	if _, err := os.Stat(done.MarkdownPath); err == nil {
		t.Fatal("expected markdown export to be skipped")
	}

	packsPath := filepath.Join(done.ReportDir, "artifacts", "review_packs.json")
	data, err := os.ReadFile(packsPath)
	if err != nil {
		t.Fatalf("read review packs: %v", err)
	}
	var packs []model.ReviewPack
	if err := json.Unmarshal(data, &packs); err != nil {
		t.Fatalf("unmarshal review packs: %v", err)
	}
	if len(packs) == 0 {
		t.Fatal("expected review packs to be generated")
	}
	if packs[0].DiffText == "" {
		t.Fatal("expected non-empty diff text")
	}
	if containsToken(packs[0].DiffText) {
		t.Fatal("expected review packs to redact bearer token")
	}
}

func initReviewRepo(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()
	runGit(t, repoDir, "init", "-b", "main")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "tester")

	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("# demo\n"), 0o644); err != nil {
		t.Fatalf("write readme: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "base")

	runGit(t, repoDir, "checkout", "-b", "feature/redaction")
	content := "package main\n\nconst api_key = \"1234567890\"\nconst header = \"Authorization: Bearer super-secret-token\"\n"
	if err := os.WriteFile(filepath.Join(repoDir, "main.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("write feature main.go: %v", err)
	}
	runGit(t, repoDir, "add", "main.go")
	runGit(t, repoDir, "commit", "-m", "feature")

	return repoDir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}

func quoteYAML(value string) string {
	return "\"" + filepath.ToSlash(value) + "\""
}

func containsToken(text string) bool {
	return strings.Contains(text, "super-secret-token")
}
