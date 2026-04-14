package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
)

// cmdMeta dispatches the `hatch meta` family of subcommands. These write
// self-describing documents to stdout so users can pipe them into a skills
// directory, a docs file, or the clipboard.
//
// Currently supported:
//
//	hatch meta skill    # print a SKILL.md teaching a coding agent about hatch
func cmdMeta(_ context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("hatch meta: missing subcommand (want: skill)")
	}
	sub := args[0]
	switch sub {
	case "skill":
		_, err := io.WriteString(stdout, metaSkillDoc)
		return err
	default:
		return fmt.Errorf("hatch meta: unknown subcommand %q (want: skill)", sub)
	}
}

// metaSkillDoc is an agentskills.io-compatible SKILL.md that teaches a
// coding agent how to use hatch in this project. It's intentionally
// comprehensive: installing, the four primitives, source layout, the
// `hatch new` helper, frontmatter fields, and the edit-regen workflow.
const metaSkillDoc = `---
name: hatch
description: Authoring rules, skills, commands, and sub-agents for this project via hatch — write once, generate for every coding agent.
---

# hatch

This project uses **hatch** (` + "`github.com/matryer/hatch`" + `) to keep a single
source of truth for the guidance it sends to coding agents. Hatch reads a
directory under ` + "`.hatch/`" + ` and produces the native files each agent expects
(Claude Code, OpenAI Codex, GitHub Copilot, OpenCode).

When you are asked to add, change, or remove rules, skills, slash commands,
or sub-agent definitions in this project, edit files under ` + "`.hatch/`" + ` — not
the generated files in ` + "`CLAUDE.md`" + `, ` + "`AGENTS.md`" + `, ` + "`.claude/`" + `, ` + "`.agents/`" + `,
` + "`.github/`" + `, or ` + "`.opencode/`" + `. After editing, run ` + "`hatch gen`" + ` from the project
root to regenerate everything.

## Install

` + "```" + `
go install github.com/matryer/hatch/cmd/hatch@latest
` + "```" + `

## The four primitives

| Kind      | Purpose                                                            | Source path                           |
|-----------|--------------------------------------------------------------------|---------------------------------------|
| ` + "`rule`" + `    | Always-on project instructions; optionally scoped with a glob      | ` + "`.hatch/rules/<slug>.md`" + `              |
| ` + "`skill`" + `   | Model-invoked in-session capability; supports sibling assets       | ` + "`.hatch/skills/<slug>/SKILL.md`" + `       |
| ` + "`command`" + ` | User-invoked slash prompt                                          | ` + "`.hatch/commands/<slug>.md`" + `           |
| ` + "`agent`" + `   | Delegated sub-agent definition                                     | ` + "`.hatch/agents/<slug>.md`" + `             |

## Creating a new source file

Use ` + "`hatch new`" + `:

` + "```" + `
hatch new rule "Always run gofmt before committing"
hatch new skill "Review pull requests"
hatch new command "Commit with generated message"
hatch new agent "Security auditor"
` + "```" + `

Each call creates the right file under ` + "`.hatch/`" + ` with a minimal template,
slugs the title into a filesystem-safe name, and reminds you to run
` + "`hatch gen`" + ` afterwards. For skills, the file is ` + "`SKILL.md`" + ` inside a
directory — place sibling assets (scripts, references) alongside and they
copy through verbatim.

You can also write files by hand under ` + "`.hatch/<kind>/...`" + ` if you prefer —
` + "`hatch new`" + ` is just a scaffolding helper.

## Frontmatter

Skills, commands, and agents carry a YAML frontmatter header. Only
` + "`description`" + ` is required.

| Field         | Required | Meaning                                                                                              |
|---------------|----------|------------------------------------------------------------------------------------------------------|
| ` + "`description`" + ` | yes      | One-sentence summary. Downstream agents either require or recommend this.                            |
| ` + "`name`" + `        | no       | Overrides the slug. Defaults to the source filename/dirname.                                         |
| ` + "`targets`" + `     | no       | List of target names (` + "`claude`" + `, ` + "`codex`" + `, ` + "`copilot`" + `, ` + "`opencode`" + `). Empty/absent → all.                 |
| ` + "`applyTo`" + `     | no       | Rules only. Glob pattern scoping the rule to matching files.                                         |
| ` + "`claude:`" + ` etc | no       | Per-target passthrough block — its contents merge into the generated file's own frontmatter.        |

Rule files are plain markdown and don't need frontmatter at all. Bodies may
use ` + "`{{agent}}`" + ` and ` + "`{{target}}`" + ` — these are substituted at generation time
with the agent's display name (e.g., "Claude Code") and short name (e.g.,
"claude").

## Workflow

1. ` + "`cd`" + ` to the project root.
2. Edit or create files under ` + "`.hatch/`" + ` — either by hand or with ` + "`hatch new`" + `.
3. Run ` + "`hatch gen`" + ` to regenerate every agent's files.
4. Commit both the ` + "`.hatch/`" + ` source changes and the regenerated output
   files together.

## Useful commands

` + "```" + `
hatch init                   # scaffold .hatch/ with one example of each primitive
hatch new <kind> <title>     # create a new source file
hatch gen                    # regenerate all target files
hatch gen claude             # regenerate only claude (positional target arg)
hatch gen claude codex       # regenerate claude and codex
hatch list                   # dry-run: show what gen would write
hatch clean                  # remove everything hatch generated
` + "```" + `

## Never edit generated files

` + "`CLAUDE.md`" + `, ` + "`AGENTS.md`" + `, ` + "`.github/copilot-instructions.md`" + `, and everything
under ` + "`.claude/`" + `, ` + "`.agents/`" + `, ` + "`.opencode/`" + `, ` + "`.github/prompts/`" + `,
` + "`.github/instructions/`" + `, and ` + "`.github/agents/`" + ` is regenerated by hatch. Any
edits you make there will be overwritten on the next ` + "`hatch gen`" + `. Put your
changes in ` + "`.hatch/`" + ` instead.

For the block-injected files (` + "`CLAUDE.md`" + `, ` + "`AGENTS.md`" + `,
` + "`.github/copilot-instructions.md`" + `), hatch only rewrites content between
` + "`<!-- hatch:begin v1 -->`" + ` and ` + "`<!-- hatch:end v1 -->`" + ` markers — surrounding
content you have written by hand is preserved across regeneration.
`
