package errors

import (
	"errors"
	"fmt"
)

// copy implement from errors package
// to be using this package instead of errors package

func New(text string) error {
	return errors.New(text)
}

func Newf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return nil
}
