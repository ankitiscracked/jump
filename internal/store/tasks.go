package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const tasksDirName = "tasks"

// Task groups snapshots for one agent/human unit of work.
type Task struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Status            string   `json:"status"`
	WorkspaceID       string   `json:"workspace_id"`
	WorkspaceName     string   `json:"workspace_name"`
	BaseSnapshotID    string   `json:"base_snapshot_id,omitempty"`
	CurrentSnapshotID string   `json:"current_snapshot_id,omitempty"`
	Snapshots         []string `json:"snapshots,omitempty"`
	FilesChanged      []string `json:"files_changed,omitempty"`
	Summary           string   `json:"summary,omitempty"`
	StartedAt         string   `json:"started_at"`
	FinishedAt        string   `json:"finished_at,omitempty"`
}

func (s *Store) tasksDir() string {
	return filepath.Join(s.root, configDirName, tasksDirName)
}

func (s *Store) taskPath(id string) string {
	return filepath.Join(s.tasksDir(), id+".json")
}

// SaveTask writes task metadata.
func (s *Store) SaveTask(task *Task) error {
	if task == nil {
		return fmt.Errorf("task is nil")
	}
	if task.ID == "" {
		return fmt.Errorf("task ID is required")
	}
	if task.StartedAt == "" {
		task.StartedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	if err := os.MkdirAll(s.tasksDir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWriteFile(s.taskPath(task.ID), data, 0644)
}

// LoadTask reads task metadata.
func (s *Store) LoadTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID is required")
	}
	data, err := os.ReadFile(s.taskPath(id))
	if err != nil {
		return nil, err
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// ListTasks returns all tasks sorted by start time.
func (s *Store) ListTasks() ([]Task, error) {
	entries, err := os.ReadDir(s.tasksDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	tasks := make([]Task, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		task, err := s.LoadTask(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			continue
		}
		tasks = append(tasks, *task)
	}
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].StartedAt == tasks[j].StartedAt {
			return tasks[i].ID < tasks[j].ID
		}
		return tasks[i].StartedAt < tasks[j].StartedAt
	})
	return tasks, nil
}

// RecordTaskSnapshot appends a snapshot and changed files to a task.
func (s *Store) RecordTaskSnapshot(taskID, snapshotID string, files []string) error {
	task, err := s.LoadTask(taskID)
	if err != nil {
		return err
	}
	task.CurrentSnapshotID = snapshotID
	if snapshotID != "" && !containsString(task.Snapshots, snapshotID) {
		task.Snapshots = append(task.Snapshots, snapshotID)
	}
	for _, file := range files {
		if file != "" && !containsString(task.FilesChanged, file) {
			task.FilesChanged = append(task.FilesChanged, file)
		}
	}
	return s.SaveTask(task)
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
