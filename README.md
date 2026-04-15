# 🥚 hatch

Write rules, skills, commands, and sub-agent definitions **once**, generate
the native files each coding agent expects.

Hatch reads a single source tree under `.hatch/` and produces the specific
files Claude Code, OpenAI Codex CLI, GitHub Copilot, Cursor, and OpenCode
each read to customise their behaviour.

## Install

```
go install github.com/matryer/hatch/cmd/hatch@latest
```

Or with mise inside the repo:

```
mise run install
```

## Quick start

A quick end-to-end walkthrough, starting from an empty project directory.

**1. Scaffold a source tree.** From the root of your project:

```
$ cd my-project
$ hatch init
```

That creates an empty `.hatch/` tree with the four primitive container
subdirectories (`_rules/`, `_skills/`, `_commands/`, `_agents/`) ready
for you to drop sources into. The leading underscore tells hatch (and
human readers) "this is a primitive container, not a path component" —
it matters once you start nesting (see "Nesting for monorepos" below).

> Prefer to start from a working example of each kind? Use `hatch init -examples` instead — it additionally writes a sample rule, skill, command, and agent so you can `hatch gen` immediately and see output.

**2. Add your own source file.** Drop in a rule you want every coding agent
in this project to follow:

```
$ hatch new rule "Always write tests"
created .hatch/_rules/always-write-tests.md
edit the file, then run `hatch gen` to write the output files.
```

Open `.hatch/_rules/always-write-tests.md` in your editor and replace the
placeholder body with the rule you want. Rules are plain markdown — no
frontmatter required.

**3. Preview what will be written.** Before you touch any generated files,
dry-run to see where everything will land:

```
$ hatch list
Claude Code (claude):
  CLAUDE.md  [block]
  .claude/skills/review-pr/SKILL.md  [file]
  ...
OpenAI Codex (codex):
  AGENTS.md  [block]
  ...
```

**4. Generate the native files.** When it looks right:

```
$ hatch gen
wrote CLAUDE.md (lines 1-3)
wrote .claude/skills/review-pr/SKILL.md (file)
wrote AGENTS.md (lines 1-3)
...
```

Block-injected files report the line range the hatch block now occupies
so you can jump straight to it in your editor; file-owned outputs show
`(file)` instead.

Your always-write-tests rule is now in `CLAUDE.md`, `AGENTS.md`, and
`.github/copilot-instructions.md` — each in the form that coding agent
expects — alongside any skills, commands, and sub-agents you had. Commit
the `.hatch/` source *and* the generated files together.

**5. Iterate.** Edit anything under `.hatch/`, run `hatch gen` again, and
every agent's files update in lockstep. `hatch clean` removes everything
hatch wrote if you ever want to start over — your `.hatch/` source tree is
never touched.

## Using the CLI

All subcommands operate on the current working directory — `cd` into your
project first. Every command is safe to re-run.

### `hatch init [-examples] [-path <relative-path>]`

Scaffolds a `.hatch/` source tree with the four primitive container
subdirectories (`_rules/`, `_skills/`, `_commands/`, `_agents/`) so you
have a home for your source files. By default the directories are empty.
Pass `-examples` to additionally write one working example of each
primitive — handy for a first-time setup where you want to `hatch gen`
and see real output right away. Pass `-path` to scaffold under a nested
scope (see "Nesting for monorepos" below).

```
hatch init                    # empty .hatch/ scaffold
hatch init -examples          # + one sample rule, skill, command, and agent
hatch init -path backend      # → .hatch/backend/_rules/, _skills/, _commands/, _agents/
```

Existing files are left alone, so all forms are safe to re-run in a
directory that's already partially set up.

### `hatch new <kind> [-path <relative-path>] [title]`

Creates a single new source file from a template. `kind` is one of `rule`,
`skill`, `command`, or `agent`. The title is slugged into a
filesystem-safe name; if you omit it, you'll be prompted on stdin. Pass
`-path` to drop the new file under a nested scope.

```
hatch new rule "Always write tests"           # → .hatch/_rules/always-write-tests.md
hatch new skill "review pr"                   # → .hatch/_skills/review-pr/SKILL.md
hatch new command commit                      # → .hatch/_commands/commit.md
hatch new agent "security auditor"            # → .hatch/_agents/security-auditor.md
hatch new rule -path backend "Database style" # → .hatch/backend/_rules/database-style.md
```

