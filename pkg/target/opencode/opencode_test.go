package opencode_test

import (
	"strings"
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/hatch/pkg/target/opencode"
	"github.com/matryer/is"
)

func TestGenerate_AllFourPrimitives(t *testing.T) {
	is := is.New(t)
	s := &source.Source{Scopes: []source.Scope{{
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
	}}}
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

func TestGenerate_ScopedPrimitivesNested(t *testing.T) {
	is := is.New(t)
	s := &source.Source{Scopes: []source.Scope{
		{Path: ""},
		{
			Path: "backend",
			Rules: []source.Primitive{
				{Kind: source.KindRule, Name: "style", Body: "OC RULE"},
			},
			Skills: []source.Primitive{
				{Kind: source.KindSkill, Name: "review", Description: "d", Body: "b"},
			},
			Commands: []source.Primitive{
				{Kind: source.KindCommand, Name: "deploy", Description: "d", Body: "b"},
			},
			Agents: []source.Primitive{
				{Kind: source.KindAgent, Name: "guard", Description: "d", Body: "b"},
			},
		},
	}}
	arts, err := opencode.New().Generate(s)
	is.NoErr(err)
	byPath := map[string]target.WriteMode{}
	contentBy := map[string]string{}
	for _, a := range arts {
		byPath[a.Path] = a.Mode
		contentBy[a.Path] = a.Content
	}
	is.Equal(byPath["backend/AGENTS.md"], target.ModeBlock)
	is.True(strings.Contains(contentBy["backend/AGENTS.md"], "OC RULE"))
	is.Equal(byPath["backend/.opencode/skills/review/SKILL.md"], target.ModeFile)
	is.Equal(byPath["backend/.opencode/commands/deploy.md"], target.ModeFile)
	is.Equal(byPath["backend/.opencode/agents/guard.md"], target.ModeFile)
}
