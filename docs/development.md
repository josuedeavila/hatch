# Development

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

`pkg/cli` is public so external tools can embed hatch's CLI with their
own `target.Set`. Every `.go` file has a matching `_test.go`; target
registration is explicit from `cmd/hatch/main.go` (no `init()` side
effects).
