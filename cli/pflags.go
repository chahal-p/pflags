package main

import (
	"encoding/base64"
	"fmt"
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
	if err == nil {
		return
	}
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

var getDesc = strings.Trim(`
pflags get:
  Get value(s) of a specific flag from parsed results.
  pflags get <flags> "$parsedData"

  Example:
    pflags get --name abc "$parsedData" 

  FLAGS:
    {FLAGS}
`, "\n")

var unparsedDesc = strings.Trim(`
pflags unparsed:
  Get non flag args
  pflags unparsed "$parsedData"

  FLAGS:
    {FLAGS}
`, "\n")

func flagGet(f *pflags.Pflags, name string) []string {
	res, err := f.Get(name)
	if err != nil {
		errorExitFromError(err)
	}
	return res
}

func parseSubCommand(internalArgs, flagArgs, externalArgs []string) {
	internalPflags := pflags.New(parseDesc)
	errorExitFromError(internalPflags.Add("d", "description", flagdef.STRING_FLAG, flagdef.DefaultValues(""), flagdef.Description("Provide desciption content for usage help\n  Specify \\{\\{\\FLAGS\\}\\} formatter to replace it with flags details")))
	errorExitFromError(internalPflags.Add("", "unrecognized-flags", flagdef.STRING_FLAG, flagdef.AllowedValues("allow", "error"), flagdef.Description("Unrecognized flags: accepted values 'allow' or 'error'")))
	errorExitFromError(internalPflags.Add("h", "help", flagdef.BOOL_FLAG, flagdef.Description("Output usage help")))

	flagsPflags := pflags.New(internalPflags.UsageHelp())
	errorExitFromError(flagsPflags.Add("s", "short", flagdef.STRING_FLAG, flagdef.DefaultValues(""), flagdef.Description("Short name for flag.")))
	errorExitFromError(flagsPflags.Add("l", "long", flagdef.STRING_FLAG, flagdef.DefaultValues(""), flagdef.Description("Long name for flag.")))
	errorExitFromError(flagsPflags.Add("t", "type", flagdef.STRING_FLAG, flagdef.Required(true), flagdef.Description("Type of flag.\n  Allowed values: string, number, bool"), flagdef.AllowedValues("string", "number", "bool")))
	errorExitFromError(flagsPflags.Add("d", "description", flagdef.STRING_FLAG, flagdef.Required(true), flagdef.Description("Description of the flag.")))
	errorExitFromError(flagsPflags.Add("r", "required", flagdef.BOOL_FLAG, flagdef.DefaultValues("false"), flagdef.Description("If a flag is required")))
	errorExitFromError(flagsPflags.Add("", "default", flagdef.STRING_FLAG, flagdef.Description("Default values\n  (Can be specified multiple times).")))
	errorExitFromError(flagsPflags.Add("a", "allowed", flagdef.STRING_FLAG, flagdef.Description("Allowed Values\n  (Can be specified multiple times).")))
	errorExitFromError(flagsPflags.Add("", "regex", flagdef.STRING_FLAG, flagdef.DefaultValues(""), flagdef.Description("Regex for string validatin\n  (Only applicable to --type=string).")))
	if hasHelpFlag(internalArgs) {
		println(flagsPflags.UsageHelp())
		return
	}
	if len(flagArgs) == 0 {
		errorExit(errors.INVALID_USAGE.Code(), "No flags provided")
	}
	if err := internalPflags.Parse(internalArgs); err != nil {
		errorExitFromError(err)
	}
	var opts []pflags.Option
	if flagGet(internalPflags, "unrecognized-flags")[0] == "allow" {
		opts = append(opts, pflags.AllowUnrecognizedFlags())
	}

	externalPflags := pflags.New(flagGet(internalPflags, "description")[0], opts...)
	for _, args := range splitArgs(flagArgs, "--") {
		if err := flagsPflags.Parse(args); err != nil {
			errorExitFromError(err)
		}
		t, err := flagdef.TypeFromString(flagGet(flagsPflags, "type")[0])
		if err != nil {
			errorExitFromError(err)
		}
		var opts []flagdef.Option
		opts = append(opts, flagdef.Description(flagGet(flagsPflags, "description")[0]))
		opts = append(opts, flagdef.DefaultValues(flagGet(flagsPflags, "default")...))
		opts = append(opts, flagdef.AllowedValues(flagGet(flagsPflags, "allowed")...))
		opts = append(opts, flagdef.StringRegex(flagGet(flagsPflags, "regex")[0]))
		if flagGet(flagsPflags, "required")[0] == "true" {
			opts = append(opts, flagdef.Required(true))
		}
		errorExitFromError(externalPflags.Add(flagGet(flagsPflags, "short")[0], flagGet(flagsPflags, "long")[0], t, opts...))
	}
	if hasHelpFlag(externalArgs) {
		println(externalPflags.UsageHelp())
		os.Exit(errors.USAGE_HELP_REQUESTED.Code())
		return
	}
	if err := externalPflags.Parse(externalArgs); err != nil {
		errorExitFromError(err)
	}
	print(base64.StdEncoding.EncodeToString(externalPflags.Parsed()))
	// fmt.Printf("\n%v", flagGet(externalPflags, "a"))
}

func getSubCommand(args []string) {
	flags := pflags.New(getDesc)
	errorExitFromError(flags.Add("n", "name", flagdef.STRING_FLAG, flagdef.Required(true), flagdef.Description("Name of flag, any one of either short or long name can be provided.")))
	errorExitFromError(flags.Add("h", "help", flagdef.BOOL_FLAG, flagdef.Description("Output usage details.")))
	if hasHelpFlag(args) {
		println(flags.UsageHelp())
		os.Exit(errors.USAGE_HELP_REQUESTED.Code())
		return
	}
	errorExitFromError(flags.Parse(args))
	nonFlagArgs := flags.NonFlagArgs()
	if len(nonFlagArgs) == 0 {
		errorExit(errors.INVALID_USAGE.Code(), "Parsed args data is not provided.")
	} else if len(nonFlagArgs) > 1 {
		errorExit(errors.INVALID_USAGE.Code(), "Only 1 non-flag arg should be given.")
	}
	flagName := flagGet(flags, "name")[0]
	parsedData, err := base64.StdEncoding.DecodeString(nonFlagArgs[0])
	if err != nil {
		errorExit(errors.INTERNAL_ERROR.Code(), err.Error())
	}
	vals, gErr := pflags.GetFromParsedBytes(flagName, parsedData)
	if gErr != nil {
		errorExitFromError(gErr)
	}
	for _, val := range vals {
		println(val)
	}
}

func unparsedSubCommand(args []string) {
	flags := pflags.New(unparsedDesc)
	errorExitFromError(flags.Add("h", "help", flagdef.BOOL_FLAG, flagdef.Description("Output usage details.")))
	if hasHelpFlag(args) {
		println(flags.UsageHelp())
		os.Exit(errors.USAGE_HELP_REQUESTED.Code())
		return
	}
	if hasHelpFlag(args) {
		println(flags.UsageHelp())
		os.Exit(errors.USAGE_HELP_REQUESTED.Code())
		return
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

	switch subCmd {
	case "parse":
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
		parseSubCommand(internalArgs, flagArgs, externalArgs)
	case "get":
		getSubCommand(args)
	case "unparsed":
		unparsedSubCommand(args)
	default:
		errorExit(errors.INVALID_USAGE.Code(), fmt.Sprintf("Unrecognized subcommand: %s", subCmd))
	}
}
