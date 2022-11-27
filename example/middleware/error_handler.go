package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"github.com/mingolm/ginflat/httperrors"
)

type ErrResponse struct {
	Msg       string
	ErrorType string
}

func ErrHandler() ginflat.Middleware {
	return func(ctx *gin.Context) {
		var err error

		defer func() {
			if x := recover(); x != nil {
				if e, ok := x.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("panic: %+v", x)
				}
			}

			if err != nil {
				handle(ctx, err)
			}
		}()

		ctx.Next()

		errs := ctx.Errors.ByType(gin.ErrorTypeAny)
		if len(errs) == 0 {
			return
		}

		err = errs[0].Err
	}
}

func handle(ctx *gin.Context, err error) {
	xe := httperrors.ToResponse(err)
	ctx.Render(xe.StatusCode, ginflat.Json(gin.H{
		"success": false,
		"err_msg": xe.Data,
	}))
}
