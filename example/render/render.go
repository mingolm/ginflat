package render

import (
	"github.com/gin-gonic/gin"
	"github.com/mingolm/ginflat"
	"net/http"
	"time"
)

var instance ginflat.Render = render{}

func Render() ginflat.Render {
	return instance
}

type render struct {
}

func (r render) DataRender(ctx *gin.Context, data interface{}) {
	h := gin.H{
		"success": true,
	}

	switch v := data.(type) {
	case interface{ WithContext(*gin.Context) }:
		v.WithContext(ctx)
		h["data"] = data
	default:
		h["data"] = data
	}

	if ctx.Writer.Header().Get("Cache-Control") == "" {
		ctx.Header("Cache-Control", "no-cache")
	}

	// no-cache返回now
	h["now"] = time.Now().Unix()

	ctx.JSON(http.StatusOK, h)
}

func (r render) ErrRender(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	ctx.Abort()
}
