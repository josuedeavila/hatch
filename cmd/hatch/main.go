// Command hatch generates rules, skills, commands, and agent definitions for
// multiple coding agents (Claude Code, Codex, Copilot, OpenCode) from a single
// neutral source tree under `.hatch/`.
//
// Usage:
//
//	hatch generate [-C dir] [-targets list]    write all target outputs
//	hatch list     [-C dir] [-targets list]    dry-run (print what would be written)
//	hatch clean    [-C dir] [-targets list]    remove generated outputs
//	hatch init     [-C dir]                    scaffold .hatch/
//	hatch version
//	hatch help
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/matryer/hatch/pkg/cli"
	"github.com/matryer/hatch/pkg/target"
	"github.com/matryer/hatch/pkg/target/claude"
	"github.com/matryer/hatch/pkg/target/codex"
	"github.com/matryer/hatch/pkg/target/copilot"
	"github.com/matryer/hatch/pkg/target/opencode"
)

// version is overridden at build time via -ldflags "-X main.version=..."
var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "hatch: %s\n", err)
		os.Exit(1)
	}
}

// run is the testable entry point. It takes its dependencies explicitly so
// tests can drive the CLI with fake args and capture output.
func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	targets := target.NewSet(
		claude.New(),
		codex.New(),
		copilot.New(),
		opencode.New(),
	)
	return cli.Run(ctx, version, targets, args, stdout, stderr)
}
