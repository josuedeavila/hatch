package target_test

import (
	"strings"
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/is"
)

func TestRulesBlock_ConcatenatesBodies(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "alpha", Body: "alpha body"},
			{Kind: source.KindRule, Name: "beta", Body: "beta body"},
		},
	}
	out := target.RulesBlock(s, "claude", "Claude Code")
	is.True(strings.Contains(out, "alpha body"))
	is.True(strings.Contains(out, "beta body"))
}

func TestRulesBlock_AddsHeadingOnlyForApplyTo(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "plain", Body: "plain body"},
			{Kind: source.KindRule, Name: "go-rules", ApplyTo: "**/*.go", Body: "go body"},
		},
	}
	out := target.RulesBlock(s, "claude", "Claude Code")
	is.True(!strings.Contains(out, "## Plain"))
	is.True(strings.Contains(out, "## Go rules (applies to `**/*.go`)"))
}

func TestRulesBlock_RespectsTargetsFilter(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "r1", Body: "hidden", Targets: []string{"codex"}},
			{Kind: source.KindRule, Name: "r2", Body: "visible"},
		},
	}
	out := target.RulesBlock(s, "claude", "Claude Code")
	is.True(!strings.Contains(out, "hidden"))
	is.True(strings.Contains(out, "visible"))
}

func TestRulesBlock_BodyTemplateSubstitution(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Rules: []source.Primitive{
			{Kind: source.KindRule, Name: "r", Body: "Use {{agent}} wisely."},
		},
	}
	out := target.RulesBlock(s, "claude", "Claude Code")
	is.True(strings.Contains(out, "Use Claude Code wisely."))
}

func TestSkillsBlock_RendersSkillSections(t *testing.T) {
	is := is.New(t)
	s := &source.Source{
		Skills: []source.Primitive{
			{Kind: source.KindSkill, Name: "review-pr", Description: "Review a PR.", Body: "do it"},
		},
	}
	out := target.SkillsBlock(s, "copilot", "GitHub Copilot")
	is.True(strings.Contains(out, "## Skill: review-pr"))
	is.True(strings.Contains(out, "Review a PR."))
	is.True(strings.Contains(out, "do it"))
}
