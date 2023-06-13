package utils

type HttpError interface {
	Message() string
	Status() int
}

type myHttpError struct {
	err    error
	status int
}

func (e *myHttpError) Message() string {
	return e.err.Error()
}

func (e *myHttpError) Status() int {
	return e.status
}

func NewHttpError(err error, status int) HttpError {
	return &myHttpError{err, status}
}
