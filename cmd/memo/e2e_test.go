package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type savedMemo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type processResult struct {
	code   int
	stdout string
	stderr string
}

func buildMemo(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "memo")
	cmd := exec.Command("go", "build", "-o", path, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return path
}

func runMemo(t *testing.T, binary, cwd, home, memoPath string, args ...string) processResult {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(), "HOME="+home, "MEMO_PATH="+memoPath, "TZ=UTC")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("memo %q: %v", args, err)
		}
		code = exitErr.ExitCode()
	}
	return processResult{code, stdout.String(), stderr.String()}
}

func expectProcess(t *testing.T, got processResult, code int, stdout, stderr string) {
	t.Helper()
	if got.code != code || got.stdout != stdout || got.stderr != stderr {
		t.Errorf("process = %+v; want code=%d stdout=%q stderr=%q", got, code, stdout, stderr)
	}
}

func expectProcessError(t *testing.T, got processResult, message string) {
	t.Helper()
	if got.code != 1 || got.stdout != "" || !strings.Contains(got.stderr, message) {
		t.Errorf("process = %+v; want code=1, empty stdout and stderr containing %q", got, message)
	}
}

func readSavedMemos(t *testing.T, path string) []savedMemo {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var result []savedMemo
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestMemoProcessWorkflow(t *testing.T) {
	binary := buildMemo(t)
	cwd := t.TempDir()
	home := t.TempDir()
	pathA := filepath.Join(cwd, "nested", "a.json")
	pathB := filepath.Join(cwd, "b.json")
	run := func(path string, args ...string) processResult { return runMemo(t, binary, cwd, home, path, args...) }

	expectProcess(t, run(pathA, "list"), 0, "No memos found.\n", "")
	if _, err := os.Stat(pathA); !os.IsNotExist(err) {
		t.Fatalf("list should not create file: %v", err)
	}
	expectProcess(t, run(pathA, "add", "Original", "--body", "old-marker"), 0, "", "")
	expectProcess(t, run(pathA, "add", "Untouched", "--body", "other"), 0, "", "")
	before := readSavedMemos(t, pathA)
	if len(before) != 2 || before[0].ID != 1 || before[0].Title != "Original" || before[0].Body != "old-marker" || before[1].ID != 2 {
		t.Fatalf("added memos = %+v", before)
	}
	expectProcess(t, run(pathA, "edit", "1", "--title", "Revised", "--body", "new-marker\nsecond line"), 0, "", "")
	after := readSavedMemos(t, pathA)
	if len(after) != 2 || after[0].ID != 1 || after[0].Title != "Revised" || after[0].Body != "new-marker\nsecond line" ||
		!after[0].CreatedAt.Equal(before[0].CreatedAt) || !after[0].UpdatedAt.After(before[0].UpdatedAt) || !reflect.DeepEqual(after[1], before[1]) {
		t.Fatalf("edited memos = %+v; before = %+v", after, before)
	}
	show := fmt.Sprintf("# Revised\n\nID: 1\nCreated: %s\nUpdated: %s\n\nnew-marker\nsecond line\n",
		before[0].CreatedAt.Format("2006-01-02 15:04"), after[0].UpdatedAt.Format("2006-01-02 15:04"))
	expectProcess(t, run(pathA, "show", "1"), 0, show, "")
	for _, keyword := range []string{"Revised", "new-marker"} {
		expectProcess(t, run(pathA, "search", keyword), 0,
			fmt.Sprintf("1 Revised\t%s\n", before[0].CreatedAt.Format("2006-01-02")), "")
	}
	expectProcessError(t, run(pathA, "search", "old-marker"), "No matching memos found.")

	fileBefore, err := os.ReadFile(pathA)
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"show", "0"}, {"show", "--", "-1"},
		{"edit", "0", "--title", "Invalid"}, {"edit", "--body", "Invalid", "--", "-1"},
		{"delete", "0"}, {"delete", "--", "-1"},
	} {
		expectProcessError(t, run(pathA, args...), "idは1以上で指定してください")
	}
	fileAfter, err := os.ReadFile(pathA)
	if err != nil || !bytes.Equal(fileBefore, fileAfter) {
		t.Fatalf("invalid IDs modified file: error=%v, before=%q after=%q", err, fileBefore, fileAfter)
	}

	expectProcess(t, run(pathB, "list"), 0, "No memos found.\n", "")
	fileAfter, err = os.ReadFile(pathA)
	if err != nil || !bytes.Equal(fileBefore, fileAfter) {
		t.Fatalf("switching paths modified file A: %v", err)
	}
	expectProcess(t, run(pathB, "add", "In B"), 0, "", "")
	if memos := readSavedMemos(t, pathB); len(memos) != 1 || memos[0].Title != "In B" {
		t.Errorf("file B = %+v", memos)
	}
	if _, err := os.Stat(filepath.Join(home, ".memo", "memos.json")); !os.IsNotExist(err) {
		t.Errorf("default path was touched: %v", err)
	}
}

