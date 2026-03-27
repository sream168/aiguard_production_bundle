package workspace

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
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
	resolvedRoot, err := normalizeRoot(root)
	if err != nil {
		return Layout{}, err
	}
	layout := layoutForRoot(resolvedRoot)
	for _, dir := range []string{layout.Root, layout.Repos, layout.Worktrees, layout.Cache, layout.Reports, layout.Logs} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return layout, err
		}
	}
	return layout, nil
}

func Clear(root string) (Layout, []string, error) {
	resolvedRoot, err := normalizeRoot(root)
	if err != nil {
		return Layout{}, nil, err
	}
	layout := layoutForRoot(resolvedRoot)
	if err := validateClearRoot(layout.Root); err != nil {
		return layout, nil, err
	}
	removed := []string{layout.Repos, layout.Worktrees, layout.Cache, layout.Reports, layout.Logs}
	for _, dir := range removed {
		if !isWithinRoot(layout.Root, dir) {
			return layout, removed, errors.New("工作区路径非法，拒绝清理")
		}
		if err := os.RemoveAll(dir); err != nil {
			return layout, removed, err
		}
	}
	layout, err = Prepare(layout.Root)
	return layout, removed, err
}

func RepoKey(repo string) string {
	sum := sha1.Sum([]byte(repo))
	return hex.EncodeToString(sum[:])
}

func layoutForRoot(root string) Layout {
	return Layout{
		Root:      root,
		Repos:     filepath.Join(root, "repos"),
		Worktrees: filepath.Join(root, "worktrees"),
		Cache:     filepath.Join(root, "cache"),
		Reports:   filepath.Join(root, "reports"),
		Logs:      filepath.Join(root, "logs"),
	}
}

func normalizeRoot(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return "", errors.New("工作区路径为空")
	}
	return filepath.Abs(filepath.Clean(root))
}

func validateClearRoot(root string) error {
	if isFilesystemRoot(root) {
		return errors.New("拒绝清理文件系统根目录作为工作区")
	}
	if home, err := os.UserHomeDir(); err == nil && samePath(root, home) {
		return errors.New("拒绝清理用户主目录作为工作区")
	}
	return nil
}

func isFilesystemRoot(path string) bool {
	return filepath.Dir(path) == path
}

func samePath(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	cleanA := filepath.Clean(a)
	cleanB := filepath.Clean(b)
	return strings.EqualFold(cleanA, cleanB)
}

func isWithinRoot(root, candidate string) bool {
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}
