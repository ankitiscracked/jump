package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ankitiscracked/jmp/internal/config"
	"github.com/ankitiscracked/jmp/internal/store"
	"github.com/ankitiscracked/jmp/internal/ui"
)

func init() {
	register(func(root *cobra.Command) {
		root.AddCommand(newEventsCmd())
		root.AddCommand(newWatchCmd())
	})
}

func newEventsCmd() *cobra.Command {
	var since int64
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "events",
		Short: "Show project event log",
		Long: `Show project-local coordination events.

Events are append-only records written under .jmp/events when jmp changes
project state. The first event type is snapshot_created.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEvents(since, jsonOutput)
		},
	}

	cmd.Flags().Int64Var(&since, "since", 0, "Only show events with ID greater than this value")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON")

	return cmd
}

func newWatchCmd() *cobra.Command {
	var since int64
	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Stream project events",
		Long:  "Stream new project-local coordination events as they are appended.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(since, interval)
		},
	}

	cmd.Flags().Int64Var(&since, "since", 0, "Start after this event ID")
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "Polling interval")

	return cmd
}

func runEvents(since int64, jsonOutput bool) error {
	s, err := openEventStore()
	if err != nil {
		return err
	}
	events, err := s.ListEvents(since)
	if err != nil {
		return err
	}
	if jsonOutput {
		data, err := json.MarshalIndent(events, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	printEvents(events)
	return nil
}

func runWatch(since int64, interval time.Duration) error {
	s, err := openEventStore()
	if err != nil {
		return err
	}
	last := since
	for {
		events, err := s.ListEvents(last)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			printEvents(events)
			last = events[len(events)-1].ID
		}
		time.Sleep(interval)
	}
}

func openEventStore() (*store.Store, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if projectRoot, _, err := config.FindProjectRootFrom(cwd); err == nil {
		return store.OpenAt(projectRoot), nil
	}
	if wsRoot, err := config.FindWorkspaceRoot(); err == nil {
		return store.OpenFromWorkspace(wsRoot), nil
	}
	return nil, fmt.Errorf("not in a project or workspace")
}

func printEvents(events []store.Event) {
	if len(events) == 0 {
		fmt.Println("No events.")
		return
	}
	for _, e := range events {
		fmt.Printf("%d  %s  %s", e.ID, e.Time, ui.Bold(e.Type))
		if e.WorkspaceName != "" {
			fmt.Printf("  %s", e.WorkspaceName)
		}
		if e.SnapshotID != "" {
			fmt.Printf("  %s", shortEventID(e.SnapshotID))
		}
		if len(e.FilesChanged) > 0 {
			fmt.Printf("  %d files", len(e.FilesChanged))
		}
		if e.Message != "" {
			fmt.Printf("  %q", e.Message)
		}
		fmt.Println()
	}
}

func shortEventID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:12]
}
