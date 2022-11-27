package httperrors

import (
	"github.com/pkg/errors"
	"net/http"
)

type Response struct {
	StatusCode int          `json:"status_code"`
	Data       ResponseData `json:"data"`
}

type ResponseData struct {
	Message       string            `json:"message"`
	ErrorType     string            `json:"error_type,omitempty"`
	ErrorMetadata map[string]string `json:"error_metadata,omitempty"`
}

func ToResponse(err error) *Response {
	resp := &Response{
		StatusCode: http.StatusInternalServerError,
		Data: ResponseData{
			Message: err.Error(),
		},
	}

	var codeE interface{ Code() int }
	if errors.As(err, &codeE) {
		resp.StatusCode = codeE.Code()
	}

	var ErrorTypeE interface{ ErrorType() string }
	if errors.As(err, &ErrorTypeE) {
		resp.Data.ErrorType = ErrorTypeE.ErrorType()
	}

	return resp
}
