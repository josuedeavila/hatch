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

Use `hatch new <kind> [title]` to scaffold a new source file, `hatch list`
to preview what `gen` would write, and `hatch clean` to remove everything
hatch generated.

## Docs

- [CLI reference](docs/cli.md) — every subcommand and flag
- [Source layout, nesting, and frontmatter](docs/sources.md) — how `.hatch/` is organised and what frontmatter fields do
- [Targets and output mapping](docs/targets.md) — per-agent file locations and nested-scope routing
- [Development](docs/development.md) — build and test the hatch binary itself

## License

MIT — see `LICENSE`.
