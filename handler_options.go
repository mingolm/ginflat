package ginflat

type handleOptions struct {
	render      Render
	middlewares MiddlewareChain
}

type HandleOption func(*handleOptions)

func WithRender(render Render) HandleOption {
	return func(o *handleOptions) {
		o.render = render
	}
}

func WithMiddlewares(middlewares ...Middleware) HandleOption {
	return func(o *handleOptions) {
		o.middlewares = middlewares
	}
}
