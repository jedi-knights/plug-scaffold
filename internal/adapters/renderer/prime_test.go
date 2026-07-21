package renderer

import (
	"strings"
	"testing"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
)

func TestPrime_Render_EmitsExpectedTree(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}

	if err := NewPrime().Render(spec, rec); err != nil {
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

func TestPrime_Render_InitUsesSingletonMetatable(t *testing.T) {
	// The distinguishing prime pattern: init.lua declares a class table
	// with __index=self, then returns a metatable-backed singleton
	// instance so callers get the same object on every require. Colon
	// syntax on setup carries `self`.
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewPrime().Render(spec, rec)

	init := string(rec.files["lua/harpoon/init.lua"])
	for _, want := range []string{
		"harpoon.__index = harpoon",
		"function harpoon:setup(opts)",
		"return setmetatable({}, harpoon)",
	} {
		if !strings.Contains(init, want) {
			t.Errorf("init.lua missing prime-style pattern %q, got:\n%s", want, init)
		}
	}
}

func TestPrime_Render_PluginRegistersNoKeymaps(t *testing.T) {
	// The prime invariant: the plugin registers no default keymaps.
	// vim.keymap.set in plugin/ would be a policy violation regardless
	// of whether plug-audit's plug-mapping rule fires (it only fires on
	// <leader> keys, but this style rejects ALL default keymaps).
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewPrime().Render(spec, rec)

	plugin := string(rec.files["plugin/harpoon.lua"])
	if strings.Contains(plugin, "vim.keymap.set") {
		t.Errorf("plugin/harpoon.lua registers a keymap; prime style forbids default bindings, got:\n%s", plugin)
	}
}

func TestPrime_Render_PluginGuardUsesModuleIdent(t *testing.T) {
	spec, err := domain.NewPluginSpec("git-worktree.nvim", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewPrime().Render(spec, rec)

	plugin := string(rec.files["plugin/git-worktree.lua"])
	if !strings.Contains(plugin, "vim.g.loaded_git_worktree") {
		t.Errorf("plugin/git-worktree.lua missing loaded_git_worktree guard, got:\n%s", plugin)
	}
	if strings.Contains(plugin, "loaded_git-worktree") {
		t.Errorf("plugin/git-worktree.lua leaked kebab into a Vim var name, got:\n%s", plugin)
	}
}

func TestPrime_Render_InitUsesModuleIdentAsClassName(t *testing.T) {
	// Kebab module names must snake_case for the class identifier —
	// Lua identifiers cannot contain hyphens. Regression guard for
	// git-worktree.nvim producing invalid `local git-worktree = {}`.
	spec, err := domain.NewPluginSpec("git-worktree.nvim", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewPrime().Render(spec, rec)

	init := string(rec.files["lua/git-worktree/init.lua"])
	if !strings.Contains(init, "local git_worktree = {}") {
		t.Errorf("init.lua class ident should be snake_case, got:\n%s", init)
	}
	if strings.Contains(init, "local git-worktree = {}") {
		t.Errorf("init.lua leaked hyphen into a Lua identifier, got:\n%s", init)
	}
}

func TestPrime_Render_HealthCallsCheckhealthStart(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewPrime().Render(spec, rec)

	health := string(rec.files["lua/harpoon/health.lua"])
	if !strings.Contains(health, `vim.health.start("harpoon")`) {
		t.Errorf("health.lua missing vim.health.start, got:\n%s", health)
	}
}

func TestPrime_Render_StaticFilesEmittedVerbatim(t *testing.T) {
	spec, err := domain.NewPluginSpec("harpoon", "Ada", "jedi-knights", domain.StylePrime)
	if err != nil {
		t.Fatalf("NewPluginSpec: %v", err)
	}
	rec := &recorder{}
	_ = NewPrime().Render(spec, rec)

	mk := string(rec.files["Makefile"])
	if !strings.Contains(mk, "PlenaryBustedDirectory") {
		t.Errorf("Makefile missing plenary target, got:\n%s", mk)
	}
	if strings.Contains(mk, "{{") {
		t.Errorf("Makefile contains unresolved template markers (should not — non-.tmpl file)")
	}
}
