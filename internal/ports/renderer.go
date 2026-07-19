package ports

import "github.com/jedi-knights/plug-scaffold/internal/domain"

// Renderer expands a PluginSpec into a stream of files, delegating the
// actual write to the supplied FileEmitter.
//
// Passing the emitter in (rather than returning a slice) keeps rendering
// stream-oriented: renderers can be composed by handing the same emitter
// to each one in turn, and no renderer needs to hold every file in
// memory at once. The order of Emit calls within a single Render is not
// significant to callers.
type Renderer interface {
	Render(spec domain.PluginSpec, out FileEmitter) error
}
