package cli_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestMeta_SkillOutputsSKILLMd(t *testing.T) {
	is := is.New(t)
	stdout, _, err := invoke(t, "meta", "skill")
	is.NoErr(err)

	// Frontmatter present.
	is.True(strings.HasPrefix(stdout, "---\n"))
	is.True(strings.Contains(stdout, "name: hatch"))
	is.True(strings.Contains(stdout, "description:"))

	// Body mentions the essentials a coding agent needs to know.
	is.True(strings.Contains(stdout, "go install github.com/matryer/hatch"))
	is.True(strings.Contains(stdout, "hatch new"))
	is.True(strings.Contains(stdout, "hatch gen"))
	is.True(strings.Contains(stdout, ".hatch/"))
	is.True(strings.Contains(stdout, "rule"))
	is.True(strings.Contains(stdout, "skill"))
	is.True(strings.Contains(stdout, "command"))
	is.True(strings.Contains(stdout, "agent"))
	is.True(strings.Contains(stdout, "Never edit generated files"))
}

func TestMeta_MissingSubcommandErrors(t *testing.T) {
	is := is.New(t)
	_, _, err := invoke(t, "meta")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "missing subcommand"))
}

func TestMeta_UnknownSubcommandErrors(t *testing.T) {
	is := is.New(t)
	_, _, err := invoke(t, "meta", "nonsense")
	is.True(err != nil)
}
