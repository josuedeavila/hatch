package claude_test

import (
	"strings"
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/hatch/pkg/target/claude"
	"github.com/matryer/is"
)

func TestEmit_RulesAsBlockInCLAUDE(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "coding-style", Body: "Write clean code."},
		},
	}
	arts, err := claude.New().Emit(s)
	is.NoErr(err)
	var blk *target.Artifact
	for i := range arts {
		if arts[i].Path == "CLAUDE.md" {
			blk = &arts[i]
			break
		}
	}
	is.True(blk != nil)                  // CLAUDE.md artifact produced
	is.Equal(blk.Mode, target.ModeBlock) // block-injected
	is.True(strings.Contains(blk.Content, "Write clean code."))
}

func TestEmit_ApplyToRuleGetsHeading(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "go-rules", ApplyTo: "**/*.go", Body: "Go body."},
		},
	}
	arts, err := claude.New().Emit(s)
	is.NoErr(err)
	var blk *target.Artifact
	for i := range arts {
		if arts[i].Path == "CLAUDE.md" {
			blk = &arts[i]
		}
	}
	is.True(blk != nil)
	is.True(strings.Contains(blk.Content, "## Go rules (applies to `**/*.go`)"))
	is.True(strings.Contains(blk.Content, "Go body."))
}

func TestEmit_SkillBecomesFileUnderDotClaude(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Skills: []source.Primitive{{
			Kind:        source.KindSkill,
			Name:        "review-pr",
			Description: "Review a PR.",
			Body:        "do stuff\n",
		}},
	}
	arts, err := claude.New().Emit(s)
	is.NoErr(err)
	var skill *target.Artifact
	for i := range arts {
		if strings.HasSuffix(arts[i].Path, "review-pr/SKILL.md") {
			skill = &arts[i]
			break
		}
	}
	is.True(skill != nil)
	is.Equal(skill.Mode, target.ModeFile)
	is.True(strings.HasPrefix(skill.Content, "---\n"))
	is.True(strings.Contains(skill.Content, "name: review-pr"))
	is.True(strings.Contains(skill.Content, "description: Review a PR."))
	is.True(strings.Contains(skill.Content, "do stuff"))
}

func TestEmit_CommandBecomesDotClaudeCommand(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Commands: []source.Primitive{{
			Kind:        source.KindCommand,
			Name:        "commit",
			Description: "Commit with a generated message.",
			Body:        "body\n",
		}},
	}
	arts, err := claude.New().Emit(s)
	is.NoErr(err)
	found := false
	for _, a := range arts {
		if a.Path == ".claude/commands/commit.md" {
			found = true
			is.Equal(a.Mode, target.ModeFile)
		}
	}
	is.True(found)
}

func TestEmit_AgentBecomesDotClaudeAgent(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Agents: []source.Primitive{{
			Kind:        source.KindAgent,
			Name:        "security-auditor",
			Description: "Security review.",
			Body:        "body\n",
		}},
	}
	arts, err := claude.New().Emit(s)
	is.NoErr(err)
	found := false
	for _, a := range arts {
		if a.Path == ".claude/agents/security-auditor.md" {
			found = true
		}
	}
	is.True(found)
}

func TestEmit_RespectsTargetsFilter(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Skills: []source.Primitive{{
			Kind:        source.KindSkill,
			Name:        "only-opencode",
			Description: "d",
			Body:        "b",
			Targets:     []string{"opencode"},
		}},
	}
	arts, err := claude.New().Emit(s)
	is.NoErr(err)
	// The only skill opts out of claude — no artifacts expected.
	is.Equal(len(arts), 0)
}
