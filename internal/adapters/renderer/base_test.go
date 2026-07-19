package renderer

import (
	"strings"
	"testing"

	"github.com/jedi-knights/plug-scaffold/internal/domain"
)

type recorder struct {
	files map[string][]byte
}

func (r *recorder) Emit(path string, content []byte) error {
	if r.files == nil {
		r.files = map[string][]byte{}
	}
	r.files[path] = content
	return nil
}

func fixtureSpec(t *testing.T) domain.PluginSpec {
	t.Helper()
	spec, err := domain.NewPluginSpec("harpoon.nvim", "Ada Lovelace", "jedi-knights", domain.StyleOmar)
	if err != nil {
		t.Fatalf("fixtureSpec: %v", err)
	}
	return spec
}

func TestBase_Render_EmitsExpectedFiles(t *testing.T) {
	base := &Base{Year: 2026}
	rec := &recorder{}

	if err := base.Render(fixtureSpec(t), rec); err != nil {
		t.Fatalf("Render: %v", err)
	}

	want := []string{"LICENSE", "README.md", ".gitignore"}
	for _, path := range want {
		if _, ok := rec.files[path]; !ok {
			t.Errorf("missing emitted file: %s", path)
		}
	}
	if len(rec.files) != len(want) {
		t.Errorf("emitted %d files, want %d", len(rec.files), len(want))
	}
}

func TestBase_Render_LicenseStampsAuthorAndYear(t *testing.T) {
	base := &Base{Year: 2026}
	rec := &recorder{}
	_ = base.Render(fixtureSpec(t), rec)

	license := string(rec.files["LICENSE"])
	if !strings.Contains(license, "Copyright (c) 2026 Ada Lovelace") {
		t.Errorf("LICENSE missing copyright line, got:\n%s", license)
	}
}

func TestBase_Render_ReadmeUsesModuleNameNotRepoName(t *testing.T) {
	// The README's `require("...")` snippet must use the derived module
	// name ("harpoon"), not the repo name ("harpoon.nvim"). Getting
	// this wrong is the top reason a scaffolded plugin fails to load.
	base := &Base{Year: 2026}
	rec := &recorder{}
	_ = base.Render(fixtureSpec(t), rec)

	readme := string(rec.files["README.md"])
	if !strings.Contains(readme, `require("harpoon")`) {
		t.Errorf("README require snippet should use module name 'harpoon', got:\n%s", readme)
	}
	if strings.Contains(readme, `require("harpoon.nvim")`) {
		t.Errorf("README require snippet leaked the .nvim suffix, got:\n%s", readme)
	}
}

func TestBase_Render_ReadmeIncludesOrgInInstall(t *testing.T) {
	base := &Base{Year: 2026}
	rec := &recorder{}
	_ = base.Render(fixtureSpec(t), rec)

	readme := string(rec.files["README.md"])
	if !strings.Contains(readme, `"jedi-knights/harpoon.nvim"`) {
		t.Errorf("README install snippet should use org/repo, got:\n%s", readme)
	}
}

func TestNewBase_YearIsCurrent(t *testing.T) {
	// NewBase should stamp *some* recent year — we don't pin an exact
	// value, just guard against a zero-year regression.
	b := NewBase()
	if b.Year < 2025 {
		t.Errorf("NewBase().Year = %d, want >= 2025", b.Year)
	}
}
