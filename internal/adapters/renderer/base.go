// Package renderer holds Renderer implementations that emit files for a
// given PluginSpec.
package renderer

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
	"github.com/jedi-knights/plug-scaffold/internal/ports"
)

// Base renders the style-agnostic files every plug-scaffold skeleton
// ships: LICENSE, README.md, .gitignore. Style-specific renderers layer
// on top of Base by sharing the same emitter.
type Base struct {
	// Year is stamped into the LICENSE copyright line. Injected so
	// tests can freeze the year without stubbing the clock.
	Year int
}

// NewBase returns a Base renderer that stamps the current year into
// LICENSE.
func NewBase() *Base {
	return &Base{Year: time.Now().UTC().Year()}
}

// Render emits LICENSE, README.md, and .gitignore.
func (b *Base) Render(spec domain.PluginSpec, out ports.FileEmitter) error {
	files := []struct {
		path string
		tmpl string
	}{
		{"LICENSE", licenseTemplate},
		{"README.md", readmeTemplate},
		{".gitignore", gitignoreTemplate},
	}

	data := struct {
		Spec domain.PluginSpec
		Year int
	}{spec, b.Year}

	for _, f := range files {
		content, err := render(f.tmpl, data)
		if err != nil {
			return fmt.Errorf("plug-scaffold: rendering %s: %w", f.path, err)
		}
		if err := out.Emit(f.path, content); err != nil {
			return err
		}
	}
	return nil
}

func render(tmpl string, data any) ([]byte, error) {
	t, err := template.New("plug-scaffold").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Compile-time proof Base satisfies the port.
var _ ports.Renderer = (*Base)(nil)

const licenseTemplate = `MIT License

Copyright (c) {{.Year}} {{.Spec.Author}}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

const readmeTemplate = "# {{.Spec.RepoName}}\n" +
	"\n" +
	"> _Scaffolded with [`plug-scaffold`](https://github.com/jedi-knights/plug-scaffold) — style: `{{.Spec.Style}}`._\n" +
	"\n" +
	"## Install\n" +
	"\n" +
	"```lua\n" +
	"-- lazy.nvim\n" +
	"{ \"{{.Spec.Org}}/{{.Spec.RepoName}}\", opts = {} }\n" +
	"```\n" +
	"\n" +
	"## Usage\n" +
	"\n" +
	"```lua\n" +
	"require(\"{{.Spec.ModuleName}}\").setup({\n" +
	"  -- your config here\n" +
	"})\n" +
	"```\n" +
	"\n" +
	"## License\n" +
	"\n" +
	"MIT. See [LICENSE](./LICENSE).\n"

const gitignoreTemplate = `# Neovim
*.swp
*.swo

# Test/coverage output
luacov.stats.out
luacov.report.out

# Editor
.idea/
.vscode/

# macOS
.DS_Store
`
