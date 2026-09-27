package storage_test

import (
	"os"
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

func TestStoragePaths(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		want func(cwd, home string) string
	}{
		{
			name: "home prefix",
			path: "~/memos/test.json",
			want: func(_, home string) string { return filepath.Join(home, "memos", "test.json") },
		},
		{
			name: "bare tilde",
			path: "~",
			want: func(cwd, _ string) string { return filepath.Join(cwd, "~") },
		},
		{
			name: "named user",
			path: "~user/memos.json",
			want: func(cwd, _ string) string { return filepath.Join(cwd, "~user", "memos.json") },
		},
		{
			name: "relative",
			path: "data/memos.json",
			want: func(cwd, _ string) string { return filepath.Join(cwd, "data", "memos.json") },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cwd := t.TempDir()
			home := t.TempDir()
			t.Chdir(cwd)
			t.Setenv("HOME", home)

			memos := []Memo{{ID: 1, Title: "Work"}}
			mo := StorageOperatorImpl{}
			if err := mo.SaveMemos(tc.path, memos); err != nil {
				t.Fatalf("SaveMemos(%q) error = %v", tc.path, err)
			}
			wantPath := tc.want(cwd, home)
			if _, err := os.Stat(wantPath); err != nil {
				t.Fatalf("saved file %q: %v", wantPath, err)
			}

			got, err := mo.LoadMemos(tc.path)
			if err != nil {
				t.Fatalf("LoadMemos(%q) error = %v", tc.path, err)
			}
			if !reflect.DeepEqual(got, memos) {
				t.Errorf("LoadMemos(%q) = %v, want %v", tc.path, got, memos)
			}
		})
	}
}
