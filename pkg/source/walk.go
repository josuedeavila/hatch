package source

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load walks `<root>/.hatch/` and returns a populated Source. Each
// directory under .hatch/ that contains at least one of the four
// primitive container subdirectories (_rules/, _skills/, _commands/,
// _agents/) becomes a Scope; the path between .hatch/ and the
// primitive container becomes the Scope's Path (forward-slash joined,
// "" for the root). The returned Source always contains at least the
// root scope, even if empty.
func Load(root string) (*Source, error) {
	srcRoot := filepath.Join(root, ".hatch")
	info, err := os.Stat(srcRoot)
	if err != nil {
		return nil, fmt.Errorf("source directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", srcRoot)
	}

	s := &Source{Scopes: []Scope{{Path: ""}}}

	// Always load the root scope's primitive containers up front, so the
	// root scope is populated even if it has no nested siblings.
	if err := loadScope(&s.Scopes[0], srcRoot); err != nil {
		return nil, err
	}

	// Walk the rest of .hatch/ looking for nested scopes. A nested scope
	// is any directory that contains at least one of the four primitive
	// container subdirectories.
	err = filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == srcRoot {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		// The four known primitive container names are loaded by their
		// parent scope's loadScope call. Never recurse inside them — this
		// also prevents pathological layouts like _rules/nested/_rules.
		if name == RulesDir || name == SkillsDir || name == CommandsDir || name == AgentsDir {
			return filepath.SkipDir
		}
		// Hidden directories (`.git`, `.cache`, etc.) are ignored.
		if strings.HasPrefix(name, ".") {
			return filepath.SkipDir
		}

		rel, relErr := filepath.Rel(srcRoot, path)
		if relErr != nil {
			return relErr
		}
		scopePath := filepath.ToSlash(rel)

		nested := Scope{Path: scopePath}
		if err := loadScope(&nested, path); err != nil {
			return err
		}
		// Only register the scope if it actually loaded any primitives.
		// Pure passthrough containers (e.g. .hatch/services/ when only
		// services/api/ and services/web/ have primitives) are dropped.
		if len(nested.Rules) > 0 || len(nested.Skills) > 0 || len(nested.Commands) > 0 || len(nested.Agents) > 0 {
			s.Scopes = append(s.Scopes, nested)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort: root scope first, then lexicographic by Path. Within each
	// scope, sort each primitive slice by name.
	sort.SliceStable(s.Scopes, func(i, j int) bool {
		if s.Scopes[i].Path == "" {
			return true
		}
		if s.Scopes[j].Path == "" {
			return false
		}
		return s.Scopes[i].Path < s.Scopes[j].Path
	})
	for i := range s.Scopes {
		sortByName(s.Scopes[i].Rules)
		sortByName(s.Scopes[i].Skills)
		sortByName(s.Scopes[i].Commands)
		sortByName(s.Scopes[i].Agents)
	}
	return s, nil
}

// loadScope loads any primitive containers found directly under dir into
// sc. Missing containers are silently skipped (they're optional).
func loadScope(sc *Scope, dir string) error {
	if err := loadFlatDir(sc, filepath.Join(dir, RulesDir), KindRule); err != nil {
		return err
	}
	if err := loadFlatDir(sc, filepath.Join(dir, CommandsDir), KindCommand); err != nil {
		return err
	}
	if err := loadFlatDir(sc, filepath.Join(dir, AgentsDir), KindAgent); err != nil {
		return err
	}
	return loadSkillsDir(sc, filepath.Join(dir, SkillsDir))
}

func loadFlatDir(sc *Scope, dir string, kind Kind) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		derived := strings.TrimSuffix(e.Name(), ".md")
		p, err := parsePrimitive(kind, derived, path, data)
		if err != nil {
			return err
		}
		if p.Name == "" {
			p.Name = derived
		}
		switch kind {
		case KindRule:
			sc.Rules = append(sc.Rules, p)
		case KindCommand:
			sc.Commands = append(sc.Commands, p)
		case KindAgent:
			sc.Agents = append(sc.Agents, p)
		}
	}
	return nil
}

func loadSkillsDir(sc *Scope, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillDir := filepath.Join(dir, e.Name())
		skillFile := filepath.Join(skillDir, "SKILL.md")
		data, err := os.ReadFile(skillFile)
		if err != nil {
			return fmt.Errorf("skill %q: %w", e.Name(), err)
		}
		p, err := parsePrimitive(KindSkill, e.Name(), skillDir, data)
		if err != nil {
			return err
		}
		if p.Name == "" {
			p.Name = e.Name()
		}
		sc.Skills = append(sc.Skills, p)
	}
	return nil
}

func sortByName(ps []Primitive) {
	sort.SliceStable(ps, func(i, j int) bool { return ps[i].Name < ps[j].Name })
}
