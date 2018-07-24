package errors

// ErrInvalidVLAN is the sturcture
type ErrInvalidVLAN struct {
	message string
}

// NewErrInvalidVLAN will return error if new an invalid VLAN
func NewErrInvalidVLAN(message string) *ErrInvalidVLAN {
	return &ErrInvalidVLAN{
		message: message,
	}
}

// Error will reture the error message
func (e *ErrInvalidVLAN) Error() string {
	return e.message
}
