package main

// CustomError ...
type CustomError struct {
	message string
}

// NewCustomError creates CustomError struct
func NewCustomError(message string) *CustomError {
	return &CustomError{
		message: message,
	}
}

func (e *CustomError) Error() string {
	return e.message
}

func (e *CustomError) Is(target error) bool {
	_, ok := target.(*CustomError)
	return ok
}

// IsCustomError evaluates if the given error is CustomError
func IsCustomError(err error) bool {
	if _, ok := err.(*CustomError); ok {
		return true
	}

	return false
}
