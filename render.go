package ginflat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Render interface {
	Data(ctx *gin.Context, data interface{})
	Error(ctx *gin.Context, err error)
}

var defaultRender Render = JsonRender{}

type JsonRender struct {
}

func (r JsonRender) Data(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)

}
func (r JsonRender) Error(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
}
