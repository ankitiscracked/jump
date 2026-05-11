package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ankitiscracked/jmp/internal/config"
	"github.com/ankitiscracked/jmp/internal/store"
)

const currentTaskFileName = "current-task.json"

type currentTaskRef struct {
	TaskID string `json:"task_id"`
}

func (ws *Workspace) currentTaskPath() string {
	return filepath.Join(ws.root, config.ConfigDirName, currentTaskFileName)
}

// CurrentTaskID returns the active task for this workspace, if one exists.
func (ws *Workspace) CurrentTaskID() string {
	data, err := os.ReadFile(ws.currentTaskPath())
	if err != nil {
		return ""
	}
	var ref currentTaskRef
	if err := json.Unmarshal(data, &ref); err != nil {
		return ""
	}
	return ref.TaskID
}

// SetCurrentTaskID sets this workspace's active task pointer.
func (ws *Workspace) SetCurrentTaskID(taskID string) error {
	data, err := json.MarshalIndent(currentTaskRef{TaskID: taskID}, "", "  ")
	if err != nil {
		return err
	}
	return store.AtomicWriteFile(ws.currentTaskPath(), data, 0644)
}

// ClearCurrentTask clears this workspace's active task pointer.
func (ws *Workspace) ClearCurrentTask() error {
	if err := os.Remove(ws.currentTaskPath()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
