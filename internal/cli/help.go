package cli

import (
	"fmt"
	"io"
)

func PrintHelp(stderr io.Writer) {
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
