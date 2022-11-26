package ginflat

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

func Json(data interface{}) render.Render {
	return jrender{
		data: data,
	}
}

type Render interface {
	DataRender(ctx *gin.Context, data interface{})
	ErrRender(ctx *gin.Context, err error)
}

var (
	defaultRender   Render = jrender{}
	jsonContentType        = []string{"application/json; charset=utf-8"}
)

type jrender struct {
	data interface{}
}

type jsonResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (r jrender) DataRender(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &jsonResponse{
		Success: true,
		Data:    data,
	})

}
func (r jrender) ErrRender(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, &jsonResponse{
		Success: false,
		Data:    err.Error(),
	})
}

func (r jrender) Render(w http.ResponseWriter) (err error) {
	bs, err := json.Marshal(&r.data)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

func (r jrender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentType
	}
}
