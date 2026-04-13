package source

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load walks `<root>/.hatch/` and returns a populated Source. Source
// primitives live in subdirectories (rules/, skills/, commands/, agents/)
// directly under `.hatch/`; any other files (e.g. config.yaml) are ignored.
func Load(root string) (*Source, error) {
	srcRoot := filepath.Join(root, ".hatch")
	info, err := os.Stat(srcRoot)
	if err != nil {
		return nil, fmt.Errorf("source directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", srcRoot)
	}

	s := &Source{}

	if err := loadFlatDir(s, filepath.Join(srcRoot, "rules"), KindRule); err != nil {
		return nil, err
	}
	if err := loadFlatDir(s, filepath.Join(srcRoot, "commands"), KindCommand); err != nil {
		return nil, err
	}
	if err := loadFlatDir(s, filepath.Join(srcRoot, "agents"), KindAgent); err != nil {
		return nil, err
	}
	if err := loadSkillsDir(s, filepath.Join(srcRoot, "skills")); err != nil {
		return nil, err
	}

	sortByName(s.Rules)
	sortByName(s.Skills)
	sortByName(s.Commands)
	sortByName(s.Agents)
	return s, nil
}

func loadFlatDir(src *Source, dir string, kind Kind) error {
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
			src.Rules = append(src.Rules, p)
		case KindCommand:
			src.Commands = append(src.Commands, p)
		case KindAgent:
			src.Agents = append(src.Agents, p)
		}
	}
	return nil
}

func loadSkillsDir(src *Source, dir string) error {
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
		src.Skills = append(src.Skills, p)
	}
	return nil
}

func sortByName(ps []Primitive) {
	sort.SliceStable(ps, func(i, j int) bool { return ps[i].Name < ps[j].Name })
}