func TestMemoProcessPaths(t *testing.T) {
	binary := buildMemo(t)
	for _, tc := range []struct {
		name string
		path string
		want func(cwd, home string) string
	}{
		{"home expansion", "~/memos/test.json", func(_, home string) string { return filepath.Join(home, "memos", "test.json") }},
		{"bare tilde", "~", func(cwd, _ string) string { return filepath.Join(cwd, "~") }},
		{"named user", "~user/memos.json", func(cwd, _ string) string { return filepath.Join(cwd, "~user", "memos.json") }},
		{"relative", "data/memos.json", func(cwd, _ string) string { return filepath.Join(cwd, "data", "memos.json") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cwd, home := t.TempDir(), t.TempDir()
			wantPath := tc.want(cwd, home)
			expectProcess(t, runMemo(t, binary, cwd, home, tc.path, "add", "From path"), 0, "", "")
			memos := readSavedMemos(t, wantPath)
			if len(memos) != 1 || memos[0].Title != "From path" {
				t.Fatalf("memos at %q = %+v", wantPath, memos)
			}
			expectProcess(t, runMemo(t, binary, cwd, home, tc.path, "list"), 0,
				fmt.Sprintf("1 From path\t%s\n", memos[0].CreatedAt.Format("2006-01-02")), "")
		})
	}
}

func TestMemoProcessStorageFailures(t *testing.T) {
	binary := buildMemo(t)
	for _, invalid := range []struct{ name, data string }{
		{"syntax", "[{"},
		{"null document", "null"},
		{"second value", "[] []"},
		{"missing field", `[{"id":1,"title":"Old","body":"","created_at":"2026-01-01T00:00:00Z"}]`},
	} {
		t.Run(invalid.name, func(t *testing.T) {
			cwd, home := t.TempDir(), t.TempDir()
			path := filepath.Join(cwd, "memos.json")
			if err := os.WriteFile(path, []byte(invalid.data), 0600); err != nil {
				t.Fatal(err)
			}
			for _, args := range [][]string{{"list"}, {"add", "New"}} {
				expectProcessError(t, runMemo(t, binary, cwd, home, path, args...), "メモを開くことができませんでした")
			}
			after, err := os.ReadFile(path)
			if err != nil || string(after) != invalid.data {
				t.Fatalf("invalid file changed: data=%q error=%v", after, err)
			}
		})
	}

	t.Run("cannot create parent", func(t *testing.T) {
		cwd, home := t.TempDir(), t.TempDir()
		parent := filepath.Join(cwd, "parent")
		if err := os.WriteFile(parent, []byte("unchanged"), 0600); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(parent, "memos.json")
		// LoadMemos は親が通常ファイルなら ENOTDIR で先に失敗する。
		// SaveMemos の mkdir 失敗はストレージテストで別途検証する。
		expectProcessError(t, runMemo(t, binary, cwd, home, path, "add", "New"), "メモを開くことができませんでした")
		data, err := os.ReadFile(parent)
		if err != nil || string(data) != "unchanged" {
			t.Fatalf("parent modified: %q, %v", data, err)
		}
		if _, err := os.Stat(path); err == nil {
			t.Error("destination was created")
		}
		if matches, err := filepath.Glob(filepath.Join(cwd, ".memos-*.json")); err != nil || len(matches) != 0 {
			t.Errorf("temporary files = %v, error = %v", matches, err)
		}
	})

	t.Run("cannot create dangling symlink parent after loading", func(t *testing.T) {
		cwd, home := t.TempDir(), t.TempDir()
		parent := filepath.Join(cwd, "parent")
		if err := os.Symlink(filepath.Join(cwd, "missing-target"), parent); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(parent, "memos.json")
		expectProcessError(t, runMemo(t, binary, cwd, home, path, "add", "New"), "メモの保存に失敗しました")
		if _, err := os.Stat(filepath.Join(cwd, "missing-target")); !os.IsNotExist(err) {
			t.Errorf("symlink target unexpectedly exists: %v", err)
		}
		if matches, err := filepath.Glob(filepath.Join(cwd, ".memos-*.json")); err != nil || len(matches) != 0 {
			t.Errorf("temporary files = %v, error = %v", matches, err)
		}
	})
}

