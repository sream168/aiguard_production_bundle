package gitops

import "testing"

func TestPickTargetBranchPreference(t *testing.T) {
	branches := []string{"feature/login", "develop", "main", "master"}
	if got := pickTargetBranch(branches); got != "master" {
		t.Fatalf("expected master, got %s", got)
	}
}

func TestPickSourceBranchSkipsStableBranches(t *testing.T) {
	branches := []string{"master", "main", "develop", "feature/login", "feature/order"}
	target := pickTargetBranch(branches)
	if got := pickSourceBranch(branches, target); got != "feature/login" {
		t.Fatalf("expected feature/login, got %s", got)
	}
}

func TestNormalizeBranchName(t *testing.T) {
	cases := map[string]string{
		"origin/feature/login": "feature/login",
		"origin/HEAD":          "",
		"main":                 "main",
	}
	for raw, want := range cases {
		if got := normalizeBranchName(raw, "origin/"); got != want {
			t.Fatalf("normalizeBranchName(%q) = %q, want %q", raw, got, want)
		}
	}
}
