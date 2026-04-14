// Package opencode generates hatch output for OpenCode (sst/opencode).
//
// OpenCode reads rules from AGENTS.md (or CLAUDE.md as a fallback), skills
// from .opencode/skills/<name>/SKILL.md, agents from .opencode/agents/<name>.md,
// and user commands from .opencode/commands/<name>.md. See
// https://opencode.ai/docs/ for the full surface.
package opencode

import (
	"path/filepath"
	"strings"

	"github.com/matryer/hatch/pkg/render"
	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
)

const (
	name        = "opencode"
	displayName = "OpenCode"
)

// Target is the OpenCode generator.
type Target struct{}

// New returns an OpenCode target.
func New() Target { return Target{} }

func (Target) Name() string        { return name }
func (Target) DisplayName() string { return displayName }

func (t Target) Generate(s *source.Source) ([]target.Artifact, error) {
	var out []target.Artifact

	// Rules → block inside AGENTS.md (shared with Codex; identical content).
	if body := target.RulesBlock(s, name, displayName); body != "" {
		out = append(out, target.Artifact{
			Path:    "AGENTS.md",
			Mode:    target.ModeBlock,
			Content: body,
		})
	}

	// Skills → .opencode/skills/<name>/SKILL.md.
	for _, sk := range s.Skills {
		if !sk.HasTarget(name) {
			continue
		}
		content, err := renderSkill(sk, displayName, name)
		if err != nil {
			return nil, err
		}
		dest := filepath.Join(".opencode", "skills", sk.Name)
		out = append(out, target.Artifact{
			Path:    filepath.Join(dest, "SKILL.md"),
			Mode:    target.ModeFile,
			Content: content,
		})
		assets, err := target.CopySkillAssets(sk, dest)
		if err != nil {
			return nil, err
		}
		out = append(out, assets...)
	}

	// Commands → .opencode/commands/<name>.md.
	for _, c := range s.Commands {
		if !c.HasTarget(name) {
			continue
		}
		content, err := renderSkill(c, displayName, name)
		if err != nil {
			return nil, err
		}
		out = append(out, target.Artifact{
			Path:    filepath.Join(".opencode", "commands", c.Name+".md"),
			Mode:    target.ModeFile,
			Content: content,
		})
	}

	// Agents → .opencode/agents/<name>.md.
	for _, a := range s.Agents {
		if !a.HasTarget(name) {
			continue
		}
		content, err := renderSkill(a, displayName, name)
		if err != nil {
			return nil, err
		}
		out = append(out, target.Artifact{
			Path:    filepath.Join(".opencode", "agents", a.Name+".md"),
			Mode:    target.ModeFile,
			Content: content,
		})
	}

	return out, nil
}

// renderSkill produces a markdown file with YAML frontmatter (name +
// description + per-target overrides). The same shape is reused for
// OpenCode's commands and agents.
func renderSkill(p source.Primitive, displayName, targetName string) (string, error) {
	fields := []render.Field{
		{Key: "name", Value: p.Name},
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
