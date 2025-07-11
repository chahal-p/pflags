package errors

import (
	"fmt"
)

type Code struct {
	code int
	name string
}

func (o Code) Code() int {
	return o.code
}

func NewCode(code int, name string) Code {
	return Code{
		code: code,
		name: name,
	}
}

var (
	ERROR                = NewCode(1, "ERROR")
	INVALID_USAGE        = NewCode(2, "INVALID_USAGE")
	INVAID_VALUE         = NewCode(4, "INVAID_VALUE")
	NOT_FOUND            = NewCode(40, "NOT_FOUND")
	INTERNAL_ERROR       = NewCode(99, "INTERNAL_ERROR")
	USAGE_HELP_REQUESTED = NewCode(100, "USAGE_HELP_REQUESTED")
)

type Error struct {
	code Code
	msg  string
}

func NewError(code Code, msgFormat string, a ...any) *Error {
	msg := msgFormat
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}

	return &Error{
		code: code,
		msg:  msg,
	}
}

func (o *Error) Error() string {
	return fmt.Sprintf("%s: %s", o.code.name, o.msg)
}

func (o *Error) Code() Code {
	return o.code
}
