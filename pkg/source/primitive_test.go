package source

import (
	"testing"

	"github.com/matryer/is"
)

func TestKindConstants(t *testing.T) {
	is := is.New(t)
	is.Equal(string(KindRule), "rule")
	is.Equal(string(KindSkill), "skill")
	is.Equal(string(KindCommand), "command")
	is.Equal(string(KindAgent), "agent")
}

func TestPrimitive_HasTarget_EmptyMeansAll(t *testing.T) {
	is := is.New(t)
	p := &Primitive{}
	is.True(p.HasTarget("claude"))
	is.True(p.HasTarget("anything"))
}

func TestPrimitive_HasTarget_ExplicitList(t *testing.T) {
	is := is.New(t)
	p := &Primitive{Targets: []string{"claude", "codex"}}
	is.True(p.HasTarget("claude"))
	is.True(p.HasTarget("codex"))
	is.True(!p.HasTarget("opencode"))
	is.True(!p.HasTarget("copilot"))
}

func TestSource_ZeroValue(t *testing.T) {
	is := is.New(t)
	s := &Source{}
	is.Equal(len(s.Rules), 0)
	is.Equal(len(s.Skills), 0)
	is.Equal(len(s.Commands), 0)
	is.Equal(len(s.Agents), 0)
}
