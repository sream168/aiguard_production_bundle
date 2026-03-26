package workspace

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
)

type Layout struct {
	Root      string
	Repos     string
	Worktrees string
	Cache     string
	Reports   string
	Logs      string
}

func Prepare(root string) (Layout, error) {
	layout := Layout{
		Root:      root,
		Repos:     filepath.Join(root, "repos"),
		Worktrees: filepath.Join(root, "worktrees"),
		Cache:     filepath.Join(root, "cache"),
		Reports:   filepath.Join(root, "reports"),
		Logs:      filepath.Join(root, "logs"),
	}
	for _, dir := range []string{layout.Root, layout.Repos, layout.Worktrees, layout.Cache, layout.Reports, layout.Logs} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return layout, err
		}
	}
	return layout, nil
}

func Clear(root string) (Layout, []string, error) {
	layout := Layout{
		Root:      root,
		Repos:     filepath.Join(root, "repos"),
		Worktrees: filepath.Join(root, "worktrees"),
		Cache:     filepath.Join(root, "cache"),
		Reports:   filepath.Join(root, "reports"),
		Logs:      filepath.Join(root, "logs"),
	}
	removed := []string{layout.Repos, layout.Worktrees, layout.Cache, layout.Reports, layout.Logs}
	for _, dir := range removed {
		if err := os.RemoveAll(dir); err != nil {
			return layout, removed, err
		}
	}
	layout, err := Prepare(root)
	return layout, removed, err
}

func RepoKey(repo string) string {
	sum := sha1.Sum([]byte(repo))
	return hex.EncodeToString(sum[:])
}
