package copilot

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestComposeApplyTo(t *testing.T) {
	cases := []struct {
		name      string
		scope     string
		userGlob  string
		want      string
		wantError string // substring, "" means no error
	}{
		{
			name:     "root no glob — empty result, root has no implicit prefix",
			scope:    "",
			userGlob: "",
			want:     "",
		},
		{
			name:     "root with glob — passes through unchanged",
			scope:    "",
			userGlob: "**/*.go",
			want:     "**/*.go",
		},
		{
			name:     "scoped no glob — implicit <path>/** glob",
			scope:    "backend",
			userGlob: "",
			want:     "backend/**",
		},
		{
			name:     "scoped with relative glob — prepended",
			scope:    "backend",
			userGlob: "*.go",
			want:     "backend/*.go",
		},
		{
			name:     "deep scope with glob — prepended",
			scope:    "services/api",
			userGlob: "handlers/*.go",
			want:     "services/api/handlers/*.go",
		},
		{
			name:      "scoped + absolute glob — error",
			scope:     "backend",
			userGlob:  "/etc/foo",
			wantError: "absolute or unanchored",
		},
		{
			name:      "scoped + **/-prefixed glob — error",
			scope:     "backend",
			userGlob:  "**/foo.go",
			wantError: "absolute or unanchored",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			got, err := composeApplyTo(tc.scope, tc.userGlob, "myrule")
			if tc.wantError != "" {
				is.True(err != nil)
				is.True(strings.Contains(err.Error(), tc.wantError))
				is.True(strings.Contains(err.Error(), "myrule"))
				return
			}
			is.NoErr(err)
			is.Equal(got, tc.want)
		})
	}
}
