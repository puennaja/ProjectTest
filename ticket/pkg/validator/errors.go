package validator

import (
	"bytes"
	"strings"
)

type fieldError struct {
	field   string
	message string
}

type FieldError interface {
	Error() string
	Field() string
}

type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	var fe *fieldError
	for i := 0; i < len(ve); i++ {
		fe = ve[i].(*fieldError)
		buff.WriteString(fe.Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

func (fe *fieldError) Error() string {
	return fe.message
}

func (fe *fieldError) Field() string {
	return fe.field
}

func NewFieldError(field, message string) FieldError {
	return &fieldError{
		field:   field,
		message: message,
	}
}
