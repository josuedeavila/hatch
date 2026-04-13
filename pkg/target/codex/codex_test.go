package codex_test

import (
	"strings"
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/hatch/pkg/target/codex"
	"github.com/matryer/is"
)

func TestEmit_RulesAsBlockInAGENTS(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "style", Body: "rule body"},
		},
	}
	arts, err := codex.New().Emit(s)
	is.NoErr(err)
	var blk *target.Artifact
	for i := range arts {
		if arts[i].Path == "AGENTS.md" {
			blk = &arts[i]
		}
	}
	is.True(blk != nil)
	is.Equal(blk.Mode, target.ModeBlock)
	is.True(strings.Contains(blk.Content, "rule body"))
}

func TestEmit_SkillUsesAgentskillsPath(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Skills: []source.Primitive{{
			Kind: source.KindSkill, Name: "review-pr",
			Description: "d", Body: "b",
		}},
	}
	arts, err := codex.New().Emit(s)
	is.NoErr(err)
	found := false
	for _, a := range arts {
		if a.Path == ".agents/skills/review-pr/SKILL.md" {
			found = true
		}
	}
	is.True(found) // uses .agents/skills/ not .codex/skills/
}

func TestEmit_CommandsAndAgentsSkipped(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Commands: []source.Primitive{{Kind: source.KindCommand, Name: "c", Description: "d", Body: "b"}},
		Agents:   []source.Primitive{{Kind: source.KindAgent, Name: "a", Description: "d", Body: "b"}},
	}
	arts, err := codex.New().Emit(s)
	is.NoErr(err)
	for _, a := range arts {
		is.True(!strings.Contains(a.Path, "commands"))
		is.True(!strings.Contains(a.Path, "agents/"))
	}
}
