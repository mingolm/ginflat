package httperrors

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrNotFound        = newWithCode(http.StatusNotFound)
	ErrInvalidArgument = newWithCode(http.StatusBadRequest)
)

type ChainError struct {
	stack bool
	cause error
}

func (e ChainError) Error() string {
	return e.cause.Error()
}

func (e ChainError) Cause() error {
	return e.cause
}

func (e ChainError) Unwrap() error {
	return e.cause
}

func (e ChainError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", e.cause)
			return
		}
		_, _ = fmt.Fprintf(s, "%v", e.cause)
	case 's':
		_, _ = fmt.Fprintf(s, "%s", e.cause)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.cause)
	}
}

func (e ChainError) Msg(msg string) *ChainError {
	if !e.stack {
		e.stack = true
		e.cause = errors.WithStack(e.cause)
	}

	e.cause = errors.WithMessage(e.cause, msg)
	return &e
}

func (e ChainError) ErrorType(typ string) *ChainError {
	if !e.stack {
		e.stack = true
		e.cause = errors.WithStack(e.cause)
	}

	e.cause = WithErrorType(e.cause, typ)
	return &e
}

func NewWithCode(code int, msg string) *ChainError {
	return &ChainError{
		stack: true,
		cause: errors.Wrap(newWithCode(code), msg),
	}
}

type codeError struct {
	statusCode int
}

func (e codeError) Code() int {
	return e.statusCode
}

func (e codeError) Error() string {
	return http.StatusText(e.statusCode)
}

func newWithCode(code int) *ChainError {
	return &ChainError{
		stack: false,
		cause: codeError{statusCode: code},
	}
}
