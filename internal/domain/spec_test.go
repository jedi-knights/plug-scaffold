package domain

import (
	"strings"
	"testing"
)

func TestNewPluginSpec_Valid(t *testing.T) {
	tests := []struct {
		name           string
		repoName       string
		wantModuleName string
	}{
		{"bare name", "harpoon", "harpoon"},
		{"nvim suffix stripped", "harpoon.nvim", "harpoon"},
		{"kebab preserved", "git-worktree.nvim", "git-worktree"},
		{"digits allowed", "go-task.nvim", "go-task"},
		{"no suffix", "telescope", "telescope"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := NewPluginSpec(tc.repoName, "Ada Lovelace", "jedi-knights", StyleOmar)
			if err != nil {
				t.Fatalf("NewPluginSpec: unexpected error %v", err)
			}
			if spec.ModuleName != tc.wantModuleName {
				t.Errorf("ModuleName = %q, want %q", spec.ModuleName, tc.wantModuleName)
			}
			if spec.RepoName != tc.repoName {
				t.Errorf("RepoName = %q, want %q (RepoName must be preserved verbatim)", spec.RepoName, tc.repoName)
			}
		})
	}
}

func TestNewPluginSpec_RejectsInvalidRepoName(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
	}{
		{"empty", ""},
		{"leading hyphen", "-harpoon"},
		{"trailing hyphen", "harpoon-"},
		{"uppercase", "Harpoon"},
		{"underscore", "har_poon"},
		{"double hyphen", "har--poon"},
		{"leading digit", "1harpoon"},
		{"space", "har poon"},
		{"only suffix", ".nvim"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPluginSpec(tc.repoName, "Ada", "org", StyleOmar)
			if err == nil {
				t.Fatalf("NewPluginSpec(%q): expected error, got nil", tc.repoName)
			}
			if !strings.Contains(err.Error(), "plugin name") {
				t.Errorf("error %q should mention 'plugin name'", err)
			}
		})
	}
}

func TestNewPluginSpec_RejectsEmptyAuthor(t *testing.T) {
	_, err := NewPluginSpec("harpoon", "", "org", StyleOmar)
	if err == nil || !strings.Contains(err.Error(), "author") {
		t.Fatalf("expected author error, got %v", err)
	}
}

func TestNewPluginSpec_RejectsInvalidOrg(t *testing.T) {
	tests := []string{"", "-jedi", "jedi-", "jedi--knights", strings.Repeat("a", 40)}
	for _, org := range tests {
		t.Run(org, func(t *testing.T) {
			_, err := NewPluginSpec("harpoon", "Ada", org, StyleOmar)
			if err == nil {
				t.Fatalf("NewPluginSpec(org=%q): expected error, got nil", org)
			}
		})
	}
}

func TestNewPluginSpec_RejectsUnknownStyle(t *testing.T) {
	_, err := NewPluginSpec("harpoon", "Ada", "org", Style("nightfox"))
	if err == nil {
		t.Fatal("NewPluginSpec: expected error for unknown style, got nil")
	}
}
