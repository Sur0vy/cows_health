package storage

type ExistError struct {
	message string
}

type EmptyError struct {
	message string
}

func NewExistError(msg string) *ExistError {
	return &ExistError{
		message: msg,
	}
}

func (e *ExistError) Error() string {
	return e.message
}

func NewEmptyError(msg string) *EmptyError {
	return &EmptyError{
		message: msg,
	}
}

func (e *EmptyError) Error() string {
	return e.message
}
