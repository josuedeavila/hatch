package cli_test

import (
	"os"
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

func TestMeta_SkillWritesToClaudeSkillsDir(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	stdout, _, err := invoke(t, "meta", "skill", "-targets", "claude")
	is.NoErr(err)
	is.True(strings.Contains(stdout, ".claude/skills/hatch/SKILL.md"))

	body, err := os.ReadFile(".claude/skills/hatch/SKILL.md")
	is.NoErr(err)
	s := string(body)
	is.True(strings.HasPrefix(s, "---\n"))
	is.True(strings.Contains(s, "name: hatch"))
	is.True(strings.Contains(s, "description:"))
	is.True(strings.Contains(s, "hatch gen"))
	is.True(strings.Contains(s, "Never edit generated files"))
}

func TestMeta_SkillWritesToCodexAgentskillsPath(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "meta", "skill", "-targets", "codex")
	is.NoErr(err)
	// Codex uses the agentskills.io standard path.
	_, err = os.Stat(".agents/skills/hatch/SKILL.md")
	is.NoErr(err)
}

func TestMeta_SkillWritesToOpencodeSkillsDir(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "meta", "skill", "-targets", "opencode")
	is.NoErr(err)
	_, err = os.Stat(".opencode/skills/hatch/SKILL.md")
	is.NoErr(err)
}

func TestMeta_SkillMultipleTargets(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "meta", "skill", "-targets", "claude,codex,opencode")
	is.NoErr(err)
	for _, path := range []string{
		".claude/skills/hatch/SKILL.md",
		".agents/skills/hatch/SKILL.md",
		".opencode/skills/hatch/SKILL.md",
	} {
		_, err := os.Stat(path)
		is.NoErr(err)
	}
}

func TestMeta_SkillCopilotInlinesIntoInstructions(t *testing.T) {
	// Copilot has no native skill primitive; hatch inlines skill bodies
	// into .github/copilot-instructions.md. Meta skill should reach
	// Copilot via that same path.
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	_, _, err := invoke(t, "meta", "skill", "-targets", "copilot")
	is.NoErr(err)

	body, err := os.ReadFile(".github/copilot-instructions.md")
	is.NoErr(err)
	s := string(body)
	is.True(strings.Contains(s, "<!-- hatch:begin v1 -->"))
	is.True(strings.Contains(s, "hatch"))
	is.True(strings.Contains(s, "Never edit generated files"))
}

func TestMeta_SkillUnknownTargetErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	_, _, err := invoke(t, "meta", "skill", "-targets", "nosuch")
	is.True(err != nil)
}