Skill, command, and agent templates include a `description:` frontmatter
field pre-filled with a `TODO` — fill it in before running `hatch gen`,
since downstream agents either require or recommend it.

### `hatch gen [-targets list]`

Reads `.hatch/` and writes the native files each coding agent expects.
With no flags, every registered target is generated. Narrow the run with
`-targets`, a comma-separated list of target names:

```
hatch gen                         # every target
hatch gen -targets claude         # only Claude Code
hatch gen -targets claude,codex   # Claude Code and Codex
```

File-owned outputs under `.claude/`, `.agents/`, `.github/`, `.cursor/`,
and `.opencode/` are overwritten from scratch. Block-injected files like
`CLAUDE.md`, `AGENTS.md`, and `.github/copilot-instructions.md` have only
the hatch-managed block replaced — any content you wrote around it is
preserved across regeneration.

### `hatch list [-targets list]`

Dry-run: prints every path `hatch gen` would write, without touching the
filesystem. Takes the same `-targets` flag as `gen`. Useful for previewing
changes before committing, or in CI to assert the checked-in generated
files match the current `.hatch/` source.

```
hatch list
hatch list -targets claude
```

### `hatch clean [-targets list]`

Removes everything hatch generated. File-owned outputs are deleted; the
hatch-managed block is stripped from block-injected files (and the file
itself is deleted if it becomes empty). Your `.hatch/` source tree is
never touched.

```
hatch clean                       # remove every target's output
hatch clean -targets claude       # only Claude Code's output
```

### Auto-injected meta skill

Every `hatch gen` run automatically writes a `hatch` SKILL.md into each
target's native skill location (`.claude/skills/hatch/SKILL.md`,
`.agents/skills/hatch/SKILL.md`, `.opencode/skills/hatch/SKILL.md`,
`.cursor/rules/skill-hatch.mdc`, and an inlined section in
`.github/copilot-instructions.md`). The meta skill teaches the coding
agent how `.hatch/` is structured, so it can extend your source tree on
its own instead of editing the generated files by mistake.

Two ways to customise:

- **Override the content.** Write your own `.hatch/_skills/hatch/SKILL.md` —
  the user version wins.
- **Opt out entirely.** Pass `-no-hatch-skill` to `hatch gen`, `hatch list`,
  or `hatch clean` to skip the auto-injection for that run. For the opt-out
  to stick across regenerations, pass it consistently (and to `clean` too,
  otherwise clean will remove meta skill files that never got written).

### `hatch version`, `hatch help`

Print the version or the built-in usage summary. `-v`/`--version` and
`-h`/`--help` are accepted as aliases.

## Source layout

```
.hatch/
  _rules/
    coding-style.md         # always-on project instructions
    testing.md              # may have an applyTo: "**/*_test.go" glob
  _skills/
    review-pr/              # skills are directories, not single files
      SKILL.md
      scripts/review.sh     # sibling assets copy through verbatim
  _commands/
    commit.md               # user-invoked slash prompts
  _agents/
    security-auditor.md     # delegated sub-agents
```

