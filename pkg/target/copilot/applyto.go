package copilot

import (
	"fmt"
	"strings"
)

// composeApplyTo combines a scope path with the user-supplied applyTo
// glob from a rule's frontmatter, producing the effective glob that
// goes into the generated `.github/instructions/<n>.instructions.md`
// file.
//
// Rules:
//   - Empty scope (root): pass userGlob through unchanged. This keeps
//     today's root-only behaviour bit-for-bit identical.
//   - Scoped, no userGlob: emit "<scopePath>/**" — every file under
//     the scope.
//   - Scoped, relative userGlob: prepend "<scopePath>/". The user
//     wrote a glob meaning "files matching X in this scope".
//   - Scoped, absolute or **/-prefixed userGlob: error. These are
//     ambiguous (does "**/*.go" mean "every Go file in the repo" or
//     "every Go file in this scope"?). The user must move the rule to
//     the root scope, or remove applyTo to inherit "<scopePath>/**".
//
// ruleName is used only for the error message, to point the user at
// the offending rule.
func composeApplyTo(scopePath, userGlob, ruleName string) (string, error) {
	if scopePath == "" {
		return userGlob, nil
	}
	if userGlob == "" {
		return scopePath + "/**", nil
	}
	if strings.HasPrefix(userGlob, "/") || strings.HasPrefix(userGlob, "**/") {
		return "", fmt.Errorf("rule %q at path %q: explicit applyTo %q is absolute or unanchored; remove applyTo to inherit the path glob, or move the rule to the root .hatch/_rules/", ruleName, scopePath, userGlob)
	}
	return scopePath + "/" + userGlob, nil
}

// scopeSlug converts a forward-slash scope path into a flat slug
// suitable for use in a filename. Used by Copilot to disambiguate
// scoped primitives that would otherwise collide at the root paths
// (which is the only place Copilot reads `.github/`).
//
// services/api -> services-api
// backend      -> backend
// ""           -> ""
func scopeSlug(scopePath string) string {
	if scopePath == "" {
		return ""
	}
	return strings.ReplaceAll(scopePath, "/", "-")
}
