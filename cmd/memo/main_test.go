package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveMemoPath_UsesMEMOPath(t *testing.T) {
	want := filepath.Join(t.TempDir(), "work-memos.json")
	t.Setenv("MEMO_PATH", want)

	if got := resolveMemoPath(); got != want {
		t.Errorf("resolveMemoPath() = %q, want %q", got, want)
	}
}

func TestResolveMemoPath_Default(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv("MEMO_PATH", "")
		if err := os.Unsetenv("MEMO_PATH"); err != nil {
			t.Fatal(err)
		}
		if got := resolveMemoPath(); got != "~/.memo/memos.json" {
			t.Errorf("resolveMemoPath() = %q, want %q", got, "~/.memo/memos.json")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv("MEMO_PATH", "")
		if got := resolveMemoPath(); got != "~/.memo/memos.json" {
			t.Errorf("resolveMemoPath() = %q, want %q", got, "~/.memo/memos.json")
		}
	})

	for _, tc := range []struct {
		name string
		path string
	}{
		{name: "spaces", path: "   "},
		{name: "tabs and newline", path: "\t \n"},
		{name: "Unicode whitespace", path: "\u3000"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MEMO_PATH", tc.path)
			if got := resolveMemoPath(); got != "~/.memo/memos.json" {
				t.Errorf("resolveMemoPath() = %q, want %q", got, "~/.memo/memos.json")
			}
		})
	}
}

func TestResolveMemoPath_PreservesNonemptyPathsForStorage(t *testing.T) {
	for _, path := range []string{"~/memos/test.json", "~", "~user/memos.json", "data/memos.json", " data/memos.json "} {
		t.Run(path, func(t *testing.T) {
			t.Setenv("MEMO_PATH", path)
			if got := resolveMemoPath(); got != path {
				t.Errorf("resolveMemoPath() = %q, want %q", got, path)
			}
		})
	}
}
