package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ankitiscracked/jmp/internal/store"
	"github.com/ankitiscracked/jmp/internal/ui"
	"github.com/ankitiscracked/jmp/internal/workspace"
)

func init() {
	register(func(root *cobra.Command) { root.AddCommand(newTaskCmd()) })
}

func newTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Group snapshots into a unit of work",
	}
	cmd.AddCommand(newTaskStartCmd())
	cmd.AddCommand(newTaskStatusCmd())
	cmd.AddCommand(newTaskFinishCmd())
	cmd.AddCommand(newTaskListCmd())
	return cmd
}

func newTaskStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start <name>",
		Short: "Start a task in the current workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTaskStart(args[0])
		},
	}
}

func newTaskStatusCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "status [task-id]",
		Short: "Show current or named task",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return runTaskStatus(id, jsonOutput)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON")
	return cmd
}

func newTaskFinishCmd() *cobra.Command {
	var summary string
	cmd := &cobra.Command{
		Use:   "finish",
		Short: "Finish the current task",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTaskFinish(summary)
		},
	}
	cmd.Flags().StringVar(&summary, "summary", "", "Task summary")
	return cmd
}

func newTaskListCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List project tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTaskList(jsonOutput)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON")
	return cmd
}

func runTaskStart(name string) error {
	ws, err := workspace.Open()
	if err != nil {
		return fmt.Errorf("not in a workspace directory - run 'jmp workspace init' first")
	}
	defer ws.Close()

	if existing := ws.CurrentTaskID(); existing != "" {
		return fmt.Errorf("task already active: %s\nRun 'jmp task finish' first.", existing)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	task := &store.Task{
		ID:                newTaskID(name),
		Name:              name,
		Status:            "active",
		WorkspaceID:       ws.WorkspaceID(),
		WorkspaceName:     ws.WorkspaceName(),
		BaseSnapshotID:    ws.CurrentSnapshotID(),
		CurrentSnapshotID: ws.CurrentSnapshotID(),
		StartedAt:         now,
	}
	if err := ws.Store().SaveTask(task); err != nil {
		return err
	}
	if err := ws.SetCurrentTaskID(task.ID); err != nil {
		return err
	}
	_ = ws.Store().WriteEvent(store.Event{
		Type:          "task_started",
		Time:          now,
		TaskID:        task.ID,
		WorkspaceID:   ws.WorkspaceID(),
		WorkspaceName: ws.WorkspaceName(),
		SnapshotID:    ws.CurrentSnapshotID(),
		Message:       name,
	})

	fmt.Printf("Started task %s\n", ui.Bold(task.ID))
	fmt.Printf("  Name:      %s\n", task.Name)
	fmt.Printf("  Workspace: %s\n", task.WorkspaceName)
	if task.BaseSnapshotID != "" {
		fmt.Printf("  Base:      %s\n", shortEventID(task.BaseSnapshotID))
	}
	return nil
}

func runTaskStatus(id string, jsonOutput bool) error {
	if id != "" {
		s, err := openEventStore()
		if err != nil {
			return err
		}
		task, err := s.LoadTask(id)
		if err != nil {
			return err
		}
		return printTask(task, jsonOutput)
	}

	ws, err := workspace.Open()
	if err != nil {
		return fmt.Errorf("not in a workspace directory - run 'jmp workspace init' first")
	}
	defer ws.Close()

	id = ws.CurrentTaskID()
	if id == "" {
		return fmt.Errorf("no active task")
	}
	task, err := ws.Store().LoadTask(id)
	if err != nil {
		return err
	}
	return printTask(task, jsonOutput)
}

func runTaskFinish(summary string) error {
	ws, err := workspace.Open()
	if err != nil {
		return fmt.Errorf("not in a workspace directory - run 'jmp workspace init' first")
	}
	defer ws.Close()

	id := ws.CurrentTaskID()
	if id == "" {
		return fmt.Errorf("no active task")
	}
	task, err := ws.Store().LoadTask(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	task.Status = "finished"
	task.CurrentSnapshotID = ws.CurrentSnapshotID()
	task.FinishedAt = now
	task.Summary = summary
	if err := ws.Store().SaveTask(task); err != nil {
		return err
	}
	if err := ws.ClearCurrentTask(); err != nil {
		return err
	}
	_ = ws.Store().WriteEvent(store.Event{
		Type:          "task_finished",
		Time:          now,
		TaskID:        task.ID,
		WorkspaceID:   ws.WorkspaceID(),
		WorkspaceName: ws.WorkspaceName(),
		SnapshotID:    ws.CurrentSnapshotID(),
		FilesChanged:  task.FilesChanged,
		Message:       task.Name,
	})

	fmt.Printf("Finished task %s\n", ui.Bold(task.ID))
	return nil
}

func runTaskList(jsonOutput bool) error {
	s, err := openEventStore()
	if err != nil {
		return err
	}
	tasks, err := s.ListTasks()
	if err != nil {
		return err
	}
	if jsonOutput {
		data, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks.")
		return nil
	}
	for _, task := range tasks {
		fmt.Printf("%s  %s  %s  %s\n", task.ID, task.Status, task.WorkspaceName, task.Name)
	}
	return nil
}

func printTask(task *store.Task, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	fmt.Printf("Task:      %s\n", ui.Bold(task.ID))
	fmt.Printf("Name:      %s\n", task.Name)
	fmt.Printf("Status:    %s\n", task.Status)
	fmt.Printf("Workspace: %s\n", task.WorkspaceName)
	if task.BaseSnapshotID != "" {
		fmt.Printf("Base:      %s\n", shortEventID(task.BaseSnapshotID))
	}
	if task.CurrentSnapshotID != "" {
		fmt.Printf("Current:   %s\n", shortEventID(task.CurrentSnapshotID))
	}
	fmt.Printf("Snapshots: %d\n", len(task.Snapshots))
	fmt.Printf("Files:     %d\n", len(task.FilesChanged))
	if task.Summary != "" {
		fmt.Printf("Summary:   %s\n", task.Summary)
	}
	return nil
}

func newTaskID(name string) string {
	slug := strings.Trim(strings.ToLower(name), " \t\r\n")
	var b strings.Builder
	lastDash := false
	for _, r := range slug {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		out = "task"
	}
	if len(out) > 40 {
		out = out[:40]
	}
	return fmt.Sprintf("task-%s-%d", out, time.Now().UTC().UnixNano())
}
