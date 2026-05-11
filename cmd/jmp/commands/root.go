package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "0.0.1"
	BuildTime = "dev"
	GitCommit = "unknown"
)

const (
	groupHappyPath    = "happy-path"
	groupInspect      = "inspect"
	groupIntegrations = "integrations"
	groupAdvanced     = "advanced"
)

var rootCmd = newRootCmd()

type registrar func(*cobra.Command)

var registrars []registrar

func register(r registrar) {
	registrars = append(registrars, r)
	if rootCmd != nil {
		r(rootCmd)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jmp",
		Short: "jmp - local snapshots and workspaces for coding agents",
		Long: `jmp is a local-first CLI for coding agents and humans working in
parallel.

Mental model:
  project    contains one or more workspaces
  workspace  is an isolated checkout for one agent or human
  snapshot   is an immutable checkpoint of a workspace
  task       groups snapshots into one unit of work
  events     let nearby agents notice workspace activity

Happy path:
  jmp project create myapp
  cd myapp/main
  jmp task start "fix auth"
  jmp snapshot -m "checkpoint"
  jmp workspace create agent-ui
  jmp drift main
  jmp merge main
  jmp task finish --summary "fixed auth"
  jmp git export --init`,
	}
	cmd.AddGroup(
		&cobra.Group{ID: groupHappyPath, Title: "Happy Path Commands"},
		&cobra.Group{ID: groupInspect, Title: "Inspect Commands"},
		&cobra.Group{ID: groupIntegrations, Title: "Integration Commands"},
		&cobra.Group{ID: groupAdvanced, Title: "Advanced Commands"},
	)
	return cmd
}

func NewRootCmd() *cobra.Command {
	cmd := newRootCmd()
	for _, r := range registrars {
		r(cmd)
	}
	configureCommandGroups(cmd)
	return cmd
}

func Execute() error {
	if len(os.Args) > 1 {
		rootCmd.SetArgs(rewriteArgs(os.Args[1:]))
	}
	configureCommandGroups(rootCmd)
	return rootCmd.Execute()
}

func configureCommandGroups(cmd *cobra.Command) {
	groups := map[string]string{
		"project":   groupHappyPath,
		"workspace": groupHappyPath,
		"task":      groupHappyPath,
		"snapshot":  groupHappyPath,
		"status":    groupHappyPath,
		"drift":     groupHappyPath,
		"diff":      groupHappyPath,
		"merge":     groupHappyPath,
		"restore":   groupHappyPath,
		"log":       groupHappyPath,
		"events":    groupHappyPath,
		"watch":     groupHappyPath,
		"git":       groupHappyPath,

		"info":    groupInspect,
		"dag":     groupInspect,
		"history": groupInspect,

		"github":  groupIntegrations,
		"backend": groupIntegrations,
		"pull":    groupIntegrations,
		"sync":    groupIntegrations,
		"agents":  groupIntegrations,
		"config":  groupIntegrations,

		"edit":      groupAdvanced,
		"drop":      groupAdvanced,
		"squash":    groupAdvanced,
		"rebase":    groupAdvanced,
		"gc":        groupAdvanced,
		"ui":        groupAdvanced,
		"conflicts": groupAdvanced,
		"version":   groupAdvanced,
	}
	for _, child := range cmd.Commands() {
		if groupID, ok := groups[child.Name()]; ok {
			child.GroupID = groupID
		}
	}
}

func rewriteArgs(args []string) []string {
	rewritten := make([]string, 0, len(args))
	for _, arg := range args {
		switch {
		case arg == "-am":
			rewritten = append(rewritten, "--agent-message")
		case strings.HasPrefix(arg, "-am="):
			rewritten = append(rewritten, "--agent-message"+arg[len("-am"):])
		default:
			rewritten = append(rewritten, arg)
		}
	}
	return rewritten
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("jmp version %s\n", Version)
			fmt.Printf("  Build time: %s\n", BuildTime)
			fmt.Printf("  Git commit: %s\n", GitCommit)
		},
	}
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(newVersionCmd()) })
}
