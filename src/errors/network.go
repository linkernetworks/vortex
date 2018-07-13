package errors

type ErrInvalidVLAN struct {
	message string
}

func NewErrInvalidVLAN(message string) *ErrInvalidVLAN {
	return &ErrInvalidVLAN{
		message: message,
	}
}

func (e *ErrInvalidVLAN) Error() string {
	return e.message
}
