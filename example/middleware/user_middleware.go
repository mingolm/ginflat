package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"net/http"
)

func UserMiddleware() ginflat.Middleware {
	return func(ctx *gin.Context) {
		// test
		if id := ctx.Request.FormValue("id"); id == "" {
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("user is not found"))
		}

		ctx.Next()
		return
	}
}
