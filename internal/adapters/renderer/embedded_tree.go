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

// EmbeddedTree renders a plugin skeleton from an embedded template
// tree rooted at fsRoot. It is the shared implementation behind every
// style renderer — NewOmar, NewTJ, and NewPrime are thin factories
// that inject the tree.
//
// Path handling:
//   - Segments equal to "MODULE" rewrite to spec.ModuleName.
//   - Segments starting with "MODULE" (e.g. "MODULE_spec.lua.tmpl")
//     rewrite as a prefix.
//   - Files ending in ".tmpl" pass through text/template with the
//     spec as the data root; other files emit verbatim.
type EmbeddedTree struct {
	fsys   embed.FS
	fsRoot string
}

// Render walks the embedded tree and emits each file through out.
func (t *EmbeddedTree) Render(spec domain.PluginSpec, out ports.FileEmitter) error {
	data := struct {
		Spec domain.PluginSpec
	}{spec}

	return fs.WalkDir(t.fsys, t.fsRoot, func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(entryPath, t.fsRoot+"/")
		outPath := rewritePath(rel, spec)

		content, err := fs.ReadFile(t.fsys, entryPath)
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

// Compile-time proof EmbeddedTree satisfies the port.
var _ ports.Renderer = (*EmbeddedTree)(nil)
