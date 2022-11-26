package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"net/http"
)

func ErrHandler() ginflat.Middleware {
	return func(ctx *gin.Context) {
		ctx.Next()

		errs := ctx.Errors.ByType(gin.ErrorTypeAny)
		if len(errs) == 0 {
			return
		}

		err := errs[0].Err
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
}
