package argsdef

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/chahal-p/pargs/errors"
)

type ArgType string

const (
	BOOL_ARG   = ArgType("boolean")
	NUMBER_ARG = ArgType("number")
	STRING_ARG = ArgType("string")
)

type ArgDef struct {
	shortName   string
	longName    string
	argType     ArgType
	required    bool
	defaultVals []string
	allowedVals []string
	numRange    []float64
	strRegex    *regexp.Regexp
}

func NewBool(shortName string, longName string) *ArgDef {
	return &ArgDef{
		shortName:   shortName,
		longName:    longName,
		argType:     BOOL_ARG,
		allowedVals: []string{"true", "false"},
	}
}

func NewNumber(shortName string, longName string, allowedValues []float64, lowerBound float64, upperbound float64) *ArgDef {
	var allowed []string
	for _, x := range allowedValues {
		allowed = append(allowed, fmt.Sprint(x))
	}
	return &ArgDef{
		shortName:   shortName,
		longName:    longName,
		argType:     NUMBER_ARG,
		allowedVals: allowed,
		numRange:    []float64{lowerBound, upperbound},
	}
}

func NewString(shortName string, longName string, AllowedValues []string, regex *regexp.Regexp) *ArgDef {
	var allowed []string
	for _, x := range AllowedValues {
		allowed = append(allowed, fmt.Sprint(x))
	}
	return &ArgDef{
		shortName:   shortName,
		longName:    longName,
		argType:     STRING_ARG,
		allowedVals: allowed,
		strRegex:    regex,
	}
}

func (o *ArgDef) ShortName() string {
	return o.shortName
}

func (o *ArgDef) LongName() string {
	return o.longName
}

func (o *ArgDef) Type() ArgType {
	return o.argType
}

func (o *ArgDef) Required() bool {
	return o.required
}

func (o *ArgDef) DefaultValues() []string {
	return o.defaultVals
}

func (o *ArgDef) Validate(val string) *errors.Error {
	switch o.argType {
	case BOOL_ARG:
		_, err := strconv.ParseBool(val)
		if err != nil {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s can not be parsed as boolean", val))
		}
		if len(o.allowedVals) > 0 && slices.Contains(o.allowedVals, val) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, allowed values: %v", val, o.allowedVals))
		}
	case NUMBER_ARG:
		parsedVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s can not be parsed as number", val))
		}
		if len(o.allowedVals) > 0 && slices.Contains(o.allowedVals, val) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, allowed values: %v", val, o.allowedVals))
		}
		if parsedVal < o.numRange[0] || parsedVal > o.numRange[1] {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, number should be with in range %v", val, o.numRange))
		}
	case STRING_ARG:
		if len(o.allowedVals) > 0 && slices.Contains(o.allowedVals, val) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, allowed values: %v", val, o.allowedVals))
		}
		if !o.strRegex.Match([]byte(val)) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, string should be matched by regex %q", val, o.strRegex.String()))
		}
	default:
		return errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("%s is not a valid argument type.", o.argType))
	}
	return nil
}
