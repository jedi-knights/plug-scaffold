package renderer

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
	"github.com/jedi-knights/plug-scaffold/internal/ports"
)

//go:embed templates/omar
var omarFS embed.FS

// Omar renders the omar-style Neovim plugin skeleton: DI-shaped
// M.setup(opts, deps), detector.should_load() gate, health module,
// plenary-busted test harness.
//
// Composes on top of Base — the caller runs Base first (LICENSE,
// README, .gitignore) then Omar for the Lua source tree. Both share
// the same emitter; ordering has no meaning because each renderer
// emits a disjoint set of paths.
type Omar struct{}

// NewOmar returns an Omar renderer. Stateless — the same instance is
// safe to reuse across concurrent Render calls.
func NewOmar() *Omar { return &Omar{} }

// Render walks the embedded omar template tree. Paths with the
// "MODULE" placeholder segment are rewritten to the plugin's module
// name; files whose name ends in ".tmpl" are executed through
// text/template with the PluginSpec as the data root.
func (o *Omar) Render(spec domain.PluginSpec, out ports.FileEmitter) error {
	const root = "templates/omar"

	data := struct {
		Spec domain.PluginSpec
	}{spec}

	return fs.WalkDir(omarFS, root, func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(entryPath, root+"/")
		outPath := rewritePath(rel, spec)

		content, err := fs.ReadFile(omarFS, entryPath)
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

// rewritePath applies the two rewrites the template tree encodes:
//   - "MODULE" as a path segment → spec.ModuleName
//   - trailing ".tmpl" suffix → stripped
//
// The MODULE substitution is segment-scoped so a file literally
// named "MODULE_spec" doesn't have the "MODULE" prefix rewritten
// mid-segment. Using path.Base/Dir keeps the join platform-neutral.
func rewritePath(rel string, spec domain.PluginSpec) string {
	segments := strings.Split(rel, "/")
	for i, seg := range segments {
		if seg == "MODULE" {
			segments[i] = spec.ModuleName
		} else if strings.HasPrefix(seg, "MODULE") {
			// e.g. "MODULE_spec.lua.tmpl" → "harpoon_spec.lua.tmpl"
			// Preserved as a prefix rewrite so tests/MODULE_spec.lua.tmpl
			// works without needing MODULE as a separate segment.
			segments[i] = spec.ModuleName + strings.TrimPrefix(seg, "MODULE")
		}
	}
	joined := path.Join(segments...)
	return strings.TrimSuffix(joined, ".tmpl")
}

// Compile-time proof Omar satisfies the port.
var _ ports.Renderer = (*Omar)(nil)
