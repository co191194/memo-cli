package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
	case "list":
	case "show":
	case "search":
	case "delete":
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  memo <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add     Add a new memo")
	fmt.Println("  list    List memos")
	fmt.Println("  show    Show a memo")
	fmt.Println("  search  Search memos")
	fmt.Println("  delete  delete a memo")
}
