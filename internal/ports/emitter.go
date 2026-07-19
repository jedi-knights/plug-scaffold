// Package ports declares the interfaces the domain uses to reach the
// outside world. Adapters (under internal/adapters/) implement these;
// the domain never imports adapters directly.
package ports

// FileEmitter writes a single file relative to a target root. The path
// is expected to use forward slashes; the emitter is responsible for
// translating to the host filesystem's separator if that matters.
//
// Implementations must create parent directories as needed. They must
// refuse to write outside the target root — the CLI relies on this to
// prevent template placeholders from escaping into the user's home dir.
type FileEmitter interface {
	Emit(relPath string, content []byte) error
}
