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

func errorExitFromError(err *errors.Error) {
	errorExit(err.Code().Code(), err.Error())
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

func parse(internalArgs, flagArgs, externalArgs []string) {
	internalPflags := pflags.New(parseDesc)
	internalPflags.Add("d", "description", flagdef.STRING_FLAG, flagdef.Description("Provide desciption content for usage help\n  Specify \\{\\{\\FLAGS\\}\\} formatter to replace it with flags details"))
	internalPflags.Add("h", "help", flagdef.STRING_FLAG, flagdef.Description("Output usage help"))

	flagsPflags := pflags.New(internalPflags.UsageHelp())
	flagsPflags.Add("s", "short", flagdef.STRING_FLAG, flagdef.Description("Short name for flag."))
	flagsPflags.Add("l", "long", flagdef.STRING_FLAG, flagdef.Description("Long name for flag."))
	flagsPflags.Add("t", "type", flagdef.STRING_FLAG, flagdef.Required(), flagdef.Description("Type of flag.\n  Allowed values: string, number, bool"), flagdef.AllowedValues("string", "number", "bool"))
	flagsPflags.Add("d", "description", flagdef.STRING_FLAG, flagdef.Required(), flagdef.Description("Description of the flag."))
	flagsPflags.Add("r", "required", flagdef.STRING_FLAG, flagdef.DefaultValues("false"), flagdef.Description("If a flag is required"))
	flagsPflags.Add("", "default", flagdef.STRING_FLAG, flagdef.Description("Default values\n  (Can be specified multiple times)."))
	flagsPflags.Add("a", "allowed", flagdef.STRING_FLAG, flagdef.Description("Allowed Values\n  (Can be specified multiple times)."))
	flagsPflags.Add("", "regex", flagdef.STRING_FLAG, flagdef.Description("Regex for string validatin\n  (Only applicable to --type=string)."))
	if hasHelpFlag(internalArgs) {
		println(flagsPflags.UsageHelp())
		return
	}
	if err := internalPflags.Parse(internalArgs); err != nil {
		errorExitFromError(err)
	}
	externalPflags := pflags.New(internalPflags.Get("description")[0])
	for _, args := range splitArgs(flagArgs, "--") {
		if err := flagsPflags.Parse(args); err != nil {
			errorExitFromError(err)
		}
		t, err := flagdef.TypeFromString(flagsPflags.Get("type")[0])
		if err != nil {
			errorExitFromError(err)
		}
		var opts []flagdef.Option
		opts = append(opts, flagdef.AllowedValues(flagsPflags.Get("allowed")...))
		opts = append(opts, flagdef.DefaultValues(flagsPflags.Get("default")...))
		opts = append(opts, flagdef.StringRegex(flagsPflags.Get("regex")[0]))
		opts = append(opts, flagdef.Description(flagsPflags.Get("description")[0]))
		externalPflags.Add(flagsPflags.Get("short")[0], flagsPflags.Get("long")[0], t, opts...)
	}
	if err := externalPflags.Parse(externalArgs); err != nil {
		errorExitFromError(err)
	}

}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func splitArgs(args []string, cutstring string) [][]string {
	remainingArgs := args[:]
	result := make([][]string, 0)
	for {
		delIndx := slices.Index(remainingArgs, cutstring)
		if delIndx == -1 {
			result = append(result, remainingArgs)
			return result
		} else {
			result = append(result, remainingArgs[:delIndx])
			remainingArgs = remainingArgs[delIndx+1:]
		}
	}
}

func main() {
	if hasHelpFlag(os.Args[1:2]) {
		println(rootHelp)
		return
	}
	subCmd := os.Args[1]
	args := os.Args[2:]
	splittedArgs := splitArgs(args, "----")
	internalArgs := splittedArgs[0]
	flagArgs := make([]string, 0)
	externalArgs := make([]string, 0)
	if len(splittedArgs) > 1 {
		flagArgs = splittedArgs[1]
	}
	if len(splittedArgs) > 2 {
		externalArgs = splittedArgs[2]
	}

	if len(internalArgs) == 0 {
		errorExit(errors.INVALID_USAGE.Code(), "No flags provided")
	}

	switch subCmd {
	case "parse":
		parse(internalArgs, flagArgs, externalArgs)
	case "get":

	case "unparsed":

	default:
		errorExit(errors.INVALID_USAGE.Code(), "Unrecognized subcommand.")
	}
}
