package parse

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chahal-p/pargs/argsdef"
	"github.com/chahal-p/pargs/errors"
)

type Result struct {
	ArgsValuesForID map[int][]string
	ArgsNameToID    map[string]int
	Unparsed        []string
}

func Parse(argsDef []argsdef.ArgDef, cmdArgs []string, opts ...Option) (*Result, *errors.Error) {
	conf := &config{}
	for _, opt := range opts {
		opt.set(conf)
	}
	result := &Result{
		ArgsValuesForID: make(map[int][]string),
		ArgsNameToID:    make(map[string]int),
		Unparsed:        make([]string, 0),
	}
	var argValues []string = nil
	var err *errors.Error = nil
	remaining := slices.Clone(cmdArgs)
	for id, arg := range argsDef {
		result.ArgsNameToID[arg.ShortName()] = id
		result.ArgsNameToID[arg.LongName()] = id
		argValues, remaining, err = argValue(arg, remaining)
		if err != nil {
			return nil, err
		}
		result.ArgsValuesForID[id] = argValues
	}
	if conf.errForUnrecongnized {
		for _, arg := range remaining {
			if strings.Trim(arg, "-") != "" {
				if strings.HasPrefix(arg, "--") {
					return nil, errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("Unrecognized arg: %s", arg))
				}
			}
		}
	}
	result.Unparsed = remaining
	return result, nil
}

func argValue(arg argsdef.ArgDef, args []string) ([]string, []string, *errors.Error) {
	var values []string
	var remaining []string
	size := len(args)
	for i := 0; i < size; {
		if args[i] == fmt.Sprintf("-%s", arg.ShortName()) || args[i] == fmt.Sprintf("--%s", arg.LongName()) {
			i++
			if arg.Type() != argsdef.BOOL_ARG {
				if err := arg.Validate(args[i]); err == nil {
					values = append(values, args[i])
					i++
				} else {
					values = append(values, "true")
				}
			} else {
				if err := arg.Validate(args[i]); err != nil {
					return nil, args, err
				}
				values = append(values, args[i])
				i++
			}
		} else if strings.HasPrefix(args[i], fmt.Sprintf("-%s=", arg.ShortName())) || strings.HasPrefix(args[i], fmt.Sprintf("--%s=", arg.LongName())) {
			eqIndx := strings.Index(args[i], "=")
			value := args[i][eqIndx+1:]
			if err := arg.Validate(value); err != nil {
				return nil, args, err
			}
			values = append(values, value)
			i++
		} else {
			remaining = append(remaining, args[i])
			i++
		}
	}
	if len(values) == 0 {
		if len(arg.DefaultValues()) > 0 {
			values = arg.DefaultValues()
		} else if arg.Required() {
			msg := fmt.Sprintf("Required args missing: -%s/--%s", arg.ShortName(), arg.LongName())
			if arg.ShortName() == "" {
				msg = fmt.Sprintf("Required args missing: --%s", arg.LongName())
			}
			if arg.LongName() == "" {
				msg = fmt.Sprintf("Required args missing: -%s", arg.ShortName())
			}
			return nil, args, errors.NewError(errors.INVALID_USAGE, msg)
		}
	}
	return values, remaining, nil
}

type config struct {
	errForUnrecongnized bool
}

type Option interface {
	set(conf *config)
}

type errForUnrecongnized struct{}

func (*errForUnrecongnized) set(conf *config) {
	conf.errForUnrecongnized = true
}

func ErrForUnrecongnizedArgs() Option {
	return &errForUnrecongnized{}
}
