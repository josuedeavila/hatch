package cli

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestSlugify(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"My rule", "my-rule"},
		{"MY RULE", "my-rule"},
		{"my---rule", "my-rule"},
		{"My Rule!", "my-rule"},
		{"  leading/trailing  ", "leading-trailing"},
		{"Multi   Space   Title", "multi-space-title"},
		{"with_underscores_too", "with-underscores-too"},
		{"#1 cool-thing", "1-cool-thing"},
		{"", ""},
		{"!!!", ""},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			is := is.New(t)
			is.Equal(slugify(tc.in), tc.want)
		})
	}
}

func TestSlugify_TruncatesLongTitle(t *testing.T) {
	is := is.New(t)
	long := strings.Repeat("word ", 40) // 200 chars of "word " repeated
	got := slugify(long)
	is.True(len(got) <= maxSlugLength)
	// Must not end with a trailing hyphen after truncation.
	is.True(!strings.HasSuffix(got, "-"))
}
