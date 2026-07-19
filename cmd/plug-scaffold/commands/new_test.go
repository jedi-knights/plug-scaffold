package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNew_EmitsExpectedFiles(t *testing.T) {
	dir := t.TempDir()

	if err := runNew("harpoon.nvim", "Ada Lovelace", "jedi-knights", "omar", dir); err != nil {
		t.Fatalf("runNew: %v", err)
	}

	target := filepath.Join(dir, "harpoon.nvim")
	for _, name := range []string{"LICENSE", "README.md", ".gitignore"} {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Errorf("expected %s: %v", name, err)
		}
	}
}

func TestRunNew_RefusesExistingTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "harpoon.nvim")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}

	err := runNew("harpoon.nvim", "Ada", "jedi-knights", "omar", dir)
	if err == nil {
		t.Fatal("runNew: expected error on existing target, got nil")
	}
	if !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Errorf("error should mention overwrite refusal, got: %v", err)
	}
}

func TestRunNew_RejectsInvalidStyle(t *testing.T) {
	dir := t.TempDir()
	err := runNew("harpoon", "Ada", "jedi-knights", "nightfox", dir)
	if err == nil {
		t.Fatal("runNew: expected error on invalid style")
	}
}

func TestRunNew_RejectsInvalidPluginName(t *testing.T) {
	dir := t.TempDir()
	err := runNew("Harpoon", "Ada", "jedi-knights", "omar", dir)
	if err == nil {
		t.Fatal("runNew: expected error on uppercase plugin name")
	}
}

func TestRunNew_RejectsMissingAuthorWhenGitUnset(t *testing.T) {
	// Force git to have no user.name by pointing $HOME and $XDG_CONFIG_HOME
	// at empty dirs and passing --no-system via env override. Simpler:
	// isolate git config resolution by running in an empty HOME.
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")

	err := runNew("harpoon", "", "jedi-knights", "omar", dir)
	if err == nil {
		t.Fatal("runNew: expected error when both --author and git user.name are unset")
	}
	if !strings.Contains(err.Error(), "author is required") {
		t.Errorf("error should mention author requirement, got: %v", err)
	}
}
