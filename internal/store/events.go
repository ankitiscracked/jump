package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const eventsDirName = "events"

// Event is an append-only project-local coordination record.
type Event struct {
	ID                  int64    `json:"id"`
	Type                string   `json:"type"`
	Time                string   `json:"time"`
	TaskID              string   `json:"task_id,omitempty"`
	WorkspaceID         string   `json:"workspace_id,omitempty"`
	WorkspaceName       string   `json:"workspace_name,omitempty"`
	SourceWorkspaceID   string   `json:"source_workspace_id,omitempty"`
	SourceWorkspaceName string   `json:"source_workspace_name,omitempty"`
	SnapshotID          string   `json:"snapshot_id,omitempty"`
	ParentSnapshotIDs   []string `json:"parent_snapshot_ids,omitempty"`
	FilesChanged        []string `json:"files_changed,omitempty"`
	Message             string   `json:"message,omitempty"`
	Agent               string   `json:"agent,omitempty"`
}

func (s *Store) eventsDir() string {
	return filepath.Join(s.root, configDirName, eventsDirName)
}

// WriteEvent appends an event to the project-local event log.
func (s *Store) WriteEvent(e Event) error {
	if e.Type == "" {
		return fmt.Errorf("event type is required")
	}
	if e.Time == "" {
		e.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}
	if e.ID == 0 {
		e.ID = time.Now().UTC().UnixNano()
	}
	if err := os.MkdirAll(s.eventsDir(), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}

	base := sanitizeEventType(e.Type)
	for i := int64(0); ; i++ {
		id := e.ID + i
		path := filepath.Join(s.eventsDir(), fmt.Sprintf("%020d-%s.json", id, base))
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if i != 0 {
			e.ID = id
			data, err = json.MarshalIndent(e, "", "  ")
			if err != nil {
				return err
			}
		}
		return AtomicWriteFile(path, data, 0644)
	}
}

// ListEvents returns project events sorted by ID, filtering IDs <= sinceID.
func (s *Store) ListEvents(sinceID int64) ([]Event, error) {
	entries, err := os.ReadDir(s.eventsDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	events := make([]Event, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.eventsDir(), entry.Name()))
		if err != nil {
			continue
		}
		var e Event
		if err := json.Unmarshal(data, &e); err != nil {
			continue
		}
		if e.ID == 0 {
			e.ID = eventIDFromName(entry.Name())
		}
		if e.ID <= sinceID {
			continue
		}
		events = append(events, e)
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].ID == events[j].ID {
			return events[i].Type < events[j].Type
		}
		return events[i].ID < events[j].ID
	})
	return events, nil
}

func eventIDFromName(name string) int64 {
	part := strings.SplitN(name, "-", 2)[0]
	id, _ := strconv.ParseInt(part, 10, 64)
	return id
}

func sanitizeEventType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	if t == "" {
		return "event"
	}
	var b strings.Builder
	for _, r := range t {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	return b.String()
}
