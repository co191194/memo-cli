package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	app := App{
		Command: &MemoCommandImpl{
			MemoPath:     "~/.memo/memos.json",
			TimeProvider: &RealTimeProvider{},
			MemoOperator: &MemoOperatorImpl{},
		},
	}
	exitCode := app.run(
		os.Stdout,
		os.Stderr,
		os.Args[1:],
	)

	os.Exit(exitCode)
}

type App struct {
	Command MemoCommand
}

func (app *App) run(
	stdout io.Writer,
	stderr io.Writer,
	args []string,
) int {
	if len(args) < 1 {
		printHelp(stderr)
		return 1
	}

	var exitCode int

	switch args[0] {
	case "add":
		exitCode = app.Command.AddMemo(stdout, stderr, args[1:])
	case "list":
		exitCode = app.Command.ListMemos(stdout, stderr, args[1:])
	case "show":
		exitCode = app.Command.ShowMemo(stdout, stderr, args[1:])
	case "search":
		exitCode = app.Command.SearchMemos(stdout, stderr, args[1:])
	case "delete":
		exitCode = app.Command.DeleteMemo(stdout, stderr, args[1:])
	default:
		printHelp(stderr)
		return 1
	}

	return exitCode
}

func printHelp(stderr io.Writer) {
	fmt.Fprintln(stderr, "Usage:")
	fmt.Fprintln(stderr, "  memo <command> [arguments]")
	fmt.Fprintln(stderr)
	fmt.Fprintln(stderr, "Commands:")
	fmt.Fprintln(stderr, "  add     Add a new memo")
	fmt.Fprintln(stderr, "  list    List memos")
	fmt.Fprintln(stderr, "  show    Show a memo")
	fmt.Fprintln(stderr, "  search  Search memos")
	fmt.Fprintln(stderr, "  delete  delete a memo")
}
