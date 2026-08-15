package renderer

import (
	"strings"
	"testing"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
)

func TestOmar_Render_EmitsExpectedTree(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}

	if err := NewOmar().Render(spec, rec); err != nil {
		t.Fatalf("Render: %v", err)
	}

	want := []string{
		"plugin/harpoon.lua",
		"lua/harpoon/init.lua",
		"lua/harpoon/config.lua",
		"lua/harpoon/detector.lua",
		"lua/harpoon/health.lua",
		"doc/harpoon.txt",
		"tests/harpoon_spec.lua",
		"scripts/minimal_init.lua",
		"Makefile",
		"neospec.toml",
		".github/workflows/ci.yml",
	}
	for _, p := range want {
		if _, ok := rec.files[p]; !ok {
			t.Errorf("missing emitted file: %s (got: %v)", p, keys(rec.files))
		}
	}
	if len(rec.files) != len(want) {
		t.Errorf("emitted %d files, want %d (got: %v)", len(rec.files), len(want), keys(rec.files))
	}
}

func TestOmar_Render_PluginGuardUsesModuleIdent(t *testing.T) {
	// git-worktree → module "git-worktree" → ident "git_worktree".
	// The vim.g.loaded_* var must use the snake ident, not the kebab
	// module name (Vim variable names can't contain hyphens).
	spec, err := domain.NewPluginSpec("git-worktree.nvim", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewOmar().Render(spec, rec)

	plugin := string(rec.files["plugin/git-worktree.lua"])
	if !strings.Contains(plugin, "vim.g.loaded_git_worktree") {
		t.Errorf("plugin/git-worktree.lua missing loaded_git_worktree guard, got:\n%s", plugin)
	}
	if strings.Contains(plugin, "loaded_git-worktree") {
		t.Errorf("plugin/git-worktree.lua leaked kebab into a Vim var name, got:\n%s", plugin)
	}
}

func TestOmar_Render_InitRequiresModuleSubpaths(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewOmar().Render(spec, rec)

	init := string(rec.files["lua/harpoon/init.lua"])
	for _, want := range []string{
		`require("harpoon.config")`,
		`require("harpoon.detector")`,
	} {
		if !strings.Contains(init, want) {
			t.Errorf("init.lua missing %q, got:\n%s", want, init)
		}
	}
}

func TestOmar_Render_HealthCallsCheckhealthStart(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewOmar().Render(spec, rec)

	health := string(rec.files["lua/harpoon/health.lua"])
	if !strings.Contains(health, `vim.health.start("harpoon")`) {
		t.Errorf("health.lua missing vim.health.start, got:\n%s", health)
	}
}

func TestOmar_Render_StaticFilesEmittedVerbatim(t *testing.T) {
	// Makefile and scripts/minimal_init.lua have no .tmpl suffix — they
	// should land as-is with no template execution.
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewOmar().Render(spec, rec)

	mk := string(rec.files["Makefile"])
	if !strings.Contains(mk, "PlenaryBustedDirectory") {
		t.Errorf("Makefile missing plenary target, got:\n%s", mk)
	}
	if strings.Contains(mk, "{{") {
		t.Errorf("Makefile contains unresolved template markers (should not — non-.tmpl file)")
	}
}

func TestOmar_Render_VimdocTagsUseModuleName(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewOmar().Render(spec, rec)

	doc := string(rec.files["doc/harpoon.txt"])
	for _, want := range []string{
		"*harpoon.txt*",
		"*harpoon-contents*",
		"vim:tw=78:ts=8:noet:ft=help:norl:",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("doc/harpoon.txt missing %q, got:\n%s", want, doc)
		}
	}
}

func TestOmar_Render_CIWorkflowPreservesGitHubExprSyntax(t *testing.T) {
	// .github/workflows/ci.yml is deliberately non-.tmpl so that GitHub
	// Actions ${{ github.ref }} expressions pass through verbatim rather
	// than being interpreted as Go template actions. Renaming the file
	// to .yml.tmpl would break the workflow at CI time — freeze the
	// verbatim-emit behaviour here.
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewOmar().Render(spec, rec)

	ci := string(rec.files[".github/workflows/ci.yml"])
	if !strings.Contains(ci, "${{ github.ref }}") {
		t.Errorf("ci.yml lost the GitHub Actions expression; may have been template-rendered:\n%s", ci)
	}
	if !strings.Contains(ci, "jedi-knights/neospec@v0") {
		t.Errorf("ci.yml missing neospec action pin, got:\n%s", ci)
	}
}

func keys(m map[string][]byte) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
