package xerrors

import (
	"errors"
	"fmt"
	"runtime"
)

// xerrors global var
var (
	_caller bool
)

// Kind of errors
type Kind int16

// kind of errors
const (
	KindOK Kind = iota
	KindNotFound
	KindBadRequest
	KindUnauthorized
	KindInternalError
)

// Op is the operation when error happens
type Op string

// String  value of Op
func (op Op) String() string {
	return string(op)
}

// Fields of errors
type Fields map[string]interface{}

// Errors of xerrors
type Errors struct {
	Err      error
	InnerErr error
	kind     Kind
	op       Op
}

// New errors
func New(v ...interface{}) error {
	var (
		xerr = &Errors{}
		file string
		line int
	)

	// only cal _caller when xerrors _caller is true
	if _caller {
		_, file, line, _ = runtime.Caller(1)
	}

	for _, arg := range v {
		switch val := arg.(type) {
		case Op:
			xerr.op = val

		case string:
			if _caller {
				xerr.Err = fmt.Errorf("%s: %s: [file=%s, line=%d]", val, xerr.op, file, line)
				continue
			}
			xerr.Err = fmt.Errorf("%s: %s", val, xerr.op)

		case Kind:
			xerr.kind = val

		case *Errors:
			val.op = xerr.op
			// copy the errors
			xerr = val

			if _caller {
				xerr.Err = fmt.Errorf("error executing %s: [file=%s, line=%d] \n%w", xerr.op, file, line, val.Err)
				continue
			}
			xerr.Err = fmt.Errorf("error executing %s: %w", xerr.op, val.Err)

		case error:
			if _caller {
				xerr.Err = fmt.Errorf("%w: %s: [file=%s, line=%d]", val, xerr.op, file, line)
				continue
			}
			xerr.Err = fmt.Errorf("%w: %s", val, xerr.op)

		default:
			continue
		}
	}
	return xerr
}

// Error return string of error
func (e *Errors) Error() string {
	return e.Err.Error()
}

// Unwrap errors
func (e *Errors) Unwrap() error {
	return e.Err
}

// Kind of errors
func (e *Errors) Kind() Kind {
	return e.kind
}

// Is wrap the errors is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As wrap the error as
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap error
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// XUnwrap return errors with xerror package type
func XUnwrap(err error) *Errors {
	xerr, ok := err.(*Errors)
	if ok {
		return xerr
	}

	return nil
}

// SetCaller to print the stack-trace of the error
func SetCaller(c bool) {
	_caller = c
}
