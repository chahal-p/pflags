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
	UNKNOWN     = FlagType("unknown")
	BOOL_FLAG   = FlagType("boolean")
	NUMBER_FLAG = FlagType("number")
	STRING_FLAG = FlagType("string")
)

func TypeFromString(t string) (FlagType, *errors.Error) {
	switch t {
	case "boolean", "bool":
		return BOOL_FLAG, nil
	case "number":
		return NUMBER_FLAG, nil
	case "string":
		return STRING_FLAG, nil
	default:
		return UNKNOWN, errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Unrecognized type: %s", t))
	}

}

type FlagDef struct {
	shortName   string
	longName    string
	desc        string
	flagType    FlagType
	required    bool
	defaultVals []string
	allowedVals []string
	numRange    []float64
	strRegex    *regexp.Regexp
}

func New(shortName string, longName string, flagType FlagType, opts ...Option) (*FlagDef, *errors.Error) {
	flag := &FlagDef{
		shortName: shortName,
		longName:  longName,
		flagType:  flagType,
	}

	for _, opt := range opts {
		if err := opt(flag); err != nil {
			return nil, err
		}
	}

	if len(flag.allowedVals) > 0 && (slices.Contains([]FlagType{BOOL_FLAG}, flag.flagType)) {
		return nil, errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("Allowed values can not be provided for type %s", flag.flagType))
	}

	if len(flag.numRange) > 0 && (slices.Contains([]FlagType{BOOL_FLAG, STRING_FLAG}, flag.flagType)) {
		return nil, errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("Number range can not be provided for type %s", flag.flagType))
	}

	if flag.strRegex != nil && (slices.Contains([]FlagType{BOOL_FLAG, NUMBER_FLAG}, flag.flagType)) {
		return nil, errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("String regex can not be provided for type %s", flag.flagType))
	}
	return flag, nil
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

func (o *FlagDef) Description() string {
	return o.desc
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
		if o.strRegex != nil && !o.strRegex.Match([]byte(val)) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, string should be matched by regex %q", val, o.strRegex.String()))
		}
	default:
		return errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("%s is not a valid flag type.", o.flagType))
	}
	return nil
}

type Option func(*FlagDef) *errors.Error

func Description(desc string) Option {
	return func(fc *FlagDef) *errors.Error {
		fc.desc = desc
		return nil
	}
}

func DefaultValues(vals ...string) Option {
	return func(fc *FlagDef) *errors.Error {
		fc.defaultVals = vals
		return nil
	}
}

func Required() Option {
	return func(fc *FlagDef) *errors.Error {
		fc.required = true
		return nil
	}
}

func AllowedValues(vals ...string) Option {
	return func(fc *FlagDef) *errors.Error {
		fc.allowedVals = vals
		return nil
	}
}

func NumberRange(lowerBound, upperBound float64) Option {
	return func(fc *FlagDef) *errors.Error {
		fc.numRange = []float64{lowerBound, upperBound}
		return nil
	}
}

func StringRegex(regex string) Option {
	return func(fc *FlagDef) *errors.Error {
		complied, err := regexp.Compile(regex)
		if err != nil {
			return errors.NewError(errors.INVALID_USAGE, err.Error())
		}
		fc.strRegex = complied
		return nil
	}
}
