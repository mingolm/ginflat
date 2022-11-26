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

func NewHandler(opts ...HandleOption) *Handler {
	o := &handleOptions{
		render: defaultRender,
	}
	for _, opt := range opts {
		opt(o)
	}
	return &Handler{
		opts: o,
	}
}

type Handler struct {
	opts         *handleOptions
	hostHandlers sync.Map
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

func (h *Handler) RegisterController(host string, controllers ...Controller) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(h.opts.middlewares...)

	h.hostHandlers.Store(host, &ginHandler{
		engine:      engine,
		controllers: controllers,
	})
}

func (h *Handler) InitControllers(ctx context.Context) (err error) {
	h.hostHandlers.Range(func(key, value interface{}) bool {
		gh := value.(*ginHandler)
		host := key.(string)
		for _, ctrl := range gh.controllers {
			if err = ctrl.Init(ctx); err != nil {
				err = fmt.Errorf("%s controller init failed: %w", host, err)
				return false
			}

			ctrl.InitRouter(NewRouter(&gh.engine.RouterGroup, h.opts.render))
		}
		return true
	})

	return err
}
