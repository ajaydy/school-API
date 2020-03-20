package helpers

import "github.com/pkg/errors"

type (
	Error struct {
		Err        error
		StatusCode int
		Message    string
	}
)

func (e *Error) Error() string {
	return e.Message
}

func ErrorWrap(err error, prefix, suffix, message string, status int) *Error {
	logger.Err.Errorf("error : %v : %v : %v", prefix, suffix, err)
	return &Error{
		Err:        errors.Wrapf(err, "%s/%s", prefix, suffix),
		Message:    message,
		StatusCode: status,
	}
}

const (
	InternalServerError = "Internal Server Error"
	BadRequestMessage   = "Bad Request"
)
