package main

import (
	"os"

	"github.com/co191194/memo-cli/internal/cli"
	"github.com/co191194/memo-cli/internal/storage"
)

func main() {
	app := cli.App{
		Command: &cli.MemoCommandImpl{
			MemoPath:        "~/.memo/memos.json",
			TimeProvider:    &cli.RealTimeProvider{},
			StorageOperator: &storage.StorageOperatorImpl{},
		},
	}
	exitCode := app.Run(
		os.Stdout,
		os.Stderr,
		os.Args[1:],
	)

	os.Exit(exitCode)
}
