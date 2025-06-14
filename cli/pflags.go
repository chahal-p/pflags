package main

import (
	"os"
	"slices"
	"strings"

	"github.com/chahal-p/pflags"
	"github.com/chahal-p/pflags/errors"
	"github.com/chahal-p/pflags/flagdef"
)

func errorExit(code int, msg string) {
	println(msg)
	os.Exit(code)
}

var rootHelp = strings.Trim(`
pflags: A tool to parse and extract flags from command line arguments.
  It supports below sub commands
    parse:
      It parses the args and output base64 of parsed args.
      Use, pflags parse -h/--help to get more details about parse.
    get:
      It outputs values of give flag name, name can be either short or long flag name.
      Use, pflags get -h/--help to get more details about parse.
    unparsed:
      It outputs unparsed values args.
      Use, pflags get -h/--help to get more details about parse.

  Flags:
    -h,  --help    Outputs the help
    `, "\n")

var parseDesc = strings.Trim(`
pflags parse:
  Parses the args as per given flag definitions and outputs the base64 of parsed results.
  Parsed results can be used for 'pflags get' or 'pflags unparsed' commands.

  pflags parse <FLAGS 1> ---- <FLAGS 2> -- <FLAGS 2> -- <and more...> ---- <Args to be parsed>

  Example:
    pflags parse --description "Testing foo" ---- \
        --short "a" --long "abc" --type string --required --allowed foo --allowed bar -- \
        --short "f" --long "fgh" --type number --default 123 -- \
        <and more flags...> ---- "$@"

  FLAGS 1:
    {FLAGS}

  FLAGS 2:
    {FLAGS}
`, "\n")

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func splitArgs(args []string) [][]string {
	remainingArgs := args[:]
	result := make([][]string, 0)
	for range 3 {
		delIndx := slices.Index(remainingArgs, "----")
		if delIndx == -1 {
			result = append(result, remainingArgs[:])
			remainingArgs = []string{}
		} else {
			result = append(result, remainingArgs[:delIndx])
			remainingArgs = remainingArgs[delIndx+1:]
		}
	}
	return result
}

func main() {
	if hasHelpFlag(os.Args[1:2]) {
		println(rootHelp)
		return
	}
	subCmd := os.Args[1]
	args := os.Args[2:]
	splittedArgs := splitArgs(args)
	internalArgs := splittedArgs[0]
	flagArgs := splittedArgs[1]
	externalArgs := splittedArgs[2]

	if len(internalArgs) == 0 {
		errorExit(errors.INVALID_USAGE.Code(), "No flags provided")
	}

	_ = flagArgs
	_ = externalArgs
	switch subCmd {
	case "parse":
		internalFlags := pflags.New(parseDesc)
		internalFlags.Add("d", "description", flagdef.STRING_FLAG, flagdef.Description("Provide desciption content for usage help"))
		internalFlags.Add("h", "help", flagdef.STRING_FLAG, flagdef.Description("Output usage help"))

		flagsPflags := pflags.New(internalFlags.UsageHelp())
		flagsPflags.Add("s", "short", flagdef.STRING_FLAG, flagdef.Description("Short name for flag."))
		flagsPflags.Add("l", "long", flagdef.STRING_FLAG, flagdef.Description("Long name for flag."))
		flagsPflags.Add("a", "allowed", flagdef.STRING_FLAG, flagdef.Description("Allowed Values\n  Can be specified multiple times)."))
		flagsPflags.Add("d", "default", flagdef.STRING_FLAG, flagdef.Description("Default value."))

		if hasHelpFlag(internalArgs) {
			println(flagsPflags.UsageHelp())
			return
		}
	case "get":

	case "unparsed":

	default:
		errorExit(errors.INVALID_USAGE.Code(), "Unrecognized subcommand.")
	}
}
