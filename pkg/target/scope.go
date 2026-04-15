package target

import "path/filepath"

// ScopedPath joins parts with the OS separator, prefixed by scopePath
// (which is stored forward-slash regardless of OS). An empty scopePath
// is the no-op fast path used by the root scope, so this helper is a
// drop-in replacement for `filepath.Join` everywhere a target builds an
// output path: targets pass `sc.Path` from the scope they're generating
// for and get either a root-relative or scope-prefixed path back.
func ScopedPath(scopePath string, parts ...string) string {
	if scopePath == "" {
		return filepath.Join(parts...)
	}
	return filepath.Join(filepath.FromSlash(scopePath), filepath.Join(parts...))
}
