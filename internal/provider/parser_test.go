package provider

import (
	"reflect"
	"testing"

	"aiguard/internal/config"
)

func TestParseGitLabMergeRequestWithConfiguredEndpoints(t *testing.T) {
	cfg := config.Default()
	cfg.Git.PreferredProtocol = "ssh"
	cfg.Git.GitLab.SSH.Host = "ssh.gitlab.example.com"
	cfg.Git.GitLab.SSH.Port = "2222"
	cfg.Git.GitLab.HTTPS.Host = "https.gitlab.example.com"
	cfg.Git.GitLab.HTTPS.Port = "8443"

	info := Parse("https://mr.gitlab.example.com/group/subgroup/project/-/merge_requests/123", cfg.Git)
	if info.Provider != "gitlab" {
		t.Fatalf("unexpected provider: %s", info.Provider)
	}
	if info.Path != "group/subgroup/project" {
		t.Fatalf("unexpected path: %s", info.Path)
	}
	if info.Number != "123" {
		t.Fatalf("unexpected mr number: %s", info.Number)
	}
	if info.RepoSSHURL != "ssh://git@ssh.gitlab.example.com:2222/group/subgroup/project.git" {
		t.Fatalf("unexpected ssh url: %s", info.RepoSSHURL)
	}
	if info.RepoHTTPSURL != "https://https.gitlab.example.com:8443/group/subgroup/project.git" {
		t.Fatalf("unexpected https url: %s", info.RepoHTTPSURL)
	}
	wantURLs := []string{
		"ssh://git@ssh.gitlab.example.com:2222/group/subgroup/project.git",
		"https://https.gitlab.example.com:8443/group/subgroup/project.git",
	}
	if !reflect.DeepEqual(info.RepoURLs, wantURLs) {
		t.Fatalf("unexpected repo urls: %#v", info.RepoURLs)
	}
	if info.RepoURL != wantURLs[0] {
		t.Fatalf("unexpected primary repo url: %s", info.RepoURL)
	}
}

func TestParseGitLabFallsBackToSourceHost(t *testing.T) {
	cfg := config.Default()
	info := Parse("https://gitlab.example.com/group/project/-/merge_requests/9", cfg.Git)

	if info.RepoSSHURL != "git@gitlab.example.com:group/project.git" {
		t.Fatalf("unexpected fallback ssh url: %s", info.RepoSSHURL)
	}
	if info.RepoHTTPSURL != "https://gitlab.example.com/group/project.git" {
		t.Fatalf("unexpected fallback https url: %s", info.RepoHTTPSURL)
	}
}

func TestParseGitHubPullRequest(t *testing.T) {
	cfg := config.Default()
	info := Parse("https://github.com/openai/example/pull/42", cfg.Git)

	if info.Provider != "github" {
		t.Fatalf("unexpected provider: %s", info.Provider)
	}
	if info.Owner != "openai" || info.Name != "example" {
		t.Fatalf("unexpected owner/name: %s/%s", info.Owner, info.Name)
	}
	if info.RepoURL != "git@github.com:openai/example.git" {
		t.Fatalf("unexpected repo url: %s", info.RepoURL)
	}
}
