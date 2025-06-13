package parse

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chahal-p/pflags/errors"
	"github.com/chahal-p/pflags/flagdef"
)

type Result struct {
	FlagValuesForID map[int][]string
	FlagsNameToID   map[string]int
	Unparsed        []string
}

func Parse(flagDef []flagdef.FlagDef, cmdArgs []string, opts ...Option) (*Result, *errors.Error) {
	conf := &config{}
	for _, opt := range opts {
		opt.set(conf)
	}
	result := &Result{
		FlagValuesForID: make(map[int][]string),
		FlagsNameToID:   make(map[string]int),
		Unparsed:        make([]string, 0),
	}
	var flagVals []string = nil
	var err *errors.Error = nil
	remaining := slices.Clone(cmdArgs)
	for id, flag := range flagDef {
		result.FlagsNameToID[flag.ShortName()] = id
		result.FlagsNameToID[flag.LongName()] = id
		flagVals, remaining, err = flagValue(flag, remaining)
		if err != nil {
			return nil, err
		}
		result.FlagValuesForID[id] = flagVals
	}
	if !conf.noErrorForUnrecognizedFlag {
		for _, arg := range remaining {
			if strings.Trim(arg, "-") != "" {
				if strings.HasPrefix(arg, "--") {
					return nil, errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("Unrecognized flag: %s", arg))
				}
			}
		}
	}
	result.Unparsed = remaining
	return result, nil
}

func flagValue(flag flagdef.FlagDef, cmdArgs []string) ([]string, []string, *errors.Error) {
	var values []string
	var remaining []string
	size := len(cmdArgs)
	for i := 0; i < size; {
		if cmdArgs[i] == fmt.Sprintf("-%s", flag.ShortName()) || cmdArgs[i] == fmt.Sprintf("--%s", flag.LongName()) {
			i++
			if flag.Type() != flagdef.BOOL_FLAG {
				if err := flag.Validate(cmdArgs[i]); err == nil {
					values = append(values, cmdArgs[i])
					i++
				} else {
					values = append(values, "true")
				}
			} else {
				if err := flag.Validate(cmdArgs[i]); err != nil {
					return nil, cmdArgs, err
				}
				values = append(values, cmdArgs[i])
				i++
			}
		} else if strings.HasPrefix(cmdArgs[i], fmt.Sprintf("-%s=", flag.ShortName())) || strings.HasPrefix(cmdArgs[i], fmt.Sprintf("--%s=", flag.LongName())) {
			eqIndx := strings.Index(cmdArgs[i], "=")
			value := cmdArgs[i][eqIndx+1:]
			if err := flag.Validate(value); err != nil {
				return nil, cmdArgs, err
			}
			values = append(values, value)
			i++
		} else {
			remaining = append(remaining, cmdArgs[i])
			i++
		}
	}
	if len(values) == 0 {
		if len(flag.DefaultValues()) > 0 {
			values = flag.DefaultValues()
		} else if flag.Required() {
			msg := fmt.Sprintf("Required flag missing: -%s/--%s", flag.ShortName(), flag.LongName())
			if flag.ShortName() == "" {
				msg = fmt.Sprintf("Required flag missing: --%s", flag.LongName())
			}
			if flag.LongName() == "" {
				msg = fmt.Sprintf("Required flag missing: -%s", flag.ShortName())
			}
			return nil, cmdArgs, errors.NewError(errors.INVALID_USAGE, msg)
		}
	}
	return values, remaining, nil
}

type config struct {
	noErrorForUnrecognizedFlag bool
}

type Option interface {
	set(conf *config)
}

type errForUnrecongnized struct{}

func (*errForUnrecongnized) set(conf *config) {
	conf.noErrorForUnrecognizedFlag = true
}

func NoErrorForUnrecognizedFlag() Option {
	return &errForUnrecongnized{}
}
