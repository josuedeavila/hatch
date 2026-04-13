# hatch

Write rules, skills, commands, and sub-agent definitions **once**, generate
the native files each coding agent expects.

Hatch reads a single source tree under `.hatch/src/` and produces the
specific files Claude Code, OpenAI Codex CLI, GitHub Copilot, and OpenCode
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
hatch init    # scaffold .hatch/src/ with one example of each primitive
hatch build   # write all target files
hatch list    # dry-run; show what would be written
hatch clean   # remove everything hatch generated
```

Every subcommand accepts `-C dir` to operate on a different directory and
`-targets list` (comma-separated) to narrow the set of agents to emit for.

## Source layout

```
.hatch/
  src/
    rules/
      coding-style.md         # always-on project instructions
      testing.md              # may have an applyTo: "**/*_test.go" glob
    skills/
      review-pr/              # skills are directories, not single files
        SKILL.md
        scripts/review.sh     # sibling assets copy through verbatim
    commands/
      commit.md               # user-invoked slash prompts
    agents/
      security-auditor.md     # delegated sub-agents
```

Source files are markdown. Frontmatter is YAML:

```markdown
---
description: Review a PR for correctness, style, and tests.
name: review-pr               # optional; derived from filename/dirname if absent
targets: [claude, opencode]   # optional; default = all
applyTo: "**/*.go"            # rules only; scopes the rule to a glob
claude:                       # optional per-target frontmatter passthrough
  allowed-tools: [Read, Grep]
copilot:
  model: gpt-4.1
---

Body markdown. Two template vars substitute at emission time:
- {{agent}}  → the agent display name, e.g. "Claude Code"
- {{target}} → the target short name, e.g. "claude"
```

`description` is the only required field on skills, commands, and agents.

## Target mapping

| Primitive        | Claude Code                             | Codex                                   | Copilot                                      | OpenCode                              |
| ---------------- | --------------------------------------- | --------------------------------------- | -------------------------------------------- | ------------------------------------- |
| `rule` (plain)   | block in `CLAUDE.md`                    | block in `AGENTS.md`                    | block in `.github/copilot-instructions.md`   | block in `AGENTS.md`                  |
| `rule` (applyTo) | block in `CLAUDE.md` with heading       | block in `AGENTS.md` with heading       | `.github/instructions/<n>.instructions.md`   | block in `AGENTS.md` with heading     |
| `skill`          | `.claude/skills/<n>/SKILL.md`           | `.agents/skills/<n>/SKILL.md`           | inlined into the copilot-instructions block  | `.opencode/skills/<n>/SKILL.md`       |
| `command`        | `.claude/commands/<n>.md`               | *skipped*                               | `.github/prompts/<n>.prompt.md`              | `.opencode/commands/<n>.md`           |
| `agent`          | `.claude/agents/<n>.md`                 | *skipped*                               | `.github/agents/<n>.agent.md`                | `.opencode/agents/<n>.md`             |

**Codex sub-agents and slash commands** have no first-class markdown primitive
in the Codex docs — Codex sub-agents live in TOML config, and there's no
documented slash-command file. A hatch `agent` or `command` is skipped for the
Codex target; to reach Codex, express the same content as a `skill`.

**Copilot skills:** Copilot has no documented model-discoverable skill
primitive, so hatch inlines every `skill` body as a section inside the
hatch-managed block in `.github/copilot-instructions.md`.

## Emission modes

Hatch writes two kinds of artifact:

- **File-owned** — hatch writes the whole file from scratch and owns it.
  Applies to everything under `.claude/`, `.agents/`, `.github/`, `.opencode/`.
- **Block-injected** — hatch writes a delimited block inside a file that may
  contain user-authored content around it. Applies to `CLAUDE.md`, `AGENTS.md`,
  and `.github/copilot-instructions.md`. Content outside the markers is
  preserved across `hatch build` and `hatch clean`.

The marker format is:

```markdown
<!-- hatch:begin v1 -->
...hatch-generated content...
<!-- hatch:end v1 -->
```

These are HTML comments so they are invisible in rendered markdown, and any
tool that recognises the marker can find and replace the block.

## No manifest

Hatch does not keep a state file. Everything it knows about what exists is
derived from `.hatch/src/`: `hatch clean` computes what a fresh build would
write, then deletes those files (for file-owned artifacts) or strips just
the hatch block (for block-injected files). Rename a rule, run `hatch clean`
afterwards against the renamed source, and the new paths are cleaned — the
old orphans will need a manual `rm` because the new source no longer points
at them.

## Development

With [mise](https://mise.jdx.dev/) installed:

```
mise run format   # gofmt -w .
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
cmd/hatch/                        main binary + subcommand CLI
  main.go
  internal/cli/                   subcommand dispatch + handlers
pkg/
  source/                         load .hatch/src/ into a Source
  config/                         optional .hatch/config.yaml
  render/                         deterministic YAML frontmatter + body templating
  block/                          hatch-marker block injection and stripping
  target/                         Target interface, Set, shared helpers
    claude/                       Claude Code emitter
    codex/                        OpenAI Codex CLI emitter
    copilot/                      GitHub Copilot emitter
    opencode/                     OpenCode (sst/opencode) emitter
```

Every `.go` file has a matching `_test.go`; target registration is explicit
from `cmd/hatch/main.go` (no `init()` side effects).

## Design references

- Claude Code: [memory](https://code.claude.com/docs/en/memory), [skills](https://code.claude.com/docs/en/skills), [sub-agents](https://code.claude.com/docs/en/sub-agents)
- Codex: [AGENTS.md](https://developers.openai.com/codex/guides/agents-md), [skills](https://developers.openai.com/codex/skills), [config reference](https://developers.openai.com/codex/config-reference)
- Copilot: [custom instructions](https://docs.github.com/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot), [custom agents](https://docs.github.com/en/copilot/reference/custom-agents-configuration)
- OpenCode: [rules](https://opencode.ai/docs/rules/), [skills](https://opencode.ai/docs/skills/), [agents](https://opencode.ai/docs/agents/)
- Shared skill standard: [agentskills.io](https://agentskills.io)

## License

MIT — see `LICENSE`.
