package source

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// splitFrontmatter divides a markdown file into its YAML frontmatter (between
// `---` delimiters at the start) and body. If no frontmatter is present, the
// whole input is returned as the body. CRLF line endings are normalised to
// LF before parsing so downstream emission is platform-independent.
func splitFrontmatter(data []byte) (front []byte, body []byte, err error) {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	sep := []byte("---")
	trimmed := bytes.TrimLeft(data, " \t\n")
	if !bytes.HasPrefix(trimmed, sep) {
		return nil, data, nil
	}
	rest := trimmed[len(sep):]
	if !bytes.HasPrefix(rest, []byte("\n")) {
		return nil, data, nil
	}
	rest = bytes.TrimLeft(rest, "\n")
	end := bytes.Index(rest, append([]byte("\n"), sep...))
	if end < 0 {
		return nil, nil, fmt.Errorf("unterminated frontmatter")
	}
	front = rest[:end]
	body = rest[end+1+len(sep):]
	body = bytes.TrimLeft(body, "\n")
	return front, body, nil
}

// parsePrimitive parses a raw markdown file into a Primitive. kind and name
// are supplied by the caller (derived from the filesystem walk); any frontmatter
// `name` field wins over the derived value.
func parsePrimitive(kind Kind, derivedName string, sourcePath string, data []byte) (Primitive, error) {
	p := Primitive{
		Kind:       kind,
		Name:       derivedName,
		SourcePath: sourcePath,
	}

	front, body, err := splitFrontmatter(data)
	if err != nil {
		return p, fmt.Errorf("%s: %w", sourcePath, err)
	}
	p.Body = string(body)

	if len(bytes.TrimSpace(front)) == 0 {
		if kind != KindRule && p.Name == "" {
			return p, fmt.Errorf("%s: missing name", sourcePath)
		}
		return p, nil
	}

	var root yaml.Node
	if err := yaml.Unmarshal(front, &root); err != nil {
		return p, fmt.Errorf("%s: frontmatter: %w", sourcePath, err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return p, nil
	}
	top := root.Content[0]
	if top.Kind != yaml.MappingNode {
		return p, fmt.Errorf("%s: frontmatter must be a mapping", sourcePath)
	}

	known := map[string]bool{
		"name":        true,
		"description": true,
		"applyTo":     true,
		"targets":     true,
	}
	overrides := map[string]*yaml.Node{}

	for i := 0; i < len(top.Content); i += 2 {
		k := top.Content[i]
		v := top.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		key := k.Value
		switch key {
		case "name":
			// Only override the derived name when the frontmatter supplies a
			// non-empty value — an empty `name:` should not wipe the name
			// derived from the filename.
			if v.Kind == yaml.ScalarNode && v.Value != "" {
				p.Name = v.Value
			}
		case "description":
			if v.Kind == yaml.ScalarNode {
				p.Description = v.Value
			}
		case "applyTo":
			if v.Kind == yaml.ScalarNode {
				p.ApplyTo = v.Value
			}
		case "targets":
			switch v.Kind {
			case yaml.SequenceNode:
				for _, item := range v.Content {
					if item.Kind == yaml.ScalarNode {
						p.Targets = append(p.Targets, item.Value)
					}
				}
			case yaml.ScalarNode:
				// Accept `targets: claude` as a singleton so a common typo
				// doesn't silently fall back to "all targets".
				if v.Value != "" {
					p.Targets = []string{v.Value}
				}
			}
		case "hatch":
			// reserved; ignore for now
		default:
			if !known[key] {
				overrides[key] = v
			}
		}
	}
	if len(overrides) > 0 {
		p.Overrides = overrides
	}

	if kind != KindRule {
		if p.Description == "" {
			return p, fmt.Errorf("%s: missing description (required for %s)", sourcePath, kind)
		}
		if p.Name == "" {
			return p, fmt.Errorf("%s: missing name", sourcePath)
		}
	}
	return p, nil
}
