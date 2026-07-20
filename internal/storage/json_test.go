package storage_test

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/co191194/memo-cli/internal/memo"
	"github.com/co191194/memo-cli/internal/storage"
)

type Memo = memo.Memo
type StorageOperatorImpl = storage.StorageOperatorImpl

func TestSaveAndLoadMemos(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "memos.json")

	now := time.Date(2026, 7, 11, 10, 0, 0, 0, time.UTC)
	expected := []Memo{
		{
			ID:        1,
			Title:     "I study Go lang",
			Body:      "using testing package",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mo := StorageOperatorImpl{}

	if err := mo.SaveMemos(filePath, expected); err != nil {
		t.Fatalf("SaveMemos() error = %v", err)
	}

	actual, err := mo.LoadMemos(filePath)
	if err != nil {
		t.Fatalf("LoadMemos() error = %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("LoadMemos() = %v, expected = %v", actual, expected)
	}
}

func TestLoadMemos_FileDoesNotExist(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-exist.json")

	mo := StorageOperatorImpl{}

	actual, err := mo.LoadMemos(filePath)
	if err != nil {
		t.Fatalf("LoadMemos() error = %v", err)
	}

	if len(actual) != 0 {
		t.Errorf("len(LoadMemos()) = %d, expected = 0", len(actual))
	}
}
