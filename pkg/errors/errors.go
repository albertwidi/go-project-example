// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// a heavy modified errors package from upspin project
// LICENSE can be found in upspin project LICENSE file.

package errors

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

// Fields of error
type Fields map[string]interface{}

// Op is operation of error
type Op string

// ToArrayInterface fields
func (f Fields) ToArrayInterface() []interface{} {
	fieldsLength := len(f)
	if fieldsLength == 0 {
		return nil
	}
	// because fields is key value, we need to create an array with len * 2
	intf := make([]interface{}, fieldsLength*2)
	counter := 0
	for key, val := range f {
		intf[counter] = key
		intf[counter+1] = val
		counter += 2
	}
	return intf
}

// Error struct
type Error struct {
	Err error
	// Codes used for Errs to identify known errors in the application
	// If the error is expected by Errs object, the errors will be shown as listed in Codes
	Code Codes
	// Traces used to add function traces to errors, this is different from context
	// While context is used to add more information about the error, traces is used
	// for easier function tracing purposes without hurting heap too much
	Traces []string
	// Fields is a fields context similar to logrus.Fields
	// Can be used for adding more context to the errors
	Fields Fields
	// Op is operaions of error
	Op Op
	// OpTraces is a trace of operations
	OpTraces []Op
}

var _ error = (*Error)(nil)

// E for creating new errors
func E(args ...interface{}) error {
	if len(args) == 0 {
		return nil
	}

	err := &Error{}
	for _, arg := range args {
		switch arg.(type) {
		case string:
			err.Err = Errorf(arg.(string))
		case *Error:
			// copy and put the errors back
			errcpy := *arg.(*Error)
			err = &errcpy
		// error should be placed below *Errs
		// implementation of Error() string will detect *Errs as error
		case error:
			err.Err = arg.(error)
		case Codes:
			err.Code = arg.(Codes)
			errString, _ := err.Code.ErrorAndCode()
			err.Err = errors.New(errString)
		// Fields cannot be appended
		// new fields will always replace the old fields
		case Fields:
			err.Fields = arg.(Fields)
		case Op:
			// check wether current error operations is empty, but always append the traces
			if err.Op == "" {
				err.Op = arg.(Op)
			}
			err.OpTraces = append(err.OpTraces, arg.(Op))
		default:
			// the default error is unknown
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.Errs: bad call from %s:%d: %v", file, line, args)
			return Errorf("unknown type %T, value %v", arg, arg)
		}
	}
	return err
}

// WithCodes give a safer passing of codes to errors as compiler/linter will check the interface{} implementation
func WithCodes(codes Codes) error {
	return E(codes)
}

// Errorf function
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// GetFields return available fields in errors
func (e *Error) GetFields() Fields {
	return e.Fields
}

// Match errorr
// Match will match two strings error through a fuzzy matching
// Need some improvement in fuzzy matching, not all cases is covered
func Match(errs1, errs2 error) bool {
	if errs1 == nil && errs2 == nil {
		return true
	}

	if errs1 != nil {
		err1, ok := errs1.(*Error)
		if ok {
			errs1 = err1.Err
		}
	} else {
		errs1 = errors.New("nil")
	}

	if errs2 != nil {
		err2, ok := errs2.(*Error)
		if ok {
			errs2 = err2.Err
		}
	} else {
		errs2 = errors.New("nil")
	}

	if errs1.Error() != errs2.Error() {
		return false
	}
	return true
}

// Codes is interface to define error custom code.
// It have two function called ErrorAndCode which return string of error and httpcode desired from the error
// Err will return the error of code itself, so error can be implemented directly in Codes
type Codes interface {
	ErrorAndCode() (string, int)
	Err() error
}
