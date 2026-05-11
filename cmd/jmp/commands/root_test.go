package commands

import (
	"strings"
	"testing"
)

func TestRootHelpShowsHappyPath(t *testing.T) {
	var output string
	err := captureStdout(func() error {
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"--help"})
		return cmd.Execute()
	}, &output)
	if err != nil {
		t.Fatalf("help failed: %v", err)
	}
	for _, want := range []string{
		"Mental model:",
		"Happy path:",
		"Happy Path Commands",
		"jmp task start",
		"jmp workspace create",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in help output:\n%s", want, output)
		}
	}
}
