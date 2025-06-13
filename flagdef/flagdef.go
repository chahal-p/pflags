package flagdef

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/chahal-p/pflags/errors"
)

type FlagType string

const (
	BOOL_FLAG   = FlagType("boolean")
	NUMBER_FLAG = FlagType("number")
	STRING_FLAG = FlagType("string")
)

type FlagDef struct {
	shortName   string
	longName    string
	flagType    FlagType
	required    bool
	defaultVals []string
	allowedVals []string
	numRange    []float64
	strRegex    *regexp.Regexp
}

func NewBool(shortName string, longName string) *FlagDef {
	return &FlagDef{
		shortName:   shortName,
		longName:    longName,
		flagType:    BOOL_FLAG,
		allowedVals: []string{"true", "false"},
	}
}

func NewNumber(shortName string, longName string, allowedValues []float64, lowerBound float64, upperbound float64) *FlagDef {
	var allowed []string
	for _, x := range allowedValues {
		allowed = append(allowed, fmt.Sprint(x))
	}
	return &FlagDef{
		shortName:   shortName,
		longName:    longName,
		flagType:    NUMBER_FLAG,
		allowedVals: allowed,
		numRange:    []float64{lowerBound, upperbound},
	}
}

func NewString(shortName string, longName string, AllowedValues []string, regex *regexp.Regexp) *FlagDef {
	var allowed []string
	for _, x := range AllowedValues {
		allowed = append(allowed, fmt.Sprint(x))
	}
	return &FlagDef{
		shortName:   shortName,
		longName:    longName,
		flagType:    STRING_FLAG,
		allowedVals: allowed,
		strRegex:    regex,
	}
}

func (o *FlagDef) ShortName() string {
	return o.shortName
}

func (o *FlagDef) LongName() string {
	return o.longName
}

func (o *FlagDef) Type() FlagType {
	return o.flagType
}

func (o *FlagDef) Required() bool {
	return o.required
}

func (o *FlagDef) DefaultValues() []string {
	return o.defaultVals
}

func (o *FlagDef) Validate(val string) *errors.Error {
	switch o.flagType {
	case BOOL_FLAG:
		_, err := strconv.ParseBool(val)
		if err != nil {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s can not be parsed as boolean", val))
		}
		if len(o.allowedVals) > 0 && slices.Contains(o.allowedVals, val) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, allowed values: %v", val, o.allowedVals))
		}
	case NUMBER_FLAG:
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
	case STRING_FLAG:
		if len(o.allowedVals) > 0 && slices.Contains(o.allowedVals, val) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, allowed values: %v", val, o.allowedVals))
		}
		if !o.strRegex.Match([]byte(val)) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, string should be matched by regex %q", val, o.strRegex.String()))
		}
	default:
		return errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("%s is not a valid flag type.", o.flagType))
	}
	return nil
}
