package httperrors

import (
	"fmt"
	"io"
)

func WithErrorType(err error, errorType string) error {
	return &withErrorType{
		error:     err,
		errorType: errorType,
	}
}

type withErrorType struct {
	error
	errorType string
}

func (w *withErrorType) ErrorType() string {
	return w.errorType
}

func (w *withErrorType) Cause() error { return w.error }

func (w *withErrorType) Unwrap() error { return w.error }

func (w *withErrorType) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			return
		}
		_, _ = fmt.Fprintf(s, "%v", w.Cause())
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}
