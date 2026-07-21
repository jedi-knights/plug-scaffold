package renderer

import (
	"strings"
	"testing"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
)

func TestTJ_Render_EmitsExpectedTree(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StyleTJ)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}

	if err := NewTJ().Render(spec, rec); err != nil {
		t.Fatalf("Render: %v", err)
	}

	want := []string{
		"plugin/harpoon.lua",
		"lua/harpoon/init.lua",
		"lua/harpoon/config.lua",
		"lua/harpoon/health.lua",
		"doc/harpoon.txt",
		"tests/harpoon_spec.lua",
		"scripts/minimal_init.lua",
		"Makefile",
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

func TestTJ_Render_InitUsesMetatableLazy(t *testing.T) {
	// The distinguishing tj pattern: init.lua declares a submodule
	// allowlist and installs an __index hook that requires on demand.
	// If this regresses to eager requires at the top of the file, the
	// tj style loses its identity.
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StyleTJ)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewTJ().Render(spec, rec)

	init := string(rec.files["lua/harpoon/init.lua"])
	for _, want := range []string{
		"local submodules = {",
		"setmetatable(M, {",
		"__index = function",
	} {
		if !strings.Contains(init, want) {
			t.Errorf("init.lua missing tj-style pattern %q, got:\n%s", want, init)
		}
	}
	// Eager top-of-file requires would defeat the lazy load — regression
	// guard.
	if strings.Contains(init, `local config = require("harpoon.config")`) {
		t.Errorf("init.lua eagerly requires config; tj style must defer via __index")
	}
}

func TestTJ_Render_PluginGuardUsesModuleIdent(t *testing.T) {
	// Mirror the omar regression guard: kebab repo → snake ident in
	// Vim variable names.
	spec, err := domain.NewPluginSpec("git-worktree.nvim", "Ada", "jedi-knights", domain.StyleTJ)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewTJ().Render(spec, rec)

	plugin := string(rec.files["plugin/git-worktree.lua"])
	if !strings.Contains(plugin, "vim.g.loaded_git_worktree") {
		t.Errorf("plugin/git-worktree.lua missing loaded_git_worktree guard, got:\n%s", plugin)
	}
	if strings.Contains(plugin, "loaded_git-worktree") {
		t.Errorf("plugin/git-worktree.lua leaked kebab into a Vim var name, got:\n%s", plugin)
	}
}

func TestTJ_Render_VimdocUsesModuleName(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StyleTJ)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewTJ().Render(spec, rec)

	doc := string(rec.files["doc/harpoon.txt"])
	if !strings.Contains(doc, "*harpoon.txt*") {
		t.Errorf("doc/harpoon.txt missing help-file tag, got:\n%s", doc)
	}
	// Help tags for anchors should use the module name, not the repo name.
	if strings.Contains(doc, "harpoon.nvim-") {
		t.Errorf("vimdoc leaked .nvim suffix into help tags, got:\n%s", doc)
	}
}

func TestTJ_Render_HealthCallsCheckhealthStart(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StyleTJ)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewTJ().Render(spec, rec)

	health := string(rec.files["lua/harpoon/health.lua"])
	if !strings.Contains(health, `vim.health.start("harpoon")`) {
		t.Errorf("health.lua missing vim.health.start, got:\n%s", health)
	}
}

func TestTJ_Render_StaticFilesEmittedVerbatim(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StyleTJ)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewTJ().Render(spec, rec)

	mk := string(rec.files["Makefile"])
	if !strings.Contains(mk, "PlenaryBustedDirectory") {
		t.Errorf("Makefile missing plenary target, got:\n%s", mk)
	}
	if strings.Contains(mk, "{{") {
		t.Errorf("Makefile contains unresolved template markers (should not — non-.tmpl file)")
	}
}
