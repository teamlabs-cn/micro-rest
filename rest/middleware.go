package rest

import "errors"

type MiddlewareFunc func(ctx *HttpContext, next func(*HttpContext))

type MiddlewarePipeline struct {
	middlewares []MiddlewareFunc
}

func (p *MiddlewarePipeline) AddMiddleware(m MiddlewareFunc) {
	if m == nil {
		panic(errors.New("Middlware argument can not be nil"))
	}

	p.middlewares = append(p.middlewares, m)
}

// GetHandlerFunc method
func (p *MiddlewarePipeline) GetHandlerFunc(handler HandlerFunc) HandlerFunc {
	next := handler
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		next = GetMiddlewareFunc(p.middlewares[i], next)
	}
	return next
}

func GetMiddlewareFunc(current MiddlewareFunc, next func(ctx *HttpContext)) func(ctx *HttpContext) {
	return func(ctx *HttpContext) {
		current(ctx, next)
	}
}

func Mixin(handler HandlerFunc, pipelines ...MiddlewarePipeline) HandlerFunc {
	next := handler
	for i := len(pipelines) - 1; i >= 0; i-- {
		next = pipelines[i].GetHandlerFunc(next)
	}
	return next
}
