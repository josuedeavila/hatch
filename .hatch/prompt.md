1) GitHub Copilot (VS Code / Visual Studio / GitHub Copilot coding agent)

Rules (instructions)

Repository-wide: .github/copilot-instructions.md

Path-specific (VS Code + Copilot agent features): .github/instructions/*.instructions.md (with an applyTo field)

Coding agent also supports: AGENTS.md (and nesteds)

Commands

No standard repo “slash-commands” file. Treat “commands” as prompt templates you store in your own system, then paste/run as needed (or convert some into .instructions.md files when they’re really rules).

Ignore

No dedicated Copilot ignore file. Best practice is: avoid putting secrets in-repo; keep "don't touch" guidance inside instructions.

2) OpenAI Codex (CLI / agent workflows)

Rules (instructions)

AGENTS.md (and related / nested discovery)

Commands

Codex CLI behavior is configured primarily via CLI flags and ~/.codex/config.toml (user-level), not a repo command file.

Treat your “commands” as prompt templates, or as sections within AGENTS.md (“When asked to X, do Y”).

Ignore

No single canonical "codexignore" file surfaced in the official docs above; use instruction text ("do not open/generated/…").

3) Anthropic Claude Code

Rules (instructions)

CLAUDE.md (pulled automatically; commonly used as project memory/instructions)

Commands

Claude Code has interactive workflows, but no standard repo “commands” file. Convert “commands” into:

sections in CLAUDE.md (“Runbook: Release”, “How we do migrations”), or

standalone markdown runbooks you tell the agent to follow.

Ignore

Claude Code supports excluding sensitive files (mechanism documented), but not via a single universally-standard ignore filename.

Practically: explicit "never read these paths" rules in CLAUDE.md.

4) Cursor

Rules (instructions)

.cursor/rules/ (rules live as files in that folder)

Commands

Cursor supports custom slash commands, but they’re created/managed via Cursor UI (“Custom commands”), not a universally-documented repo command file.

Best conversion: generate a commands.md (your own) and optionally mirror the top ones into Cursor UI.

Ignore

Cursor supports controlling context, but a single canonical repo ignore file isn't clearly documented in the sources above; use rules ("don't read/build artifacts").

5) Windsurf (Codeium)

Rules (instructions)

.windsurfrules at repo root

Commands

“Workflows” are markdown files and are invoked via slash commands.

Convert your “commands” into one markdown file per command in the Windsurf workflows format (name maps to /command-name).

Ignore

Typically expressed as rules (“Don’t modify any files in …”) inside .windsurfrules.

6) Continue (Continue.dev)

Rules (instructions)

.continue/rules/*.md (loaded in lexicographic order)

Commands

Continue supports reusable “prompts” (selectable via /… inside the extension); these are configured via Continue’s config/prompt system rather than a single universal “commands” file.

Practical conversion: map each command to a Continue prompt entry (and/or keep them as markdown prompt files your generator manages).

Ignore

Prefer "do not touch/read" rules in .continue/rules/.

7) Cline (VS Code extension)

Rules (instructions)

.clinerules (single file) and/or .clinerules/ folder with markdown files

Commands

No single official repo “commands” file in the sources above; treat commands as prompt templates (often kept in-repo as docs) and referenced from .clinerules (“When I say ‘release’, follow docs/release.md”).

Ignore

Implement as:


explicit “never modify X” in .clinerules.

8) Roo Code

Rules (instructions)

.roo/rules-{slug}/ folders (mode-scoped rules; modern approach)

Commands

Roo’s “modes” and prompts are more configuration-driven; commands typically become either:

a mode with rules + prompt scaffolding, or

documented runbooks referenced by rules.

Ignore

Apply via rules; keep "excluded paths" in the relevant .roo/rules-{slug}/… files.

9) JetBrains AI Assistant

Rules (instructions)

“Project rules files” are created/managed from the IDE settings UI (rules are project-attached, but the filename/location is created by JetBrains tooling).

Commands

Prompts are managed via “Prompt Library” (UI).

Conversion: generate a jetbrains-prompts.md (your own) and/or provide importable text snippets, since this is not a simple “drop a file in the repo” workflow.

Ignore

JetBrains provides project usage restrictions, again UI-driven.

10) Zed

Rules (instructions)

Zed rules are stored locally via the Rules Library (not a repo file by default).

Commands

Typically handled as local prompt/rule artifacts, not repo-managed. Your best bet is to keep repo docs that you reference from the assistant.

Ignore

No canonical repo ignore file for AI rules; use documented guidance.

How to structure your converter (a pragmatic, “works across tools” approach)
1) Normalize your input

Model your source data as:

rules/: durable, always-on guidance (style, architecture, do/don’t)

commands/: named workflows/runbooks (release, migrate, add endpoint, etc.)

ignore/: paths/globs + rationale (“noise”, “secrets”, “generated”)

2) Emit per-agent outputs

Rules → files

Cursor: ./.cursor/rules/NN-topic.md

Windsurf: ./.windsurfrules (flatten and number)

Continue: ./.continue/rules/NN-topic.md

Cline: ./.clinerules or ./.clinerules/NN-topic.md

Roo: ./.roo/rules-code/… (and other mode slugs as needed)

Copilot: ./.github/copilot-instructions.md (+ optional .github/instructions/*.instructions.md)

Codex: ./AGENTS.md (and optionally nested)

Claude Code: ./CLAUDE.md

JetBrains/Zed: cannot be fully file-driven; generate “import/paste” artifacts (markdown) for humans to add via UI.

Commands → files

Windsurf: convert each command to a Workflow markdown file (becomes /name in Cascade).

Continue: convert to prompt entries (config-based) or keep as docs/runbooks/*.md and reference from rules.

Everyone else: safest cross-tool strategy is runbook markdown files in-repo:

docs/agent/runbooks/<command>.md

then add a short index section in the tool’s rules file: “When asked to <command>, follow docs/agent/runbooks/<command>.md exactly.”

Ignore → files

Universal: Keep an AI-specific ignore section inside each tool’s rules file:

“Never read/modify: pathA, pathB…”

“Treat as generated: …”

Where the tool doesn’t have a formal ignore mechanism (many don’t), rules-text is the mechanism.
