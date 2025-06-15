package parse

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chahal-p/pflags/errors"
	"github.com/chahal-p/pflags/flagdef"
)

type Result struct {
	FlagValuesForID map[int]([]string) `json:"flagValuesForID"`
	FlagsNameToID   map[string]int     `json:"flagsNameToID"`
	NonFlagArgs     []string           `json:"nonFlagArgs"`
}

func Get(name string, result *Result) ([]string, *errors.Error) {
	if result == nil {
		return nil, errors.NewError(errors.INVALID_USAGE, "Parsed result is nil.")
	}
	var vals []string
	id, ok := result.FlagsNameToID[name]
	if !ok {
		return nil, errors.NewError(errors.NOT_FOUND, fmt.Sprintf("Flag %s not found in parsed result.", name))
	}
	vals, ok = result.FlagValuesForID[id]
	if !ok {
		return nil, errors.NewError(errors.NOT_FOUND, fmt.Sprintf("ID %d for flag %s can not be found in parsed result.", id, name))
	}
	return vals, nil
}

func NonFlagArgs(result *Result) []string {
	if result == nil {
		return nil
	}
	return result.NonFlagArgs[:]
}

func Parse(flagDef []*flagdef.FlagDef, cmdArgs []string, allowUnrecognizedFlags bool) (*Result, *errors.Error) {
	result := &Result{
		FlagValuesForID: make(map[int][]string),
		FlagsNameToID:   make(map[string]int),
		NonFlagArgs:     make([]string, 0),
	}
	var flagVals []string = nil
	var err *errors.Error = nil
	remaining := slices.Clone(cmdArgs)
	for id, flag := range flagDef {
		if flag.ShortName() != "" {
			result.FlagsNameToID[flag.ShortName()] = id
		}
		if flag.LongName() != "" {
			result.FlagsNameToID[flag.LongName()] = id
		}
		flagVals, remaining, err = flagValue(flag, remaining)
		if err != nil {
			return nil, err
		}
		result.FlagValuesForID[id] = flagVals
	}
	if !allowUnrecognizedFlags {
		for _, arg := range remaining {
			if strings.Trim(arg, "-") != "" {
				if strings.HasPrefix(arg, "-") {
					return nil, errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("Unrecognized flag: %s", arg))
				}
			}
		}
	}
	result.NonFlagArgs = remaining
	return result, nil
}

func flagValue(flag *flagdef.FlagDef, cmdArgs []string) ([]string, []string, *errors.Error) {
	var values []string
	var remaining []string
	size := len(cmdArgs)
	for i := 0; i < size; {
		fn := cmdArgs[i]
		if cmdArgs[i] == fmt.Sprintf("-%s", flag.ShortName()) || cmdArgs[i] == fmt.Sprintf("--%s", flag.LongName()) {
			i++
			if flag.Type() == flagdef.BOOL_FLAG {
				var err *errors.Error
				if i < len(cmdArgs) {
					err = flag.Validate(cmdArgs[i])
				}
				if i < len(cmdArgs) && err == nil {
					values = append(values, cmdArgs[i])
					i++
				} else {
					values = append(values, "true")
				}
			} else if i < len(cmdArgs) && !strings.HasPrefix(cmdArgs[i], "-") {
				if err := flag.Validate(cmdArgs[i]); err != nil {
					return nil, cmdArgs, err
				}
				values = append(values, cmdArgs[i])
				i++
			} else {
				return nil, cmdArgs, errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("No value provided for flag %s", fn))
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
		if flag.Required() {
			msg := fmt.Sprintf("Required flag missing: -%s/--%s", flag.ShortName(), flag.LongName())
			if flag.ShortName() == "" {
				msg = fmt.Sprintf("Required flag missing: --%s", flag.LongName())
			}
			if flag.LongName() == "" {
				msg = fmt.Sprintf("Required flag missing: -%s", flag.ShortName())
			}
			return nil, cmdArgs, errors.NewError(errors.INVALID_USAGE, msg)
		}
		if len(flag.DefaultValues()) > 0 {
			values = flag.DefaultValues()
		}
	}
	return values, remaining, nil
}
