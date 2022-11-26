package ginflat

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var _ http.Handler = &Handler{}

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct {
	hostHandlers sync.Map
	hostRenders  sync.Map
	middlewares  MiddlewareChain
}

type ginHandler struct {
	engine      *gin.Engine
	controllers []Controller
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	i := strings.LastIndex(req.Host, ":")
	host := req.Host
	if i != -1 {
		host = host[:i]
	}

	if req.URL.Path == "/healthz" {
		_, _ = rw.Write([]byte(`{"status":"ok"}`))
	} else if gh, ok := h.hostHandlers.Load(host); ok {
		gh.(*ginHandler).engine.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

func (h *Handler) Use(handles ...Middleware) {
	h.middlewares = append(h.middlewares, handles...)
}

func (h *Handler) RegisterRender(host string, render Render) {
	h.hostRenders.Store(host, render)
}

func (h *Handler) RegisterController(host string, controllers ...Controller) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(h.middlewares...)

	h.hostHandlers.Store(host, &ginHandler{
		engine:      engine,
		controllers: controllers,
	})
}

func (h *Handler) InitControllers(ctx context.Context) (err error) {
	h.hostHandlers.Range(func(key, value interface{}) bool {
		gh := value.(*ginHandler)
		host := key.(string)

		// render
		rd := defaultRender
		if rdi, ok := h.hostRenders.Load(host); ok {
			rd = rdi.(Render)
		}

		for _, ctrl := range gh.controllers {
			if err = ctrl.Init(ctx); err != nil {
				err = fmt.Errorf("%s controller init failed: %w", host, err)
				return false
			}

			ctrl.InitRouter(NewRouter(&gh.engine.RouterGroup, rd))
		}
		return true
	})

	return err
}
