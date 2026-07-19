package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// PluginSpec is the validated input to the scaffold pipeline. It carries
// exactly what the template layer needs to expand a plugin skeleton — no
// filesystem paths, no I/O, no CLI concerns.
//
// Construct via NewPluginSpec so validation and module-name derivation
// happen in one place.
type PluginSpec struct {
	// RepoName is the git repository name the user typed (e.g. "harpoon"
	// or "go-task.nvim"). Used as the emitted top-level directory name.
	RepoName string

	// ModuleName is the Lua module name that Neovim will `require`
	// (e.g. "harpoon" → require("harpoon")). Derived from RepoName by
	// stripping a trailing ".nvim" suffix when present; the ".nvim"
	// convention is a distribution suffix, not part of the module path.
	ModuleName string

	// ModuleIdent is a snake_case identifier safe to use in Vim
	// variable names (e.g. vim.g.loaded_<ident>) and Lua identifier
	// contexts. Derived from ModuleName by replacing "-" with "_".
	// Neovim's own convention: `vim.g.loaded_gitsigns`, not
	// `vim.g.loaded_git-signs`.
	ModuleIdent string

	// Author is the free-form author-attribution string. Emitted verbatim
	// into LICENSE and README.
	Author string

	// Org is the GitHub organisation or user that will own the repo.
	// Emitted into README badges and CI workflows.
	Org string

	// Style selects the template tree to render.
	Style Style
}

// repoNamePattern matches lowercase kebab-case names with an optional
// .nvim suffix — the community convention for Neovim plugin repos.
// Two-word maximum is not enforced; some legitimate plugins have longer
// names (e.g. `git-worktree.nvim`).
var repoNamePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*(?:\.nvim)?$`)

// orgNameShape matches GitHub's org/user name shape: alphanumeric with
// single hyphens allowed, no leading/trailing hyphen, 1–39 chars.
// Consecutive hyphens are rejected by a separate check because RE2 has
// no lookahead. Case is preserved (GitHub is case-insensitive on lookup
// but case-preserving in URLs).
var orgNameShape = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9-]{0,37}[A-Za-z0-9])?$`)

func isValidOrgName(s string) bool {
	return orgNameShape.MatchString(s) && !strings.Contains(s, "--")
}

// NewPluginSpec validates every field and derives ModuleName. The caller
// (usually the CLI layer) is responsible for defaulting empty inputs
// before calling — this constructor rejects them.
func NewPluginSpec(repoName, author, org string, style Style) (PluginSpec, error) {
	if !repoNamePattern.MatchString(repoName) {
		return PluginSpec{}, fmt.Errorf(
			"invalid plugin name %q: must be lowercase kebab-case with an optional .nvim suffix",
			repoName,
		)
	}
	if author == "" {
		return PluginSpec{}, fmt.Errorf("author is required")
	}
	if !isValidOrgName(org) {
		return PluginSpec{}, fmt.Errorf(
			"invalid GitHub org %q: must match GitHub's org/user rules (alphanumeric with single hyphens, 1–39 chars)",
			org,
		)
	}
	if _, err := ParseStyle(string(style)); err != nil {
		return PluginSpec{}, err
	}

	moduleName := deriveModuleName(repoName)
	return PluginSpec{
		RepoName:    repoName,
		ModuleName:  moduleName,
		ModuleIdent: deriveModuleIdent(moduleName),
		Author:      author,
		Org:         org,
		Style:       style,
	}, nil
}

// deriveModuleName strips the ".nvim" distribution suffix. The remaining
// string is the Lua module name — `lua/<module>/init.lua` gets
// `require("<module>")`. Kebab-case is preserved; Neovim's require can
// take any string.
func deriveModuleName(repoName string) string {
	return strings.TrimSuffix(repoName, ".nvim")
}

// deriveModuleIdent turns a module name into a snake_case identifier
// safe for Vim variable names. Only "-" needs translation — module
// names are already lowercase-ascii per repoNamePattern.
func deriveModuleIdent(moduleName string) string {
	return strings.ReplaceAll(moduleName, "-", "_")
}
