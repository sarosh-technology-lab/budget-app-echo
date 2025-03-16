package custom_app_errors

type NotFoundError struct {
	Message string
}

func (receriver NotFoundError) Error() string {
	return receriver.Message
}

func NewNotFoundError(message string) NotFoundError {
	if message == "" {
		message = "resource not found"
	}
	return NotFoundError{
		Message: message,
	}
}