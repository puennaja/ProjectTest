package errors

import (
	"net/http"
)

type ExternalErr struct {
	Message string `json:"message"`
}

type ResponseErr struct {
	httpStatus int
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	Errors     []ExternalErr `json:"errors,omitempty"`
}

func NewResponseErr(err error, optsFuncs ...func(*ResponseErr)) *ResponseErr {
	responseErr := &ResponseErr{
		httpStatus: http.StatusInternalServerError,
		Code:       DefaultCode,
		Message:    http.StatusText(http.StatusInternalServerError),
	}
	if err != nil {
		internalErr := IsInternalError(err)
		if internalErr != nil {
			responseErr.httpStatus = internalErr.status
			responseErr.Code = internalErr.code
			responseErr.Message = internalErr.message
			responseErr.Errors = append(responseErr.Errors, ExternalErr{Message: internalErr.Error()})
		} else {
			responseErr.Errors = append(responseErr.Errors, ExternalErr{Message: err.Error()})
		}
	}
	for _, optsFunc := range optsFuncs {
		optsFunc(responseErr)
	}
	return responseErr

}

func (r *ResponseErr) HTTPStatus() int {
	return r.httpStatus
}

func (r *ResponseErr) Error() string {
	return r.Message
}

func IsResponseErr(e error) *ResponseErr {
	if err, ok := e.(*ResponseErr); ok {
		return err
	}
	return nil
}

func OptionHttpStatus(httpStatus int) func(*ResponseErr) {
	return func(e *ResponseErr) {
		e.httpStatus = httpStatus
	}
}

func OptionCode(code int) func(*ResponseErr) {
	return func(e *ResponseErr) {
		e.Code = code
	}
}

func OptionErr(err error) func(*ResponseErr) {
	return func(e *ResponseErr) {
		internalErr := IsInternalError(err)
		if internalErr != nil {
			e.Errors = append(e.Errors, ExternalErr{Message: internalErr.Error()})
		}
		e.Errors = append(e.Errors, ExternalErr{Message: err.Error()})
	}
}

func OptionErrs(errs []error) func(*ResponseErr) {
	return func(e *ResponseErr) {
		for _, err := range errs {
			internalErr := IsInternalError(err)
			if internalErr != nil {
				e.Errors = append(e.Errors, ExternalErr{Message: internalErr.Error()})
			}
			e.Errors = append(e.Errors, ExternalErr{Message: err.Error()})
		}
	}
}
