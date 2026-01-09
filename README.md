# 🥚 hatch
Write agent rules and commands once and generate for all. Support everybody's favourite LLM IDE at the same time, with one set of source files.

1. Copy the .hatch folder into your repo
2. Edit the .hatch/src/* files
3. Run ./hatch/gen.sh to generate config for all supported IDEs

**Features:**

- `/src/commands` - each file is a command the LLM can carry out
- `/src/rules` - each file becomes a rule that the LLM will follow
- `/src/.ignore` - list of files and folders the LLM should ignore

**Supported IDEs:**

- Aider
- Claude Code
- Cline
- Cody
- Continue
- Cursor
- GitHub Copilot
- VS Code
- Windsurf
- Zed
