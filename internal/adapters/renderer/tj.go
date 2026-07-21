package renderer

import "embed"

//go:embed templates/tj
var tjFS embed.FS

// NewTJ returns a renderer for the tj-style Neovim plugin skeleton:
// metatable-lazy init.lua, hand-written vimdoc under doc/,
// plenary-busted tests.
func NewTJ() *EmbeddedTree {
	return &EmbeddedTree{fsys: tjFS, fsRoot: "templates/tj"}
}
