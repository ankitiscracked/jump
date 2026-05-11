package store

import "testing"

func TestWriteAndListEvents(t *testing.T) {
	s := OpenAt(t.TempDir())

	if err := s.WriteEvent(Event{
		ID:            100,
		Type:          "snapshot_created",
		WorkspaceName: "main",
		SnapshotID:    "snap-1",
		FilesChanged:  []string{"a.txt"},
	}); err != nil {
		t.Fatalf("WriteEvent: %v", err)
	}
	if err := s.WriteEvent(Event{
		ID:            101,
		Type:          "snapshot_created",
		WorkspaceName: "feature",
		SnapshotID:    "snap-2",
		FilesChanged:  []string{"b.txt"},
	}); err != nil {
		t.Fatalf("WriteEvent 2: %v", err)
	}

	events, err := s.ListEvents(100)
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].WorkspaceName != "feature" || events[0].SnapshotID != "snap-2" {
		t.Fatalf("unexpected event: %+v", events[0])
	}
}

func TestListEventsMissingDir(t *testing.T) {
	s := OpenAt(t.TempDir())
	events, err := s.ListEvents(0)
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
}
