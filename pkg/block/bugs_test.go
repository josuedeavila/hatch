package block

import (
	"testing"

	"github.com/matryer/is"
)

// Bug 4: If a hatch block body contains a hatch marker on a line of its
// own, a subsequent Inject/Strip would match it as a block boundary and
// corrupt the file. High severity (data corruption).
//
// The correct behavior is to refuse such content up-front. Marker matching
// is line-aware, so bare-line markers are the ones that matter — inline
// mentions inside paragraphs or code fences are safe and must NOT be
// rejected.
func TestBug_InjectRefusesBareLineEndMarker(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := dir + "/CLAUDE.md"
	bad := "preamble\n<!-- hatch:end v1 -->\ntrailing"
	err := Inject(path, CurrentMarker, bad)
	is.True(err != nil) // must refuse, not silently corrupt
}

func TestBug_InjectRefusesBareLineBeginMarker(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := dir + "/CLAUDE.md"
	bad := "<!-- hatch:begin v1 -->\nfake block\n"
	err := Inject(path, CurrentMarker, bad)
	is.True(err != nil)
}

func TestBug_InjectAllowsInlineMarkerMention(t *testing.T) {
	// A paragraph that mentions the marker text inline (not on a line of
	// its own) is legitimate documentation and must be accepted. The
	// line-aware parser ignores it.
	is := is.New(t)
	dir := t.TempDir()
	path := dir + "/CLAUDE.md"
	doc := "This documents the <!-- hatch:end v1 --> marker format."
	err := Inject(path, CurrentMarker, doc)
	is.NoErr(err)
}
