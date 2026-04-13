package block

import (
	"testing"

	"github.com/matryer/is"
)

// Bug 4: If a hatch block body happens to contain the literal end-marker
// string, a subsequent Inject would match the marker inside the body as the
// end of the block, corrupting the file. High severity (data corruption).
//
// The correct behavior is to refuse to inject content that contains any
// hatch marker text, so users get a clear error instead of silent damage.
func TestBug_InjectRefusesContentContainingEndMarker(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := dir + "/CLAUDE.md"
	bad := "This documents the <!-- hatch:end v1 --> marker."
	err := Inject(path, CurrentMarker, bad)
	is.True(err != nil) // must refuse, not silently corrupt
}

func TestBug_InjectRefusesContentContainingBeginMarker(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := dir + "/CLAUDE.md"
	bad := "Usage: <!-- hatch:begin v1 --> opens a block."
	err := Inject(path, CurrentMarker, bad)
	is.True(err != nil)
}
