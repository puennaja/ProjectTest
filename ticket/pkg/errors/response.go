package errors

import (
	"net/http"
)

type ExternalErr struct {
	Field   string `json:"field,omitempty"`
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
			if internalErr.err != nil {
				responseErr.Errors = append(responseErr.Errors, ExternalErr{Message: internalErr.err.Error()})
			}
		} else {
			responseErr.Message = err.Error()
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
			return
		}

		validationErrs := IsValidationErrors(err)
		if validationErrs != nil {
			for _, v := range validationErrs {
				e.Errors = append(e.Errors, ExternalErr{Field: v.Field(), Message: v.Error()})
			}
			return
		}

		e.Errors = append(e.Errors, ExternalErr{Message: err.Error()})
	}
}
