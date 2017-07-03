package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// RouteCollection is the collection of route rule
type RouteCollection interface {
	AddRoute(method, route string, handler HandlerFunc)
}

// RouteHandler is the handler of route rule
type RouteHandler interface {
	GetHandlerFunc(routeCtx *RouteContext) (HandlerFunc, bool)
}

// RouteHandlerFunc is the function type of RouteHandler
type RouteHandlerFunc func(ctx *RouteContext) (HandlerFunc, bool)

// RouteContext is the route context for handling route
type RouteContext struct {
	RoutePath string
	Request   *http.Request
}

// RouteProvider is the route provider for routing
type RouteProvider struct {
	routeMap map[string][]RouteHandlerFunc
}

// GetHandlerFunc returns the function of request handler
func (p *RouteProvider) GetHandlerFunc(routeCtx *RouteContext) (HandlerFunc, bool) {
	r := routeCtx.Request
	if routes, ok := p.routeMap[r.Method]; ok {
		for _, routeHandler := range routes {
			if handler, ok := routeHandler(routeCtx); ok {
				return handler, true
			}
		}
	}
	return nil, false
}

// AddRoute adds a route rule with the specified method & route pattern
func (p *RouteProvider) AddRoute(method, route string, handler HandlerFunc) {
	routeHandler, err := newRouteHandlerFunc(route, handler)
	if err != nil {
		log.Fatalln(err)
	}
	p.routeMap[method] = append(p.routeMap[method], routeHandler)
}

var newRouteHandlerFunc = func(route string, handler HandlerFunc) (func(*RouteContext) (HandlerFunc, bool), error) {
	if len(route) == 0 {
		return nil, errors.New("route argument can not be nil or empty")
	}
	if handler == nil {
		return nil, errors.New("handler argument can not be nil")
	}
	matchRoute, err := getRouteMatchFunc(route)
	if err != nil {
		return nil, err
	}
	return func(routeCtx *RouteContext) (HandlerFunc, bool) {
		if routeData, ok := matchRoute(routeCtx.RoutePath); ok {
			return func(ctx *HttpContext) {
				ctx.RouteData = routeData
				handler(ctx)
			}, true
		}
		return nil, false
	}, nil
}

var getRouteMatchFunc = func(route string) (func(string) (map[string]interface{}, bool), error) {
	convMap := make(map[string]func(string) interface{})
	pattern := regexp.MustCompile(`\{([\w_]+(:[a-z]+)?)\}*`).ReplaceAllStringFunc(route, func(matched string) string {
		tokens := strings.Split(strings.Trim(matched, "{}"), ":")
		key := tokens[0]
		p := "[:word:]"
		if len(tokens) > 1 {
			if convt := getConverter(tokens[1]); convt != nil {
				convMap[key] = convt
			}
			p = getPattern(tokens[1])
		}
		return fmt.Sprintf("(?P<%s>[%s]+)", strings.Trim(key, "{}"), p)
	})
	pattern = "^/" + strings.TrimPrefix(pattern, "/") + "$"
	// fmt.Println(pattern)
	routeExp, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return func(routePath string) (map[string]interface{}, bool) {
		m := routeExp.FindStringSubmatch(routePath)
		if len := len(m); len != 0 {
			routeData := make(map[string]interface{})
			names := routeExp.SubexpNames()
			for i := 1; i < len; i++ {
				if conv, ok := convMap[names[i]]; ok {
					routeData[names[i]] = conv(m[i])
				} else {
					routeData[names[i]] = m[i]
				}
			}
			return routeData, true
		}
		return nil, false
	}, nil
}

func parseInt(input string) interface{} {
	if result, error := strconv.Atoi(input); error == nil {
		return result
	}
	return 0
}

func getConverter(cons string) func(string) interface{} {
	switch cons {
	case "int":
		return parseInt
	default:
		return nil
	}
}

func getPattern(cons string) string {
	switch cons {
	case "int":
		return "[:digit:]"
	default:
		return cons
	}
}
