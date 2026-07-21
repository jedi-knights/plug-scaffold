package renderer

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
	"github.com/jedi-knights/plug-scaffold/internal/ports"
)

//go:embed templates/tj
var tjFS embed.FS

// TJ renders the tj-style Neovim plugin skeleton: metatable-lazy
// init.lua, hand-written vimdoc under doc/, plenary-busted tests.
//
// Composes on top of Base — the caller runs Base first (LICENSE,
// README, .gitignore) then TJ for the Lua source tree. Both share
// the same emitter; ordering has no meaning because each renderer
// emits a disjoint set of paths.
type TJ struct{}

// NewTJ returns a TJ renderer. Stateless — the same instance is safe
// to reuse across concurrent Render calls.
func NewTJ() *TJ { return &TJ{} }

// Render walks the embedded tj template tree. Path handling mirrors
// Omar.Render — "MODULE" segments rewrite to spec.ModuleName and
// ".tmpl" files pass through text/template.
func (t *TJ) Render(spec domain.PluginSpec, out ports.FileEmitter) error {
	const root = "templates/tj"

	data := struct {
		Spec domain.PluginSpec
	}{spec}

	return fs.WalkDir(tjFS, root, func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(entryPath, root+"/")
		outPath := rewritePath(rel, spec)

		content, err := fs.ReadFile(tjFS, entryPath)
		if err != nil {
			return fmt.Errorf("plug-scaffold: reading %s: %w", entryPath, err)
		}

		if strings.HasSuffix(entryPath, ".tmpl") {
			content, err = render(string(content), data)
			if err != nil {
				return fmt.Errorf("plug-scaffold: rendering %s: %w", entryPath, err)
			}
		}

		return out.Emit(outPath, content)
	})
}

// Compile-time proof TJ satisfies the port.
var _ ports.Renderer = (*TJ)(nil)
