package rest

import (
	"errors"
	"net/http"
	"strings"
)

type HttpResource struct {
	RouteProvider
}

func (res *HttpResource) Get(route string, handler HandlerFunc) {
	res.AddRoute("GET", route, handler)
}

func (res *HttpResource) Post(route string, handler HandlerFunc) {
	res.AddRoute("POST", route, handler)
}

func (res *HttpResource) Put(route string, handler HandlerFunc) {
	res.AddRoute("PUT", route, handler)
}

func (res *HttpResource) Patch(route string, handler HandlerFunc) {
	res.AddRoute("PATCH", route, handler)
}

func (res *HttpResource) Delete(route string, handler HandlerFunc) {
	res.AddRoute("DELETE", route, handler)
}

type HttpContext struct {
	Response  http.ResponseWriter
	Request   *http.Request
	RouteData map[string]interface{}
}

type HandleFunc func(pattern string, handler func(http.ResponseWriter, *http.Request))

type HandlerFunc func(*HttpContext)

func Resource(path string, factory func() HttpResource) {
	globalApp.Resource(path, factory)
}

func Use(middelware MiddlewareFunc) {
	globalApp.Use(middelware)
}

func NewResource() HttpResource {
	return HttpResource{RouteProvider{routeMap: make(map[string][]RouteHandlerFunc)}}
}

type Application struct {
	RootPath           string
	HandleFuncDelegate HandleFunc

	MiddlewarePipeline
}

func (app *Application) Resource(path string, factory func() HttpResource) {
	if !strings.HasPrefix(path, "/") {
		panic(errors.New("Path must be start with '/'"))
	}

	if path == "/" {
		panic(errors.New("Path can not be '/'"))
	}

	resource := factory()
	app.HandleFunc(path, resource)
}

func (app *Application) Use(m MiddlewareFunc) {
	app.AddMiddleware(m)
}

func (app *Application) HandleFunc(path string, resource HttpResource) {
	resPath := app.RootPath + strings.TrimSuffix(path, "/")
	handler := func(ctx *HttpContext) {
		routePath := strings.TrimPrefix(ctx.Request.URL.Path, resPath)
		if routeHandler, ok := resource.GetHandlerFunc(&RouteContext{routePath, ctx.Request}); ok {
			routeHandler(ctx)
		} else {
			http.NotFound(ctx.Response, ctx.Request)
		}
	}

	pipelines := []MiddlewarePipeline{globalApp.MiddlewarePipeline}
	if app != globalApp {
		pipelines = append(pipelines, app.MiddlewarePipeline)
	}
	handler = Mixin(handler, pipelines...)

	app.HandleFuncDelegate(resPath+"/", func(w http.ResponseWriter, r *http.Request) {
		handler(&HttpContext{Request: r, Response: w})
	})
}

func NewApp(path string) Application {
	return NewAppWithHandleFunc(path, http.HandleFunc)
}

func NewAppWithHandleFunc(path string, handleFunc HandleFunc) Application {
	if !strings.HasPrefix(path, "/") {
		panic(errors.New("Path must be start with '/'"))
	}

	if path == "/" {
		panic(errors.New("Path can not be '/'"))
	}

	path = strings.TrimSuffix(path, "/")
	return Application{
		RootPath:           path,
		HandleFuncDelegate: handleFunc,
	}
}

var globalApp = &Application{HandleFuncDelegate: http.HandleFunc}
