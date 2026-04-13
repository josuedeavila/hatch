package copilot_test

import (
	"strings"
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/hatch/pkg/target/copilot"
	"github.com/matryer/is"
)

func TestGenerate_UnscopedRulesBlockInCopilotInstructions(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "style", Body: "unscoped body"},
		},
	}
	arts, err := copilot.New().Generate(s)
	is.NoErr(err)
	var blk *target.Artifact
	for i := range arts {
		if arts[i].Path == ".github/copilot-instructions.md" {
			blk = &arts[i]
		}
	}
	is.True(blk != nil)
	is.Equal(blk.Mode, target.ModeBlock)
	is.True(strings.Contains(blk.Content, "unscoped body"))
}

func TestGenerate_ScopedRuleBecomesInstructionsFile(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{{
			Kind:    source.KindRule,
			Name:    "go-rules",
			ApplyTo: "**/*.go",
			Body:    "Go-only rule",
		}},
	}
	arts, err := copilot.New().Generate(s)
	is.NoErr(err)
	var inst *target.Artifact
	for i := range arts {
		if arts[i].Path == ".github/instructions/go-rules.instructions.md" {
			inst = &arts[i]
		}
	}
	is.True(inst != nil)
	is.Equal(inst.Mode, target.ModeFile)
	is.True(strings.Contains(inst.Content, "applyTo: '**/*.go'") ||
		strings.Contains(inst.Content, "applyTo: \"**/*.go\""))
	is.True(strings.Contains(inst.Content, "Go-only rule"))
}

func TestGenerate_SkillInlinedIntoCopilotInstructionsBlock(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Skills: []source.Primitive{{
			Kind:        source.KindSkill,
			Name:        "review-pr",
			Description: "review prs",
			Body:        "skill body",
		}},
	}
	arts, err := copilot.New().Generate(s)
	is.NoErr(err)
	var blk *target.Artifact
	for i := range arts {
		if arts[i].Path == ".github/copilot-instructions.md" {
			blk = &arts[i]
		}
	}
	is.True(blk != nil) // skill should be inlined as a block here
	is.True(strings.Contains(blk.Content, "review-pr"))
	is.True(strings.Contains(blk.Content, "skill body"))
}

func TestGenerate_CommandBecomesPromptFile(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Commands: []source.Primitive{{
			Kind:        source.KindCommand,
			Name:        "commit",
			Description: "commit",
			Body:        "body",
		}},
	}
	arts, err := copilot.New().Generate(s)
	is.NoErr(err)
	found := false
	for _, a := range arts {
		if a.Path == ".github/prompts/commit.prompt.md" {
			found = true
		}
	}
	is.True(found)
}

func TestGenerate_AgentBecomesAgentFile(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Agents: []source.Primitive{{
			Kind:        source.KindAgent,
			Name:        "security",
			Description: "d",
			Body:        "body",
		}},
	}
	arts, err := copilot.New().Generate(s)
	is.NoErr(err)
	found := false
	for _, a := range arts {
		if a.Path == ".github/agents/security.agent.md" {
			found = true
		}
	}
	is.True(found)
}
