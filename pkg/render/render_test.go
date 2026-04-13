package render

import (
	"testing"

	"github.com/matryer/is"
)

func TestBody(t *testing.T) {
	is := is.New(t)
	out := Body("When running {{agent}} (target={{target}}), do X.", "Claude Code", "claude")
	is.Equal(out, "When running Claude Code (target=claude), do X.")
}

func TestBody_NoSubstitutions(t *testing.T) {
	is := is.New(t)
	out := Body("plain text", "x", "y")
	is.Equal(out, "plain text")
}
