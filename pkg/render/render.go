// Package render builds the final text content for a hatch emission: YAML
// frontmatter (deterministically ordered) plus a body with template variables
// substituted.
package render

import "strings"

// Body substitutes the tiny template vars {{agent}} and {{target}} in s.
// {{agent}} is the human display name (e.g. "Claude Code"); {{target}} is the
// short machine name (e.g. "claude").
func Body(s, agentDisplay, targetName string) string {
	s = strings.ReplaceAll(s, "{{agent}}", agentDisplay)
	s = strings.ReplaceAll(s, "{{target}}", targetName)
	return s
}
