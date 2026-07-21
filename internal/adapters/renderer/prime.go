package renderer

import "embed"

//go:embed templates/prime
var primeFS embed.FS

// NewPrime returns a renderer for the prime-style Neovim plugin
// skeleton: singleton instance on a metatable, colon-syntax methods,
// no default keymaps.
func NewPrime() *EmbeddedTree {
	return &EmbeddedTree{fsys: primeFS, fsRoot: "templates/prime"}
}
