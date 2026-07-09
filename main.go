package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) == 2 || strings.TrimSpace(os.Args[2]) == "" {
			fmt.Println("タイトルを入力してください")
			os.Exit(1)
		}
		if len(os.Args) != 3 {
			fmt.Println("Usage:")
			fmt.Println("  memo add <title>")
			os.Exit(1)
		}
		AddMemo(os.Args[2])
	case "list":
		if len(os.Args) != 2 {
			fmt.Println("Usage:")
			fmt.Println("  memo list")
			os.Exit(1)
		}
		ListMemos()
	case "show":
		fmt.Println("not implemented")
		os.Exit(1)
	case "search":
		fmt.Println("not implemented")
		os.Exit(1)
	case "delete":
		fmt.Println("not implemented")
		os.Exit(1)
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
