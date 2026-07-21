package renderer

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
	"github.com/jedi-knights/plug-scaffold/internal/ports"
)

//go:embed templates/prime
var primeFS embed.FS

// Prime renders the prime-style Neovim plugin skeleton: singleton
// instance on a metatable, colon-syntax methods, no default keymaps.
//
// Composes on top of Base — the caller runs Base first (LICENSE,
// README, .gitignore) then Prime for the Lua source tree. Both share
// the same emitter; ordering has no meaning because each renderer
// emits a disjoint set of paths.
type Prime struct{}

// NewPrime returns a Prime renderer. Stateless — the same instance is
// safe to reuse across concurrent Render calls.
func NewPrime() *Prime { return &Prime{} }

// Render walks the embedded prime template tree. Path handling mirrors
// Omar.Render — "MODULE" segments rewrite to spec.ModuleName and
// ".tmpl" files pass through text/template.
func (p *Prime) Render(spec domain.PluginSpec, out ports.FileEmitter) error {
	const root = "templates/prime"

	data := struct {
		Spec domain.PluginSpec
	}{spec}

	return fs.WalkDir(primeFS, root, func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(entryPath, root+"/")
		outPath := rewritePath(rel, spec)

		content, err := fs.ReadFile(primeFS, entryPath)
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

// Compile-time proof Prime satisfies the port.
var _ ports.Renderer = (*Prime)(nil)
