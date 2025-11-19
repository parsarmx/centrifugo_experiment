package errors

type AppError struct {
	Msg string
	Err error
}

func NewAppError(message string, err error) *AppError {
	return &AppError{Msg: message, Err: err}
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}
