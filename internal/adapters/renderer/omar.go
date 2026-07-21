package renderer

import "embed"

//go:embed templates/omar
var omarFS embed.FS

// NewOmar returns a renderer for the omar-style Neovim plugin
// skeleton: DI-shaped M.setup(opts, deps), detector.should_load()
// gate, health module, plenary-busted test harness.
//
// Composes on top of Base — the caller runs Base first (LICENSE,
// README, .gitignore) then this for the Lua source tree. Both share
// the same emitter; ordering has no meaning because each renderer
// emits a disjoint set of paths.
func NewOmar() *EmbeddedTree {
	return &EmbeddedTree{fsys: omarFS, fsRoot: "templates/omar"}
}
