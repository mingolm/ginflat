package ginflat

import (
	"fmt"
	"github.com/mingolm/ginflat/internal"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Router interface {
	// Group 路由前缀
	Group(path string) Router
	// Use 使用中间键
	Use(handlers ...Middleware) Router

	Get(path string, handler interface{})
	Post(path string, handler interface{})
	Put(path string, handler interface{})
	Delete(path string, handler interface{})
	Options(path string, handler interface{})
	Connect(path string, handler interface{})
}

func NewRouter(rg *gin.RouterGroup, render Render) Router {
	return &router{
		render:      render,
		ginRouter:   rg,
		middlewares: MiddlewareChain{},
	}
}

type router struct {
	render           Render
	ginRouter        *gin.RouterGroup
	prefixRouterPath string
	middlewares      MiddlewareChain
}

func (r *router) Group(path string) Router {
	r.prefixRouterPath = path
	return r
}

func (r *router) Use(handlers ...Middleware) Router {
	r.middlewares = append(r.middlewares, handlers...)
	return r
}

func (r *router) Get(path string, handler interface{}) {
	r.handle(http.MethodGet, path, handler)
}

func (r *router) Post(path string, handler interface{}) {
	r.handle(http.MethodPost, path, handler)
}

func (r *router) Put(path string, handler interface{}) {
	r.handle(http.MethodPut, path, handler)
}

func (r *router) Delete(path string, handler interface{}) {
	r.handle(http.MethodDelete, path, handler)
}

func (r *router) Options(path string, handler interface{}) {
	r.handle(http.MethodOptions, path, handler)
}

func (r *router) Connect(path string, handler interface{}) {
	r.handle(http.MethodConnect, path, handler)
}

func (r *router) handle(method, path string, handler interface{}) {
	path = r.prefixRouterPath + path

	switch fn := handler.(type) {
	case func(ctx *gin.Context):
		r.ginRouter.Handle(method, path, fn)
		return
	case gin.HandlerFunc:
		r.ginRouter.Handle(method, path, fn)
		return
	}

	v := reflect.ValueOf(handler)
	t := reflect.ValueOf(handler).Type()
	if t.Kind() != reflect.Func {
		panic("handler should be a function")
	}

	switch {
	case t.NumIn() != 2 || t.NumOut() != 2: // 入/出参必须2个
		fallthrough
	case t.In(0) != reflect.TypeOf(&gin.Context{}): // 第1个必须为 *gin.Context
		fallthrough
	case t.In(1).Kind() != reflect.Pointer || t.In(1).Elem().Kind() != reflect.Struct: // 第2个必须为 *Struct
		fallthrough
	case t.Out(0).Kind() != reflect.Ptr || t.In(1).Elem().Kind() != reflect.Struct: // 第1个必须为 *Struct
		fallthrough
	case t.Out(1) != reflect.TypeOf((*error)(nil)).Elem(): // 第2个必须为 error
		panic(invalidHandlerError(v))
	}

	reqTyp := t.In(1)
	reqPool := &sync.Pool{New: func() any {
		rv := reflect.New(reqTyp.Elem())
		return &rv
	}}
	reqM := reflect.New(reqTyp.Elem()).Elem()

	r.ginRouter.Use(r.middlewares...)
	r.ginRouter.Handle(method, path, func(ctx *gin.Context) {
		reqP := reqPool.Get().(*reflect.Value)
		reqVal := *reqP
		reqVal.Elem().Set(reqM)
		defer reqPool.Put(reqP)

		if err := internal.BindingRequest(ctx, &reqVal); err != nil {
			r.render.ErrRender(ctx, err)
			return
		}
		callResp := v.Call([]reflect.Value{reflect.ValueOf(ctx), reqVal})
		if callResp[1].IsNil() {
			r.render.DataRender(ctx, callResp[0].Interface())
			return
		}
		if !callResp[0].IsNil() {
			r.render.DataRender(ctx, callResp[0].Interface())
		}
		r.render.ErrRender(ctx, callResp[1].Interface().(error))
	})
}

func invalidHandlerError(handler reflect.Value) error {
	handlerName := runtime.FuncForPC(handler.Pointer()).Name()

	_, structName, funcName := parseFuncName(handlerName)

	reqName := firstCharToUppercase(funcName + "Request")
	respName := firstCharToUppercase(funcName + "Response")

	return fmt.Errorf("%s: invalid handler defined, should be\n\tfunc (c *%s) %s(ctx *gin.Context, req *%s)(resp *%s, err error)",
		handlerName,
		structName, funcName, reqName, respName,
	)
}

func parseFuncName(s string) (pkg string, structName string, funcName string) {
	l := strings.LastIndex(s, ".")
	f := strings.LastIndex(s[:l], ".")

	pkg = s[:f]

	structName = strings.Trim(s[f+1:l], "(*)")
	funcName = strings.TrimSuffix(s[l+1:], "-fm")

	return
}

func firstCharToUppercase(str string) string {
	if str[0] >= 'a' && str[0] <= 'z' {
		bs := []byte(str)
		bs[0] -= 32
		return string(bs)
	}
	return str
}
