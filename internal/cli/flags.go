package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

func newCommandFlagSet(
	name string,
	usage string,
	output io.Writer,
) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(output)

	fs.Usage = func() {
		fmt.Fprintln(output, "Usage:")
		fmt.Fprintln(output, "  "+usage)
		fs.PrintDefaults()
	}

	return fs
}

func parseCommandArgs(
	fs *flag.FlagSet,
	args []string,
	positionalCount int,
) ([]string, error) {
	separatorIndex := findEndOfOptions(args)

	if separatorIndex >= 0 {
		if err := fs.Parse(args[:separatorIndex]); err != nil {
			return nil, err
		}

		// `--`より前に位置引数があった場合
		if fs.NArg() != 0 {
			fs.Usage()

			return nil, errors.New(
				"only options may appear before --",
			)
		}

		return validatePositionals(
			fs,
			args[separatorIndex+1:],
			positionalCount,
		)
	}

	if positionalCount > 0 &&
		len(args) > 0 &&
		isOptionArgument(args[0]) {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}

		fmt.Fprintln(
			fs.Output(),
			"位置引数はオプションより前に指定してください",
		)
		fs.Usage()

		return nil, errors.New(
			"positional arguments must precede options",
		)
	}

	if len(args) < positionalCount {
		fs.Usage()
		return nil, fmt.Errorf(
			"%s requires %d positional arguments",
			fs.Name(),
			positionalCount,
		)
	}

	positionals := args[:positionalCount]

	if err := fs.Parse(args[positionalCount:]); err != nil {
		return nil, err
	}

	if fs.NArg() != 0 {
		fmt.Fprintln(fs.Output(), "引数が多すぎます")
		fs.Usage()

		return nil, fmt.Errorf(
			"unexpected arguments: %v",
			fs.Args(),
		)
	}

	return positionals, nil
}

func isOptionArgument(arg string) bool {
	return len(arg) > 1 && arg[0] == '-'
}

func findEndOfOptions(args []string) int {
	for i, arg := range args {
		if arg == "--" {
			return i
		}
	}

	return -1
}

func validatePositionals(
	fs *flag.FlagSet,
	args []string,
	expectedCount int,
) ([]string, error) {
	if len(args) < expectedCount {
		fs.Usage()

		return nil, fmt.Errorf(
			"%s requires %d positional arguments",
			fs.Name(),
			expectedCount,
		)
	}

	if len(args) > expectedCount {
		fmt.Fprintln(fs.Output(), "引数が多すぎます")
		fs.Usage()

		return nil, fmt.Errorf(
			"unexpected arguments: %v",
			args[expectedCount:],
		)
	}

	return args, nil
}
