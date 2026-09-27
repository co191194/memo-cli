package storage

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveMemos_RenameFailurePreservesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "memos.json")
	before := []byte("[\n  {\"id\": 1, \"title\": \"Original\", \"body\": \"\", \"created_at\": \"2026-01-01T00:00:00Z\", \"updated_at\": \"2026-01-01T00:00:00Z\"}\n]\n")
	if err := os.WriteFile(path, before, 0600); err != nil {
		t.Fatal(err)
	}

	injected := errors.New("rename failed")
	operator := &StorageOperatorImpl{rename: func(oldPath, newPath string) error {
		if filepath.Dir(oldPath) != dir || newPath != path {
			t.Errorf("rename(%q, %q): wrong paths", oldPath, newPath)
		}
		return injected
	}}
	if err := operator.SaveMemos(path, []Memo{{ID: 2, Title: "Replacement"}}); !errors.Is(err, injected) {
		t.Fatalf("SaveMemos() error = %v, want injected error", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Errorf("existing JSON changed: before = %q, after = %q", before, after)
	}
	if matches, err := filepath.Glob(filepath.Join(dir, ".memos-*.json")); err != nil || len(matches) != 0 {
		t.Errorf("temporary files = %v, error = %v", matches, err)
	}
}
