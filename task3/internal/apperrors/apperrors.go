package apperrors

type ServerError struct {
	Code    int    // HTTP-статус код
	Message string // Сообщение для клиента
	Err     error  // Оригинальная ошибка (для логирования)
}

// Полная ошибка
func (e *ServerError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return e.Message
	}
	return e.Message + ": " + e.Err.Error()
}

func (e *ServerError) ClientError() string {
	return e.Message
}

func New(code int, message string, err error) *ServerError {
	return &ServerError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ModifyErr(se ServerError, err error) *ServerError {
	return &ServerError{
		Code:    se.Code,
		Message: se.Message,
		Err:     err,
	}
}

var (
	ErrInvalidInput  = New(400, "Invalid input data", nil)
	ErrInternal      = New(500, "Internal server error", nil)
	ErrAlreadyExists = New(503, "Event already exists", nil)
	ErrNotFound      = New(503, "Event not found", nil)
	ErrPastDate      = New(503, "Date cannot be in the past", nil)
)
