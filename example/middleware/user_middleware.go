package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
)

func UserMiddleware() ginflat.Middleware {
	return func(ctx *gin.Context) {
		fmt.Println("user middleware test")
		ctx.Next()
	}
}