Source files are plain markdown. Skills, commands, and agents carry a small
YAML frontmatter header (see [Frontmatter](#frontmatter) below for the full
list of fields).

The leading underscore on the four primitive container directories is
intentional. It distinguishes hatch-managed primitive containers from
ordinary path components, so a monorepo can put `.hatch/backend/_rules/`
alongside `.hatch/frontend/_rules/` without ambiguity. Only those four
exact names (`_rules`, `_skills`, `_commands`, `_agents`) are recognised
as primitive containers — any other directory under `.hatch/` is treated
as a scope path component (see "Nesting for monorepos" below).

## Nesting for monorepos

Drop a path component between `.hatch/` and a primitive container, and
its outputs are emitted under the same path. So `.hatch/backend/_rules/style.md`
generates `backend/CLAUDE.md`, `backend/AGENTS.md`,
`backend/.opencode/...`, etc., and Claude Code / Codex / OpenCode pick
those up natively when working inside the `backend/` subtree:

```
.hatch/
  _rules/                          # root scope: applies repo-wide
    coding-style.md
  backend/
    _rules/
      database.md                  # only loaded when working in backend/
    _skills/
      check-migrations/
        SKILL.md
  services/api/
    _commands/
      deploy.md
```

Generated layout for the example above:

```
CLAUDE.md                          # root rules
AGENTS.md                          # root rules
backend/CLAUDE.md                  # backend rules
backend/AGENTS.md                  # backend rules
backend/.claude/skills/check-migrations/SKILL.md
services/api/CLAUDE.md             # services/api commands and rules (none here)
services/api/.claude/commands/deploy.md
…
```

GitHub Copilot and Cursor only read `.github/` and `.cursor/` from the
repo root respectively, so hatch routes their scoped output through
their native scoping mechanism (`applyTo` / `globs` frontmatter) and
slug-rewrites filenames to avoid collisions. A scoped Copilot rule
named `style` under `backend/` becomes `.github/instructions/backend-style.instructions.md`
with `applyTo: backend/**`; the same rule for Cursor becomes
`.cursor/rules/backend-style.mdc` with `globs: [backend/**]`.

Path components must avoid the four exact primitive container names
(`_rules`, `_skills`, `_commands`, `_agents`). Other names that happen
to start with `_` work, but the README-soft convention is to avoid them
in case future hatch versions add new primitive container names.

## Target mapping

| Primitive        | Claude Code                             | Codex                                   | Copilot                                      | Cursor                                       | OpenCode                              |
| ---------------- | --------------------------------------- | --------------------------------------- | -------------------------------------------- | -------------------------------------------- | ------------------------------------- |
| `rule` (plain)   | block in `CLAUDE.md`                    | block in `AGENTS.md`                    | block in `.github/copilot-instructions.md`   | `.cursor/rules/<n>.mdc` (alwaysApply: true)  | block in `AGENTS.md`                  |
| `rule` (applyTo) | block in `CLAUDE.md` with heading       | block in `AGENTS.md` with heading       | `.github/instructions/<n>.instructions.md`   | `.cursor/rules/<n>.mdc` with `globs:`        | block in `AGENTS.md` with heading     |
| `skill`          | `.claude/skills/<n>/SKILL.md`           | `.agents/skills/<n>/SKILL.md`           | inlined into the copilot-instructions block  | `.cursor/rules/skill-<n>.mdc`                | `.opencode/skills/<n>/SKILL.md`       |
| `command`        | `.claude/commands/<n>.md`               | inlined into `AGENTS.md`                | `.github/prompts/<n>.prompt.md`              | `.cursor/rules/command-<n>.mdc`              | `.opencode/commands/<n>.md`           |
| `agent`          | `.claude/agents/<n>.md`                 | inlined into `AGENTS.md`                | `.github/agents/<n>.agent.md`                | `.cursor/rules/agent-<n>.mdc`                | `.opencode/agents/<n>.md`             |

**Codex commands and sub-agents.** Codex has no first-class markdown primitive
for slash commands or sub-agents (sub-agents live in TOML config). Rather than
silently drop them, hatch inlines each one into `AGENTS.md` as a `## Commands`
or `## Sub-agents` section with per-entry instructions — so if a user asks
Codex to run a command or delegate to a sub-agent, the guidance is right there
in `AGENTS.md`.

**Copilot skills.** Copilot has no documented model-discoverable skill
primitive, so hatch inlines every `skill` body as a section inside the
hatch-managed block in `.github/copilot-instructions.md`.

**Cursor skills, commands, and sub-agents.** Cursor's only stable primitive
is the rule (`.cursor/rules/*.mdc`). Skills, commands, and sub-agents are
emitted as additional `.mdc` rule files with `alwaysApply: true` and a
kind prefix in the filename (`skill-`, `command-`, `agent-`) so the model
sees the body and the user can tell at a glance which rule files came
from which hatch primitive.

**Nested-scope routing for Copilot and Cursor.** Both only read their root
config dirs (`.github/`, `.cursor/`), so hatch routes nested-scope output
through their native scoping mechanism (`applyTo` for Copilot, `globs`
for Cursor) and slug-rewrites filenames with a `<scope>-` prefix.

## File-owned vs block-injected files

Hatch writes two kinds of generated file:

- **File-owned** — hatch writes the whole file from scratch and owns it.
  Applies to everything under `.claude/`, `.agents/`, `.github/`,
  `.cursor/`, `.opencode/`.
- **Block-injected** — hatch writes a delimited block inside a file that may
  contain user-authored content around it. Applies to `CLAUDE.md`, `AGENTS.md`,
  and `.github/copilot-instructions.md`. Content outside the markers is
  preserved across `hatch gen` and `hatch clean`.

The marker format is:

```markdown
<!-- hatch:begin v1 -->
...hatch-generated content...
<!-- hatch:end v1 -->
```

These are HTML comments so they are invisible in rendered markdown, and any
tool that recognises the marker can find and replace the block.

## Frontmatter

Skills, commands, and agents start with a YAML frontmatter block delimited by
`---`. Only `description` is required — everything else is optional.

```markdown
---
description: Review a PR for correctness, style, and tests.
name: review-pr
targets: [claude, opencode]
applyTo: "**/*.go"
claude:
  allowed-tools: [Read, Grep]
copilot:
  model: gpt-4.1
---

Body markdown. Two template vars substitute when generated files are written:
- {{agent}}  → the agent display name, e.g. "Claude Code"
- {{target}} → the target short name, e.g. "claude"
```

| Field         | Applies to              | Required    | Meaning                                                                                           |
| ------------- | ----------------------- | ----------- | ------------------------------------------------------------------------------------------------- |
| `description` | skill, command, agent   | yes         | One-sentence summary. Every downstream agent either requires or recommends this field.            |
| `name`        | skill, command, agent   | no          | Overrides the filename/dirname slug. Absent → derived from the source file path.                  |
| `targets`     | all                     | no          | List of target names this primitive should reach. Empty/absent means every target. A single string is accepted. |
| `applyTo`     | rule                    | no          | Glob pattern limiting the rule to matching paths. Copilot gets a native path-scoped instruction file; other targets wrap it in a scoped heading. |
| `claude:`     | all                     | no          | Per-target passthrough block — keys inside are merged into the generated file's own frontmatter, for that target only. Also `codex:`, `copilot:`, `opencode:`. |

Rules are plain markdown and don't need frontmatter at all. If they have one,
only `targets` and `applyTo` are meaningful.

## Development

With [mise](https://mise.jdx.dev/) installed:

```
mise run format   # go fmt ./...
mise run vet      # go vet ./...
mise run test     # go test ./...
mise run build    # go build -o bin/hatch ./cmd/hatch
mise run check    # format + vet + test
mise run install  # go install ./cmd/hatch
```

Or plain Go:

```
go test ./...
go build -o bin/hatch ./cmd/hatch
```

## Layout

```
cmd/hatch/main.go                 main binary
pkg/
  cli/                            public CLI: cli.Run(ctx, ver, targets, args, in, out, err)
  source/                         load .hatch/ into a Source
  render/                         deterministic YAML frontmatter + body templating
  block/                          hatch-marker block injection and stripping
  target/                         Target interface, Set, shared helpers
    claude/                       Claude Code generator
    codex/                        OpenAI Codex CLI generator
    copilot/                      GitHub Copilot generator
    cursor/                       Cursor generator
    opencode/                     OpenCode (sst/opencode) generator
```

`pkg/cli` is public so external tools can embed hatch's CLI with their own
`target.Set`. Every `.go` file has a matching `_test.go`; target registration
is explicit from `cmd/hatch/main.go` (no `init()` side effects).

## Design references

- Claude Code: [memory](https://code.claude.com/docs/en/memory), [skills](https://code.claude.com/docs/en/skills), [sub-agents](https://code.claude.com/docs/en/sub-agents)
- Codex: [AGENTS.md](https://developers.openai.com/codex/guides/agents-md), [skills](https://developers.openai.com/codex/skills), [config reference](https://developers.openai.com/codex/config-reference)
- Copilot: [custom instructions](https://docs.github.com/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot), [custom agents](https://docs.github.com/en/copilot/reference/custom-agents-configuration)
- Cursor: [rules](https://docs.cursor.com/context/rules)
- OpenCode: [rules](https://opencode.ai/docs/rules/), [skills](https://opencode.ai/docs/skills/), [agents](https://opencode.ai/docs/agents/)
- Shared skill standard: [agentskills.io](https://agentskills.io)

## License

MIT — see `LICENSE`.
