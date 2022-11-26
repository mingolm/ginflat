package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"net/http"
	"time"
)

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
				_ = ctx.Error(err)
				ctx.Abort()
			}
		}()

		ctx.Next()

		errs := ctx.Errors.ByType(gin.ErrorTypeAny)
		if len(errs) == 0 {
			return
		}

		err = errs[0].Err
		if err != nil {
			ctx.Render(http.StatusBadRequest, ginflat.Json(gin.H{
				"success": false,
				"data":    err.Error(),
				"now":     time.Now().Unix(),
			}))
			return
		}
	}
}
