package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type fakeMemoCommand struct {
	MemoPath     string
	TimeProvider TimeProvider
	Command      MemoCommand
}

func (cmd *fakeMemoCommand) AddMemo(stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprint(stdout, "Called AddMemo()")

	if len(args) != 1 || args[0] == "add" {
		fmt.Fprint(stderr, "Failed AddMemo()")
		return 1
	}

	return 0
}
func (cmd *fakeMemoCommand) ListMemos(stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprint(stdout, "Called ListMemos()")

	if len(args) != 0 {
		fmt.Fprint(stderr, "Failed ListMemos()")
		return 1
	}

	return 0
}
func (cmd *fakeMemoCommand) ShowMemo(stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprint(stdout, "Called ShowMemo()")

	if len(args) != 1 || args[0] == "show" {
		fmt.Fprint(stderr, "Failed ShowMemo()")
		return 1
	}

	return 0
}
func (cmd *fakeMemoCommand) SearchMemos(stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprint(stdout, "Called SearchMemos()")

	if len(args) != 1 || args[0] == "search" {
		fmt.Fprint(stderr, "Failed SearchMemos()")
		return 1
	}
	return 0
}
func (cmd *fakeMemoCommand) DeleteMemo(stdout io.Writer, stderr io.Writer, args []string) int {
	fmt.Fprint(stdout, "Called DeleteMemo()")

	if len(args) != 1 || args[0] == "delete" {
		fmt.Fprint(stderr, "Failed DeleteMemo()")
		return 1
	}
	return 0
}

func TestRun(t *testing.T) {
	app := App{
		Command: &fakeMemoCommand{},
	}

	testCases := []struct {
		testName string
		args     []string
		expected string
	}{
		{"AddMemo()を呼ぶ", []string{"add", "test"}, "Called AddMemo()"},
		{"ListMemos()を呼ぶ", []string{"list"}, "Called ListMemos()"},
		{"ShowMemo()を呼ぶ", []string{"show", "1"}, "Called ShowMemo()"},
		{"SearchMemos()を呼ぶ", []string{"search", "Jack"}, "Called SearchMemos()"},
		{"DeleteMemo()を呼ぶ", []string{"delete", "3"}, "Called DeleteMemo()"},
	}

	for _, tc := range testCases {

		t.Run(tc.testName, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if exitCode := app.run(&stdout, &stderr, tc.args); exitCode != 0 {
				t.Fatalf("actual exitCode = %d, expected = 0", exitCode)
			}

			if stdout.String() != tc.expected {
				t.Errorf("actual = %q, expected = %q", stdout.String(), tc.expected)
			}
		})
	}

}

func TestRun_PrintHelp(t *testing.T) {
	app := App{
		Command: &fakeMemoCommand{},
	}

	testCases := []struct {
		testName string
		args     []string
	}{
		{"コマンドが未入力の場合", []string{}},
		{"未定義のコマンドの場合", []string{"unknown"}},
	}

	var sb strings.Builder
	sb.WriteString("Usage:\n")
	sb.WriteString("  memo <command> [arguments]\n")
	sb.WriteString("\n")
	sb.WriteString("Commands:\n")
	sb.WriteString("  add     Add a new memo\n")
	sb.WriteString("  list    List memos\n")
	sb.WriteString("  show    Show a memo\n")
	sb.WriteString("  search  Search memos\n")
	sb.WriteString("  delete  delete a memo\n")

	expected := sb.String()

	for _, tc := range testCases {

		t.Run(tc.testName, func(t *testing.T) {

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if exitCode := app.run(&stdout, &stderr, tc.args); exitCode != 1 {
				t.Fatalf("actual exitCode = %d, expected = 1", exitCode)
			}

			if stderr.String() != expected {
				t.Errorf("actual = %q, expected = %q", stdout.String(), expected)
			}
		})
	}
}

func TestRun_FailedCommand(t *testing.T) {

	app := App{
		Command: &fakeMemoCommand{},
	}

	testCases := []struct {
		testName string
		args     []string
		expected string
	}{
		{"AddMemo()が失敗", []string{"add", "test", "test"}, "Failed AddMemo()"},
		{"ListMemos()が失敗", []string{"list", "aaaa"}, "Failed ListMemos()"},
		{"ShowMemo()が失敗", []string{"show"}, "Failed ShowMemo()"},
		{"SearchMemos()が失敗", []string{"search", "Jack", "Joe"}, "Failed SearchMemos()"},
		{"DeleteMemo()が失敗", []string{"delete", "3", "2"}, "Failed DeleteMemo()"},
	}

	for _, tc := range testCases {

		t.Run(tc.testName, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if exitCode := app.run(&stdout, &stderr, tc.args); exitCode != 1 {
				t.Fatalf("actual exitCode = %d, expected = 1", exitCode)
			}

			if stderr.String() != tc.expected {
				t.Errorf("actual = %q, expected = %q", stderr.String(), tc.expected)
			}
		})
	}
}
