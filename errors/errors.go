package errors

import (
	"fmt"
)

type Code struct {
	code int
	name string
}

func (o *Code) Code() int {
	return o.code
}

func NewCode(code int, name string) Code {
	return Code{
		code: code,
		name: name,
	}
}

var (
	INVALID_USAGE = NewCode(2, "INVALID_USAGE")
	INVAID_VALUE  = NewCode(4, "INVAID_VALUE")
)

type Error struct {
	code Code
	msg  string
}

func NewError(code Code, msg string) *Error {
	return &Error{
		code: code,
		msg:  msg,
	}
}

func (o *Error) Error() string {
	return fmt.Sprintf("Error %s: %s", o.code.name, o.msg)
}

func (o *Error) Code() *Code {
	return &o.code
}
