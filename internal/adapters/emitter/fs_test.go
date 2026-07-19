package emitter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFS_Emit_WritesFile(t *testing.T) {
	root := t.TempDir()
	fs := New(root)

	if err := fs.Emit("README.md", []byte("hello")); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("content = %q, want %q", got, "hello")
	}
}

func TestFS_Emit_CreatesParentDirs(t *testing.T) {
	root := t.TempDir()
	fs := New(root)

	if err := fs.Emit("lua/harpoon/init.lua", []byte("return {}")); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "lua", "harpoon")); err != nil {
		t.Errorf("parent dir missing: %v", err)
	}
}

func TestFS_Emit_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	fs := New(root)

	tests := []struct {
		name string
		path string
	}{
		{"parent traversal", "../evil"},
		{"nested parent", "lua/../../evil"},
		{"absolute unix", "/etc/passwd"},
		{"empty", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := fs.Emit(tc.path, []byte("x"))
			if err == nil {
				t.Fatalf("Emit(%q): expected error, got nil", tc.path)
			}
			if !strings.Contains(err.Error(), "plug-scaffold") {
				t.Errorf("error %q should be prefixed with plug-scaffold", err)
			}
		})
	}
}

func TestFS_Emit_AllowsInternalDotDot(t *testing.T) {
	// A path like "lua/foo/../bar" resolves to "lua/bar" and stays
	// under root. The cleaner should accept it.
	root := t.TempDir()
	fs := New(root)

	if err := fs.Emit("lua/foo/../bar.lua", []byte("x")); err != nil {
		t.Fatalf("Emit: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "lua", "bar.lua")); err != nil {
		t.Errorf("expected lua/bar.lua to exist: %v", err)
	}
}
