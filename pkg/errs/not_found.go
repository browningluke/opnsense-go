package errs

type NotFoundError struct {
	err *Error
}

func (e *NotFoundError) Error() string {
	return e.err.Error()
}

func (e *NotFoundError) StatusCode() int {
	return e.err.StatusCode
}

func (e *NotFoundError) ErrorMessage() string {
	return e.err.ErrorMessage
}

func NewNotFoundError() *NotFoundError {
	return &NotFoundError{
		err: &Error{
			Type:         ErrorTypeNotFound,
			StatusCode:   404,
			ErrorMessage: "Resource not found",
		},
	}
}
