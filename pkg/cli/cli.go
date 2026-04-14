// Package cli holds the subcommand dispatch and per-command handlers for
// the hatch binary. It's public so external tools can embed the same CLI —
// pass your own target.Set into Run to extend the set of supported agents.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/matryer/hatch/pkg/target"
)

// Run is the testable CLI entry point. It dispatches on args[1] and delegates
// to a per-subcommand handler. The caller supplies the target set explicitly;
// there is no global registry.
//
// `stdin` is used only by interactive subcommands (currently `hatch new` when
// no title argument is supplied). Callers that never invoke those may pass a
// nil or empty reader.
func Run(ctx context.Context, version string, targets *target.Set, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) < 2 {
		printUsage(stdout, version, targets)
		return nil
	}
	cmd, rest := args[1], args[2:]
	switch cmd {
	case "gen":
		return cmdGen(ctx, targets, rest, stdout, stderr)
	case "list":
		return cmdList(ctx, targets, rest, stdout, stderr)
	case "clean":
		return cmdClean(ctx, targets, rest, stdout, stderr)
	case "init":
		return cmdInit(ctx, rest, stdout, stderr)
	case "new":
		return cmdNew(ctx, rest, stdin, stdout, stderr)
	case "meta":
		return cmdMeta(ctx, targets, rest, stdout, stderr)
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, version)
		return nil
	case "help", "--help", "-h":
		printUsage(stdout, version, targets)
		return nil
	default:
		return fmt.Errorf("unknown command %q (try `hatch help`)", cmd)
	}
}

func printUsage(w io.Writer, version string, targets *target.Set) {
	fmt.Fprintf(w, `hatch %s

Generate rules, skills, commands, and sub-agent definitions for multiple
coding agents from a single source at .hatch/.

Usage:
  hatch gen   [-targets list]          write all target outputs
  hatch list  [-targets list]          dry-run; print what would be written
  hatch clean [-targets list]          remove generated outputs
  hatch init                           scaffold .hatch/
  hatch new <kind> [title]             create a new source file
  hatch meta skill [-targets list]     write/print a SKILL.md about hatch
  hatch version                        print version and exit
  hatch help                           this message

Flags:
  -targets list  comma-separated target names (default: all)

Flags:
  -targets list  comma-separated target names (default: all)

Registered targets: %s
`, version, strings.Join(targets.Names(), ", "))
}

// commonFlags adds the flags shared by generate/list/clean. All subcommands
// operate on the current working directory.
func commonFlags(name string, stderr io.Writer) (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	targetsList := fs.String("targets", "", "comma-separated target names (default: all)")
	return fs, targetsList
}

// selectTargets resolves a comma-separated target-name list against the
// available set. An empty value means "all".
func selectTargets(available *target.Set, list string) (*target.Set, error) {
	list = strings.TrimSpace(list)
	if list == "" {
		return available, nil
	}
	var names []string
	for _, n := range strings.Split(list, ",") {
		n = strings.TrimSpace(n)
		if n != "" {
			names = append(names, n)
		}
	}
	return available.Select(names)
}
