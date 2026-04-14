package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
)

// cmdMeta dispatches the `hatch meta` family of subcommands. These emit a
// self-describing document that teaches a coding agent about hatch itself.
//
// Currently supported:
//
//	hatch meta skill [-targets list]
//
// With no targets, the SKILL.md is written to stdout. With `-targets`,
// hatch synthesises a Source containing a single skill (the meta skill)
// and runs each named target's Generate, writing the result to the
// target's native skill location (e.g. `.claude/skills/hatch/SKILL.md`
// for claude, `.agents/skills/hatch/SKILL.md` for codex, etc.).
func cmdMeta(ctx context.Context, available *target.Set, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("hatch meta: missing subcommand (want: skill)")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "skill":
		return cmdMetaSkill(ctx, available, rest, stdout, stderr)
	default:
		return fmt.Errorf("hatch meta: unknown subcommand %q (want: skill)", sub)
	}
}

func cmdMetaSkill(_ context.Context, available *target.Set, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("meta skill", flag.ContinueOnError)
	fs.SetOutput(stderr)
	targetsList := fs.String("targets", "", "comma-separated target names; omit to print to stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// No targets → print the full SKILL.md to stdout (pipe-friendly).
	if *targetsList == "" {
		_, err := io.WriteString(stdout, metaSkillDoc)
		return err
	}

	targets, err := selectTargets(available, *targetsList)
	if err != nil {
		return err
	}

	// Synthesise a source tree containing just the meta skill, then let
	// each target's Generate drop it into its native skill location.
	src := &source.Source{
		Skills: []source.Primitive{{
			Kind:        source.KindSkill,
			Name:        metaSkillName,
			Description: metaSkillDescription,
			Body:        metaSkillBody,
		}},
	}
	for _, t := range targets.All() {
		arts, err := t.Generate(src)
		if err != nil {
			return fmt.Errorf("%s: %w", t.Name(), err)
		}
		for _, a := range arts {
			if err := writeArtifact(a); err != nil {
				return fmt.Errorf("%s: %s: %w", t.Name(), a.Path, err)
			}
			fmt.Fprintf(stdout, "wrote %s (%s)\n", a.Path, a.Mode)
		}
	}
	return nil
}

// metaSkillName, metaSkillDescription, metaSkillBody are the three pieces
// that together make up the SKILL.md printed or written for `hatch meta
// skill`. Keeping them separate lets us either serialise the full file
// (for stdout) or hand them to a target's Generate (for target-native
// placement) without reparsing.
const (
	metaSkillName        = "hatch"
	metaSkillDescription = "Authoring rules, skills, commands, and sub-agents for this project via hatch — write once, generate for every coding agent."
)

const metaSkillBody = `# hatch

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
hatch init                        # scaffold .hatch/ with one example of each primitive
hatch new <kind> <title>          # create a new source file
hatch gen                         # regenerate all target files
hatch gen -targets claude         # regenerate only one agent's files
hatch list                        # dry-run: show what gen would write
hatch clean                       # remove everything hatch generated
hatch meta skill -targets claude  # drop this SKILL.md into every target's skills dir
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

// metaSkillDoc is the full SKILL.md (frontmatter + body) emitted for
// `hatch meta skill` with no targets. Assembled from the constants above
// at package init so there's a single source of truth for the body.
const metaSkillDoc = "---\n" +
	"name: " + metaSkillName + "\n" +
	"description: " + metaSkillDescription + "\n" +
	"---\n\n" +
	metaSkillBody
