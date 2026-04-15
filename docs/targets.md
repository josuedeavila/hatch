# Targets and output mapping

## Mapping

| Primitive        | Claude Code                             | Codex                                   | Copilot                                      | Cursor                                       | OpenCode                              |
| ---------------- | --------------------------------------- | --------------------------------------- | -------------------------------------------- | -------------------------------------------- | ------------------------------------- |
| `rule` (plain)   | block in `CLAUDE.md`                    | block in `AGENTS.md`                    | block in `.github/copilot-instructions.md`   | `.cursor/rules/<n>.mdc` (alwaysApply: true)  | block in `AGENTS.md`                  |
| `rule` (applyTo) | block in `CLAUDE.md` with heading       | block in `AGENTS.md` with heading       | `.github/instructions/<n>.instructions.md`   | `.cursor/rules/<n>.mdc` with `globs:`        | block in `AGENTS.md` with heading     |
| `skill`          | `.claude/skills/<n>/SKILL.md`           | `.agents/skills/<n>/SKILL.md`           | inlined into the copilot-instructions block  | `.cursor/rules/skill-<n>.mdc`                | `.opencode/skills/<n>/SKILL.md`       |
| `command`        | `.claude/commands/<n>.md`               | inlined into `AGENTS.md`                | `.github/prompts/<n>.prompt.md`              | `.cursor/rules/command-<n>.mdc`              | `.opencode/commands/<n>.md`           |
| `agent`          | `.claude/agents/<n>.md`                 | inlined into `AGENTS.md`                | `.github/agents/<n>.agent.md`                | `.cursor/rules/agent-<n>.mdc`                | `.opencode/agents/<n>.md`             |

## Per-target notes

**Codex commands and sub-agents.** Codex has no first-class markdown
primitive for slash commands or sub-agents. Rather than drop them,
hatch inlines each one into `AGENTS.md` as a `## Commands` or
`## Sub-agents` section with per-entry instructions.

**Cursor non-rule primitives.** Cursor's only stable primitive is the
rule. Skills, commands, and sub-agents become additional `.mdc` rule
files with a kind prefix in the filename (`skill-`, `command-`,
`agent-`) so the user can tell at a glance which rule files came from
which hatch primitive.

**Nested-scope routing for Copilot and Cursor.** Both only read their
root config dirs (`.github/`, `.cursor/`), so hatch routes nested-scope
output through their native scoping mechanism (`applyTo` for Copilot,
`globs` for Cursor) and slug-rewrites filenames with a `<scope>-`
prefix. A scoped Copilot rule named `style` under `backend/` becomes
`.github/instructions/backend-style.instructions.md` with
`applyTo: backend/**`.

## File-owned vs block-injected files

Hatch writes two kinds of generated file:

- **File-owned** — hatch writes the whole file from scratch. Applies to
  everything under `.claude/`, `.agents/`, `.github/`, `.cursor/`, and
  `.opencode/`.
- **Block-injected** — hatch writes a delimited block inside a file that
  may contain user-authored content around it. Applies to `CLAUDE.md`,
  `AGENTS.md`, and `.github/copilot-instructions.md`. Content outside
  the markers is preserved across `hatch gen` and `hatch clean`.

The marker format:

```markdown
<!-- hatch:begin v1 -->
...hatch-generated content...
<!-- hatch:end v1 -->
```

## Design references

- Claude Code: [memory](https://code.claude.com/docs/en/memory), [skills](https://code.claude.com/docs/en/skills), [sub-agents](https://code.claude.com/docs/en/sub-agents)
- Codex: [AGENTS.md](https://developers.openai.com/codex/guides/agents-md), [skills](https://developers.openai.com/codex/skills), [config reference](https://developers.openai.com/codex/config-reference)
- Copilot: [custom instructions](https://docs.github.com/copilot/customizing-copilot/adding-custom-instructions-for-github-copilot), [custom agents](https://docs.github.com/en/copilot/reference/custom-agents-configuration)
- Cursor: [rules](https://docs.cursor.com/context/rules)
- OpenCode: [rules](https://opencode.ai/docs/rules/), [skills](https://opencode.ai/docs/skills/), [agents](https://opencode.ai/docs/agents/)
- Shared skill standard: [agentskills.io](https://agentskills.io)
