package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// TestRun_Version is a minimal smoke test that drives main's run()
// function end-to-end. Subcommand behavior is exercised in the cli
// package tests; this test just confirms dependency wiring.
func TestRun_Version(t *testing.T) {
	is := is.New(t)
	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{"hatch", "version"}, &stdout, &stderr)
	is.NoErr(err)
	is.True(strings.Contains(stdout.String(), "dev"))
}

func TestRun_Help(t *testing.T) {
	is := is.New(t)
	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{"hatch", "help"}, &stdout, &stderr)
	is.NoErr(err)
	is.True(strings.Contains(stdout.String(), "Registered targets"))
	is.True(strings.Contains(stdout.String(), "claude"))
	is.True(strings.Contains(stdout.String(), "codex"))
	is.True(strings.Contains(stdout.String(), "copilot"))
	is.True(strings.Contains(stdout.String(), "opencode"))
}
