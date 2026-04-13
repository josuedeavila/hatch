package render

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Field is one key/value pair to render in frontmatter. Value may be a
// scalar string, a []string, or a *yaml.Node (for passthrough from the
// source).
type Field struct {
	Key   string
	Value any
}

// Frontmatter renders a YAML frontmatter block (between `---` delimiters)
// for the given fields, in the order listed. Keys with nil or empty values
// are omitted. The output is deterministic and uses LF line endings.
//
// If fields is empty, returns an empty string — callers should write no
// frontmatter at all in that case.
func Frontmatter(fields []Field) (string, error) {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	for _, f := range fields {
		val, ok := toNode(f.Value)
		if !ok {
			continue
		}
		mapping.Content = append(mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: f.Key},
			val,
		)
	}
	if len(mapping.Content) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(mapping); err != nil {
		return "", fmt.Errorf("encode frontmatter: %w", err)
	}
	enc.Close()

	body := strings.TrimRight(buf.String(), "\n")
	return "---\n" + body + "\n---\n", nil
}

// MergeOverride merges an extra per-target mapping into the field list. Any
// keys already present in fields are replaced; the rest are appended in
// alphabetical order after the named fields.
func MergeOverride(fields []Field, override *yaml.Node) []Field {
	if override == nil || override.Kind != yaml.MappingNode {
		return fields
	}
	known := map[string]int{}
	for i, f := range fields {
		known[f.Key] = i
	}
	var extras []Field
	for i := 0; i < len(override.Content); i += 2 {
		k := override.Content[i]
		v := override.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		if idx, ok := known[k.Value]; ok {
			fields[idx].Value = v
			continue
		}
		extras = append(extras, Field{Key: k.Value, Value: v})
	}
	sort.SliceStable(extras, func(i, j int) bool { return extras[i].Key < extras[j].Key })
	return append(fields, extras...)
}

// toNode converts v to a yaml.Node suitable for rendering. Returns
// (nil, false) when the value is empty and should be omitted.
func toNode(v any) (*yaml.Node, bool) {
	switch t := v.(type) {
	case nil:
		return nil, false
	case string:
		if t == "" {
			return nil, false
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Value: t}, true
	case []string:
		if len(t) == 0 {
			return nil, false
		}
		seq := &yaml.Node{Kind: yaml.SequenceNode, Style: yaml.FlowStyle}
		for _, s := range t {
			seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: s})
		}
		return seq, true
	case *yaml.Node:
		if t == nil {
			return nil, false
		}
		return t, true
	default:
		return nil, false
	}
}