func TestMemoProcessV010Compatibility(t *testing.T) {
	binary := buildMemo(t)
	fixture, err := os.ReadFile(filepath.Join("..", "..", "testdata", "v0.1.0-memos.json"))
	if err != nil {
		t.Fatal(err)
	}
	cwd, home := t.TempDir(), t.TempDir()
	path := filepath.Join(cwd, "memos.json")
	if err := os.WriteFile(path, fixture, 0600); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) processResult { return runMemo(t, binary, cwd, home, path, args...) }
	legacyList := "4 Legacy Alpha\t2026-01-02\n2 Legacy Beta\t2026-01-03\n"
	expectProcess(t, run("list"), 0, legacyList, "")
	expectProcess(t, run("show", "2"), 0,
		"# Legacy Beta\n\nID: 2\nCreated: 2026-01-03 12:10\nUpdated: 2026-01-04 13:20\n\nlegacy-needle\n", "")
	expectProcess(t, run("search", "legacy-needle"), 0, "2 Legacy Beta\t2026-01-03\n", "")
	expectProcess(t, run("search", "Legacy"), 0, legacyList, "")

	expectProcess(t, run("add", "Title only"), 0, "", "")
	added := readSavedMemos(t, path)
	want := []savedMemo{
		{4, "Legacy Alpha", "", time.Date(2026, 1, 2, 9, 30, 0, 0, time.UTC), time.Date(2026, 1, 2, 9, 30, 0, 0, time.UTC)},
		{2, "Legacy Beta", "legacy-needle", time.Date(2026, 1, 3, 12, 10, 0, 0, time.UTC), time.Date(2026, 1, 4, 13, 20, 0, 0, time.UTC)},
	}
	if len(added) != 3 || added[2].ID != 5 || added[2].Title != "Title only" || added[2].Body != "" ||
		added[2].CreatedAt.IsZero() || !added[2].CreatedAt.Equal(added[2].UpdatedAt) || !reflect.DeepEqual(added[:2], want) {
		t.Fatalf("v0.1.0 add = %+v", added)
	}
	expectProcess(t, run("list"), 0,
		legacyList+fmt.Sprintf("5 Title only\t%s\n", added[2].CreatedAt.Format("2006-01-02")), "")
	expectProcess(t, run("delete", "5"), 0, "", "")
	expectProcess(t, run("list"), 0, legacyList, "")
	expectProcess(t, run("search", "Legacy"), 0, legacyList, "")
	if got := readSavedMemos(t, path); !reflect.DeepEqual(got, want) {
		t.Errorf("v0.1.0 memos after add/delete = %+v, want %+v", got, want)
	}
}
