// Package copilot generates hatch output for GitHub Copilot.
//
// Copilot reads project instructions from .github/copilot-instructions.md,
// path-scoped instructions from .github/instructions/<name>.instructions.md
// (with an `applyTo` glob), prompt files from .github/prompts/<name>.prompt.md,
// and custom agents from .github/agents/<name>.agent.md. Copilot has no
// documented skill primitive: hatch inlines skill bodies into the hatch block
// inside copilot-instructions.md so that skill content still reaches Copilot.
// See https://docs.github.com/copilot/ for the full surface.
package copilot

import (
	"path/filepath"
	"strings"

	"github.com/matryer/hatch/pkg/render"
	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
)

const (
	name        = "copilot"
	displayName = "GitHub Copilot"
)

// Target is the Copilot generator.
type Target struct{}

// New returns a Copilot target.
func New() Target { return Target{} }

func (Target) Name() string        { return name }
func (Target) DisplayName() string { return displayName }

func (t Target) Generate(s *source.Source) ([]target.Artifact, error) {
	var out []target.Artifact

	// Rules without applyTo and skills → combined block inside
	// .github/copilot-instructions.md.
	var buf strings.Builder
	rulesBody := unscopedRulesBlock(s, name, displayName)
	if rulesBody != "" {
		buf.WriteString(rulesBody)
	}
	skillsBody := target.SkillsBlock(s, name, displayName)
	if skillsBody != "" {
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString("## Skills\n\n")
		buf.WriteString(skillsBody)
	}
	if buf.Len() > 0 {
		out = append(out, target.Artifact{
			Path:    filepath.Join(".github", "copilot-instructions.md"),
			Mode:    target.ModeBlock,
			Content: buf.String(),
		})
	}

	// Rules with applyTo → .github/instructions/<name>.instructions.md
	// (native Copilot path-scoped instructions).
	for _, r := range s.Rules {
		if !r.HasTarget(name) || r.ApplyTo == "" {
			continue
		}
		content, err := renderScopedRule(r, displayName, name)
		if err != nil {
			return nil, err
		}
		out = append(out, target.Artifact{
			Path:    filepath.Join(".github", "instructions", r.Name+".instructions.md"),
			Mode:    target.ModeFile,
			Content: content,
		})
	}

	// Commands → .github/prompts/<name>.prompt.md.
	for _, c := range s.Commands {
		if !c.HasTarget(name) {
			continue
		}
		content, err := renderSlashPrimitive(c, displayName, name, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, target.Artifact{
			Path:    filepath.Join(".github", "prompts", c.Name+".prompt.md"),
			Mode:    target.ModeFile,
			Content: content,
		})
	}

	// Agents → .github/agents/<name>.agent.md.
	for _, a := range s.Agents {
		if !a.HasTarget(name) {
			continue
		}
		content, err := renderSlashPrimitive(a, displayName, name, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, target.Artifact{
			Path:    filepath.Join(".github", "agents", a.Name+".agent.md"),
			Mode:    target.ModeFile,
			Content: content,
		})
	}

	return out, nil
}

// unscopedRulesBlock returns rule bodies concatenated as markdown for rules
// that do NOT have an `applyTo` glob. Path-scoped rules are generated as
// separate `.github/instructions/` files instead of being inlined.
func unscopedRulesBlock(s *source.Source, targetName, displayName string) string {
	var buf strings.Builder
	first := true
	for _, r := range s.Rules {
		if !r.HasTarget(targetName) || r.ApplyTo != "" {
			continue
		}
		body := strings.TrimSpace(render.Body(r.Body, displayName, targetName))
		if body == "" {
			continue
		}
		if !first {
			buf.WriteString("\n\n")
		}
		first = false
		buf.WriteString(body)
	}
	return buf.String()
}

func renderScopedRule(p source.Primitive, displayName, targetName string) (string, error) {
	fields := []render.Field{
		{Key: "applyTo", Value: p.ApplyTo},
	}
	if p.Description != "" {
		fields = append(fields, render.Field{Key: "description", Value: p.Description})
	}
	if over, ok := p.Overrides[name]; ok {
		fields = render.MergeOverride(fields, over)
	}
	fm, err := render.Frontmatter(fields)
	if err != nil {
		return "", err
	}
	body := strings.TrimRight(render.Body(p.Body, displayName, targetName), "\n")
	if body == "" {
		return fm, nil
	}
	return fm + "\n" + body + "\n", nil
}

func renderSlashPrimitive(p source.Primitive, displayName, targetName string, _ any) (string, error) {
	fields := []render.Field{
		{Key: "description", Value: p.Description},
	}
	if over, ok := p.Overrides[name]; ok {
		fields = render.MergeOverride(fields, over)
	}
	fm, err := render.Frontmatter(fields)
	if err != nil {
		return "", err
	}
	body := strings.TrimRight(render.Body(p.Body, displayName, targetName), "\n")
	if body == "" {
		return fm, nil
	}
	return fm + "\n" + body + "\n", nil
}
