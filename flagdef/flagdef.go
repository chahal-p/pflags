package flagdef

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/chahal-p/pflags/errors"
)

type FlagType string

const (
	UNKNOWN     = FlagType("unknown")
	BOOL_FLAG   = FlagType("boolean")
	NUMBER_FLAG = FlagType("number")
	STRING_FLAG = FlagType("string")
)

func DefaultValueForType(f FlagType) string {
	switch f {
	case BOOL_FLAG:
		return "false"
	case NUMBER_FLAG:
		return "0"
	default:
		return ""
	}
}

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
	strRegex    *regexp.Regexp
}

func New(shortName string, longName string, flagType FlagType, opts ...Option) (*FlagDef, *errors.Error) {
	shortName = strings.Trim(shortName, " ")
	longName = strings.Trim(longName, " ")
	if shortName == "" && longName == "" {
		return nil, errors.NewError(errors.INVALID_USAGE, "At least one of short or long flag name is required.")
	}
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
		if val != "false" && val != "true" {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, boolean can only take true or false.", val))
		}
	case NUMBER_FLAG:
		_, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s can not be parsed as number", val))
		}
	case STRING_FLAG:
		if o.strRegex != nil && !o.strRegex.Match([]byte(val)) {
			return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, string should be matched by regex %q", val, o.strRegex.String()))
		}
	default:
		return errors.NewError(errors.INVALID_USAGE, fmt.Sprintf("%s is not a valid flag type.", o.flagType))
	}
	if len(o.allowedVals) > 0 && !slices.Contains(o.allowedVals, val) {
		return errors.NewError(errors.INVAID_VALUE, fmt.Sprintf("Invalid value: %s, allowed values: %v", val, o.allowedVals))
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

func Required(required bool) Option {
	return func(fc *FlagDef) *errors.Error {
		fc.required = required
		return nil
	}
}

func AllowedValues(vals ...string) Option {
	return func(fc *FlagDef) *errors.Error {
		fc.allowedVals = vals
		return nil
	}
}

func StringRegex(regex string) Option {
	return func(fc *FlagDef) *errors.Error {
		if regex == "" {
			return nil
		}
		complied, err := regexp.Compile(regex)
		if err != nil {
			return errors.NewError(errors.INVALID_USAGE, err.Error())
		}
		fc.strRegex = complied
		return nil
	}
}
