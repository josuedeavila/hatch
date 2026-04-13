package opencode_test

import (
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/hatch/pkg/target/opencode"
	"github.com/matryer/is"
)

func TestGenerate_AllFourPrimitives(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "style", Body: "rule body"},
		},
		Skills: []source.Primitive{
			{Kind: source.KindSkill, Name: "review-pr", Description: "d", Body: "b"},
		},
		Commands: []source.Primitive{
			{Kind: source.KindCommand, Name: "commit", Description: "d", Body: "b"},
		},
		Agents: []source.Primitive{
			{Kind: source.KindAgent, Name: "security", Description: "d", Body: "b"},
		},
	}
	arts, err := opencode.New().Generate(s)
	is.NoErr(err)

	paths := make(map[string]target.WriteMode)
	for _, a := range arts {
		paths[a.Path] = a.Mode
	}

	is.Equal(paths["AGENTS.md"], target.ModeBlock)
	is.Equal(paths[".opencode/skills/review-pr/SKILL.md"], target.ModeFile)
	is.Equal(paths[".opencode/commands/commit.md"], target.ModeFile)
	is.Equal(paths[".opencode/agents/security.md"], target.ModeFile)
}
