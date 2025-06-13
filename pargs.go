package pargs

import (
	"regexp"
)

type ArgType string

const (
	BOOL   = ArgType("boolean")
	NUMBER = ArgType("number")
	STRING = ArgType("string")
)

type Arg struct {
	name     string
	expected []string
	regex    *regexp.Regexp
	argType  ArgType
}

func Parse(args []string) error {
	return nil
}
