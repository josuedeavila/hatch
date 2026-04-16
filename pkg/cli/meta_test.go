package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// `hatch meta skill` was removed in favour of automatic injection: every
// `hatch gen` run writes the hatch meta SKILL.md into each target's
// native skill location, so users no longer need a separate command.

func TestGen_AutoInjectsMetaSkillIntoEveryTarget(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	// A bare .hatch/ tree (no user skills) is enough — gen still writes
	// the meta skill because injection happens after Load.
	is.NoErr(os.MkdirAll(".hatch/_rules", 0o755))

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	for _, p := range []string{
		".claude/skills/hatch/SKILL.md",
		".agents/skills/hatch/SKILL.md",
		".opencode/skills/hatch/SKILL.md",
	} {
		body, err := os.ReadFile(p)
		is.NoErr(err) // every native-skills target should have the meta skill
		s := string(body)
		is.True(strings.Contains(s, "name: hatch"))
		is.True(strings.Contains(s, "description:"))
		is.True(strings.Contains(s, "hatch gen"))
	}

	// Copilot has no native skill primitive — the meta skill is inlined
	// into the .github/copilot-instructions.md block instead.
	cp, err := os.ReadFile(".github/copilot-instructions.md")
	is.NoErr(err)
	is.True(strings.Contains(string(cp), "Skill: hatch"))
	is.True(strings.Contains(string(cp), "hatch gen"))

	// Cursor has no native skill primitive either — the meta skill is
	// inlined as a .cursor/rules/skill-hatch.mdc rule file.
	cu, err := os.ReadFile(".cursor/rules/skill-hatch.mdc")
	is.NoErr(err)
	is.True(strings.Contains(string(cu), "alwaysApply: true"))
	is.True(strings.Contains(string(cu), "hatch gen"))
}

func TestGen_AutoInjectedMetaSkillOmitsSourceMetadata(t *testing.T) {
	// When hatch auto-injects the meta skill (no .hatch/_skills/hatch/
	// on disk), there IS no source file — pointing the `source:`
	// metadata at .hatch/_skills/hatch/SKILL.md would be a lie. The
	// generated-by-hatch line still appears; source is dropped.
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	is.NoErr(os.MkdirAll(".hatch/_rules", 0o755))
	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	body, err := os.ReadFile(".claude/skills/hatch/SKILL.md")
	is.NoErr(err)
	s := string(body)
	is.True(strings.Contains(s, "generated: hatch@"))
	is.True(!strings.Contains(s, "source: .hatch/_skills/hatch/SKILL.md"))
}

func TestGen_UserOverriddenMetaSkillHasSourceMetadata(t *testing.T) {
	// If the user writes their own .hatch/_skills/hatch/SKILL.md, the
	// source IS real, so the `source:` metadata points at it.
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	is.NoErr(os.MkdirAll(".hatch/_skills/hatch", 0o755))
	is.NoErr(os.WriteFile(".hatch/_skills/hatch/SKILL.md", []byte(`---
description: My custom hatch skill
---

body
`), 0o644))

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	body, err := os.ReadFile(".claude/skills/hatch/SKILL.md")
	is.NoErr(err)
	is.True(strings.Contains(string(body), "source: .hatch/_skills/hatch/SKILL.md"))
}

func TestGen_UserSkillNamedHatchOverridesInjectedMetaSkill(t *testing.T) {
	// If the user writes their own .hatch/_skills/hatch/SKILL.md, that
	// version wins — injectMetaSkill is a no-op so users can override
	// the content of the meta skill.
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	is.NoErr(os.MkdirAll(".hatch/_skills/hatch", 0o755))
	is.NoErr(os.WriteFile(".hatch/_skills/hatch/SKILL.md", []byte(`---
description: My custom hatch skill
---

CUSTOM HATCH BODY
`), 0o644))

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	body, err := os.ReadFile(".claude/skills/hatch/SKILL.md")
	is.NoErr(err)
	s := string(body)
	is.True(strings.Contains(s, "CUSTOM HATCH BODY"))
	is.True(strings.Contains(s, "My custom hatch skill"))
	// The default meta skill body must not appear.
	is.True(!strings.Contains(s, "## The four primitives"))
}

