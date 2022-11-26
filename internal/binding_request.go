package internal

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

func BindingRequest(ctx *gin.Context, req *reflect.Value) error {
	return ctx.Bind(req.Interface())
}
