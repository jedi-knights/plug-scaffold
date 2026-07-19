// Package emitter implements the FileEmitter port against the local
// filesystem.
package emitter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedi-knights/plug-scaffold/internal/ports"
)

// FS writes files under a fixed Root directory. All Emit paths are
// resolved relative to Root and cleaned; any path that would escape Root
// (via absolute components, `..` segments, or symlink-shaped literals)
// is rejected before opening the file.
//
// Parent directories are created on demand with 0o755. Files are written
// with 0o644 — templates never emit executable content.
type FS struct {
	// Root is the target directory for all emissions. Callers must
	// create Root (or accept that it will be created on first Emit).
	Root string
}

// New returns an FS emitter rooted at root. The root is not created
// until the first Emit call.
func New(root string) *FS {
	return &FS{Root: root}
}

// Emit writes content to relPath under the emitter's root. See the FS
// type comment for path-safety rules.
func (f *FS) Emit(relPath string, content []byte) error {
	full, err := f.resolveSafe(relPath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(full)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("plug-scaffold: creating %s: %w", dir, err)
	}

	if err := os.WriteFile(full, content, 0o644); err != nil {
		return fmt.Errorf("plug-scaffold: writing %s: %w", full, err)
	}
	return nil
}

// resolveSafe joins relPath onto Root and rejects any path that resolves
// outside Root. This is a lexical check — it does not follow symlinks.
// A symlink planted between resolution and write could still redirect
// the write, but plug-scaffold's write targets are directories the user
// just asked to create, so the attack surface is limited to a directory
// the user already trusts.
func (f *FS) resolveSafe(relPath string) (string, error) {
	if relPath == "" {
		return "", fmt.Errorf("plug-scaffold: empty emit path")
	}
	// Reject absolute or Windows-drive paths early — filepath.Join
	// on an absolute rhs still escapes Root on some platforms.
	if filepath.IsAbs(relPath) || strings.HasPrefix(relPath, "/") {
		return "", fmt.Errorf("plug-scaffold: refusing absolute emit path %q", relPath)
	}

	rootAbs, err := filepath.Abs(f.Root)
	if err != nil {
		return "", fmt.Errorf("plug-scaffold: resolving root %s: %w", f.Root, err)
	}
	// filepath.Clean collapses `..` segments before joining; the
	// resulting path must remain under rootAbs.
	full := filepath.Clean(filepath.Join(rootAbs, relPath))

	rel, err := filepath.Rel(rootAbs, full)
	if err != nil {
		return "", fmt.Errorf("plug-scaffold: resolving %q: %w", relPath, err)
	}
	if strings.HasPrefix(rel, "..") || rel == ".." {
		return "", fmt.Errorf("plug-scaffold: refusing emit path %q that escapes root", relPath)
	}
	return full, nil
}

// Compile-time proof FS satisfies the port.
var _ ports.FileEmitter = (*FS)(nil)
