// Package source loads hatch source files from a `.hatch/` tree and parses
// them into primitives (rules, skills, commands, agents).
package source

import "gopkg.in/yaml.v3"

// Kind identifies a primitive type.
type Kind string

const (
	KindRule    Kind = "rule"
	KindSkill   Kind = "skill"
	KindCommand Kind = "command"
	KindAgent   Kind = "agent"
)

// Primitive is one hatch source file, parsed.
//
// For skills, SourcePath points to the skill directory (containing SKILL.md
// plus any sibling assets). For other kinds it points to the source .md file.
type Primitive struct {
	Kind        Kind
	Name        string
	Description string
	ApplyTo     string
	Targets     []string
	Body        string
	SourcePath  string
	Overrides   map[string]*yaml.Node
}

// HasTarget reports whether the primitive should be generated for the
// named target. An empty Targets slice means "all targets".
func (p *Primitive) HasTarget(name string) bool {
	if len(p.Targets) == 0 {
		return true
	}
	for _, t := range p.Targets {
		if t == name {
			return true
		}
	}
	return false
}

// Source is the loaded set of hatch primitives for a project.
type Source struct {
	Rules    []Primitive
	Skills   []Primitive
	Commands []Primitive
	Agents   []Primitive
}
