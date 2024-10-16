package jwt

import "errors"

type Error struct {
	InnerErr error
	Message  string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.InnerErr
}

func NewError(err error, message string) error {
	return &Error{
		InnerErr: err,
		Message:  message,
	}
}

var ErrInvalidToken = NewError(errors.New("invalid token"), "invalid token")
