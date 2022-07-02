package errors

type ExistError struct {
	message string
}

type EmptyError struct {
	message string
}

func NewExistError() *ExistError {
	return &ExistError{
		message: "entry already exist",
	}
}

func (e *ExistError) Error() string {
	return e.message
}

func NewEmptyError() *EmptyError {
	return &EmptyError{
		message: "entry is missing",
	}
}

func (e *EmptyError) Error() string {
	return e.message
}
