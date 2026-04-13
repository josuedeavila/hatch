package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestNew_RuleWithTitleArg(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	stdout, _, err := invoke(t, "new", "rule", "My Rule")
	is.NoErr(err)

	path := ".hatch/rules/my-rule.md"
	_, err = os.Stat(path)
	is.NoErr(err)

	body, err := os.ReadFile(path)
	is.NoErr(err)
	is.True(strings.Contains(string(body), "# My Rule"))

	// Output should confirm creation and remind the user to run gen.
	is.True(strings.Contains(stdout, "created"))
	is.True(strings.Contains(stdout, "my-rule.md"))
	is.True(strings.Contains(stdout, "hatch gen"))
}

func TestNew_SkillCreatesDirectoryWithSKILLMd(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "new", "skill", "Review PR")
	is.NoErr(err)

	body, err := os.ReadFile(".hatch/skills/review-pr/SKILL.md")
	is.NoErr(err)
	is.True(strings.Contains(string(body), "description:"))
	is.True(strings.Contains(string(body), "# Review PR"))
}

func TestNew_CommandFile(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "new", "command", "Commit Changes")
	is.NoErr(err)
	_, err = os.Stat(".hatch/commands/commit-changes.md")
	is.NoErr(err)
}

func TestNew_AgentFile(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "new", "agent", "Security Auditor")
	is.NoErr(err)
	_, err = os.Stat(".hatch/agents/security-auditor.md")
	is.NoErr(err)
}

func TestNew_MultiWordTitleAsSeparateArgs(t *testing.T) {
	// `hatch new rule my new rule` should join the positional args with
	// spaces so quoting is optional.
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "new", "rule", "my", "new", "rule")
	is.NoErr(err)
	_, err = os.Stat(".hatch/rules/my-new-rule.md")
	is.NoErr(err)
}

func TestNew_PromptsForTitleWhenMissing(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	stdout, _, err := invokeWithStdin(t, "My Prompted Rule\n", "new", "rule")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "rule name:"))
	_, err = os.Stat(".hatch/rules/my-prompted-rule.md")
	is.NoErr(err)
}

func TestNew_RefusesToOverwriteExisting(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "new", "rule", "dup")
	is.NoErr(err)
	_, _, err = invoke(t, "new", "rule", "dup")
	is.True(err != nil) // second call must fail
	is.True(strings.Contains(err.Error(), "already exists"))
}

func TestNew_UnknownKindErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	_, _, err := invoke(t, "new", "nope", "x")
	is.True(err != nil)
}

func TestNew_MissingKindErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	_, _, err := invoke(t, "new")
	is.True(err != nil)
}

func TestNew_EmptyTitleErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	_, _, err := invoke(t, "new", "rule", "   ")
	is.True(err != nil)
}

func TestNew_PunctuationOnlyTitleErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	_, _, err := invoke(t, "new", "rule", "!!!")
	is.True(err != nil)
}
