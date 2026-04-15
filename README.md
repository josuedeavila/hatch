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

```
$ cd my-project
$ hatch init -examples
$ hatch gen
wrote CLAUDE.md:1-3
wrote .claude/skills/review-pr/SKILL.md
wrote AGENTS.md:1-3
...
```

Edit files under `.hatch/_rules/`, `.hatch/_skills/`, `.hatch/_commands/`,
or `.hatch/_agents/`, then re-run `hatch gen`. Commit the `.hatch/` source
*and* the generated files together.

## CLI

| Command                              | What it does                                                    |
| ------------------------------------ | --------------------------------------------------------------- |
| `hatch init [-examples] [-path p]`   | scaffold `.hatch/` (optionally with example files or nested)    |
| `hatch new <kind> [-path p] [title]` | create a new rule, skill, command, or agent from a template     |
| `hatch gen [-targets list]`          | write every target's native files                               |
| `hatch list [-targets list]`         | dry-run; print what `gen` would write                           |
| `hatch clean [-targets list]`        | remove everything hatch generated                               |
| `hatch version`, `hatch help`        | print version / usage                                           |

See [docs/cli.md](docs/cli.md) for full flag reference, including
`-no-hatch-skill` and the auto-injected meta skill.

## Docs

- [CLI reference](docs/cli.md) — every subcommand and flag
- [Source layout, nesting, and frontmatter](docs/sources.md) — how `.hatch/` is organised and what frontmatter fields do
- [Targets and output mapping](docs/targets.md) — per-agent file locations and nested-scope routing
- [Development](docs/development.md) — build and test the hatch binary itself

## License

MIT — see `LICENSE`.
