package errors

type internalErr struct {
	status  int // http status code
	code    int
	message string
	err     error
}

func NewInternalErr(status, code int) *internalErr {
	return &internalErr{
		status:  status,
		code:    code,
		message: messageErr[code],
	}
}

func (i *internalErr) SetStatus(status int) *internalErr {
	i.status = status
	return i
}

func (i *internalErr) SetCode(code int) *internalErr {
	i.code = code
	return i
}

func (i *internalErr) SetMessage(message string) *internalErr {
	i.message = message
	return i
}

func (i *internalErr) SetError(err error) *internalErr {
	i.err = err
	return i
}

func (i *internalErr) Error() string {
	if i.err != nil {
		return i.message + " : " + i.err.Error()
	}
	return i.message
}

func (i *internalErr) GetCode() int {
	return i.code
}

func (i *internalErr) Unwrap() error {
	if e, ok := i.err.(*internalErr); ok {
		return e.Unwrap()
	}

	return i.err
}

func IsInternalError(e error) *internalErr {
	if i, ok := e.(*internalErr); ok {
		return i
	}
	return nil
}
