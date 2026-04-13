// Package block implements the hatch-managed block-injection protocol used
// for files that may contain user-authored content around a hatch-generated
// section (CLAUDE.md, AGENTS.md, .github/copilot-instructions.md).
//
// A hatch-managed block is delimited by HTML comments:
//
//	<!-- hatch:begin v1 -->
//	...hatch-generated content...
//	<!-- hatch:end v1 -->
//
// The `v1` token versions the marker so older blocks can be recognized and
// replaced. Everything outside the markers is user-authored and is never
// touched by hatch.
package block

import (
	"fmt"
	"os"
	"strings"
)

// CurrentMarker is the marker token used by this version of hatch.
const CurrentMarker = "v1"

// Markers returns the begin/end lines for a given marker token.
func Markers(marker string) (begin, end string) {
	return fmt.Sprintf("<!-- hatch:begin %s -->", marker),
		fmt.Sprintf("<!-- hatch:end %s -->", marker)
}

// Render wraps content in begin/end markers and returns the complete block
// text (without surrounding file content).
func Render(marker, content string) string {
	begin, end := Markers(marker)
	content = strings.TrimRight(content, "\n")
	return begin + "\n" + content + "\n" + end + "\n"
}

// Inject writes the hatch block into path. If the file exists and already
// contains a hatch block (for any recognized marker), the block is replaced
// in-place. If the file exists without a block, the block is appended. If
// the file does not exist, it is created with just the block.
//
// Content that embeds the literal begin/end marker text is rejected: such
// content would collide with the block boundaries on subsequent rebuilds
// and corrupt the file. Callers that legitimately need to document hatch
// markers should escape them before passing the content to Inject.
func Inject(path, marker, content string) error {
	if err := validateContent(content); err != nil {
		return err
	}
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	block := Render(marker, content)
	var out string
	if len(existing) == 0 {
		out = block
	} else {
		out = replaceBlock(string(existing), marker, block)
	}
	if err := os.MkdirAll(parentDir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(out), 0o644)
}

// validateContent refuses content that would break block boundaries. Any
// occurrence of the literal marker prefix (without caring about the version
// token) would cause a later Inject/Strip to find the wrong boundary.
func validateContent(content string) error {
	for _, needle := range []string{"<!-- hatch:begin", "<!-- hatch:end"} {
		if strings.Contains(content, needle) {
			return fmt.Errorf("content contains hatch marker text %q; would corrupt block boundaries", needle)
		}
	}
	return nil
}

// Strip removes the hatch block (for the given marker) from path, preserving
// surrounding user content. If the file ends up empty after stripping, it is
// removed entirely.
func Strip(path, marker string) error {
	return stripFile(path, marker)
}

// replaceBlock swaps any existing hatch block in src (matching marker) with
// block. If none is found, block is appended with a separating blank line.
func replaceBlock(src, marker, block string) string {
	begin, end := Markers(marker)
	bi := strings.Index(src, begin)
	if bi < 0 {
		// No existing block — append.
		if !strings.HasSuffix(src, "\n") {
			src += "\n"
		}
		if !strings.HasSuffix(src, "\n\n") {
			src += "\n"
		}
		return src + block
	}
	ei := strings.Index(src[bi:], end)
	if ei < 0 {
		// Malformed (begin without end); leave the file alone and append.
		return src + "\n" + block
	}
	ei += bi + len(end)
	for ei < len(src) && src[ei] == '\n' {
		ei++
	}
	tail := src[ei:]
	head := src[:bi]
	var buf strings.Builder
	buf.WriteString(head)
	buf.WriteString(block)
	if tail != "" {
		if !strings.HasSuffix(block, "\n") {
			buf.WriteString("\n")
		}
		buf.WriteString(tail)
	}
	return buf.String()
}

func stripFile(path, marker string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	begin, end := Markers(marker)
	src := string(data)
	bi := strings.Index(src, begin)
	if bi < 0 {
		return nil
	}
	ei := strings.Index(src[bi:], end)
	if ei < 0 {
		return fmt.Errorf("begin marker without matching end")
	}
	ei += bi + len(end)
	for ei < len(src) && src[ei] == '\n' {
		ei++
	}
	stripped := strings.TrimLeft(src[:bi], "")
	stripped = strings.TrimRight(stripped, "\n")
	tail := strings.TrimLeft(src[ei:], "\n")
	var result string
	switch {
	case stripped == "" && tail == "":
		return os.Remove(path)
	case stripped == "":
		result = tail
	case tail == "":
		result = stripped + "\n"
	default:
		result = stripped + "\n\n" + tail
	}
	return os.WriteFile(path, []byte(result), 0o644)
}

func parentDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