func TestClean_RemovesAutoInjectedMetaSkill(t *testing.T) {
	// Regression: `hatch clean` must remove the same meta-skill outputs
	// that `hatch gen` wrote. Before the fix, clean was calling
	// source.Load directly and never saw the injected meta skill, so
	// .claude/skills/hatch/SKILL.md etc. leaked past clean.
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	is.NoErr(os.MkdirAll(".hatch/_rules", 0o755))

	_, _, err := invoke(t, "gen")
	is.NoErr(err)
	// Sanity: meta skill outputs exist post-gen.
	for _, p := range []string{
		".claude/skills/hatch/SKILL.md",
		".agents/skills/hatch/SKILL.md",
		".opencode/skills/hatch/SKILL.md",
		".cursor/rules/skill-hatch.mdc",
	} {
		_, err := os.Stat(p)
		is.NoErr(err)
	}

	_, _, err = invoke(t, "clean")
	is.NoErr(err)
	for _, p := range []string{
		".claude/skills/hatch/SKILL.md",
		".agents/skills/hatch/SKILL.md",
		".opencode/skills/hatch/SKILL.md",
		".cursor/rules/skill-hatch.mdc",
	} {
		_, err := os.Stat(p)
		is.True(os.IsNotExist(err)) // meta skill output must be cleaned up
	}
}

func TestGen_NoHatchSkillFlagSuppressesInjection(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	is.NoErr(os.MkdirAll(".hatch/_rules", 0o755))

	_, _, err := invoke(t, "gen", "-no-hatch-skill")
	is.NoErr(err)

	// None of the meta skill outputs should exist.
	for _, p := range []string{
		".claude/skills/hatch/SKILL.md",
		".agents/skills/hatch/SKILL.md",
		".opencode/skills/hatch/SKILL.md",
		".cursor/rules/skill-hatch.mdc",
	} {
		_, err := os.Stat(p)
		is.True(os.IsNotExist(err))
	}
	// And the Copilot block must not contain a "Skill: hatch" section.
	cp, err := os.ReadFile(".github/copilot-instructions.md")
	// File may or may not exist depending on whether any other content
	// was produced; only assert the absence of the meta skill heading
	// if the file exists.
	if err == nil {
		is.True(!strings.Contains(string(cp), "Skill: hatch"))
	}
}

func TestList_NoHatchSkillFlagHidesInjection(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	is.NoErr(os.MkdirAll(".hatch/_rules", 0o755))

	stdout, _, err := invoke(t, "list", "-no-hatch-skill")
	is.NoErr(err)
	is.True(!strings.Contains(stdout, "skills/hatch"))
	is.True(!strings.Contains(stdout, "skill-hatch"))
}

func TestClean_NoHatchSkillFlagLeavesInjectedFilesAlone(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	is.NoErr(os.MkdirAll(".hatch/_rules", 0o755))

	// First gen WITH injection so the meta skill files exist.
	_, _, err := invoke(t, "gen")
	is.NoErr(err)
	// Sanity: meta skill exists post-gen.
	_, err = os.Stat(".claude/skills/hatch/SKILL.md")
	is.NoErr(err)

	// Now clean WITH the flag — meta skill outputs must survive because
	// clean's view of the source no longer includes them.
	_, _, err = invoke(t, "clean", "-no-hatch-skill")
	is.NoErr(err)
	for _, p := range []string{
		".claude/skills/hatch/SKILL.md",
		".agents/skills/hatch/SKILL.md",
		".opencode/skills/hatch/SKILL.md",
		".cursor/rules/skill-hatch.mdc",
	} {
		_, err := os.Stat(p)
		is.NoErr(err)
	}
}

func TestRun_MetaCommandRemoved(t *testing.T) {
	is := is.New(t)
	_, _, err := invoke(t, "meta")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "unknown command"))
}
