package node

import (
	"github.com/gin-gonic/gin"
)

type RMethod string

const (
	MethodAny     RMethod = "any"
	MethodGet     RMethod = "get"
	MethodPost    RMethod = "post"
	MethodPut     RMethod = "put"
	MethodPatch   RMethod = "patch"
	MethodDelete  RMethod = "delete"
	MethodHead    RMethod = "head"
	MethodOptions RMethod = "options"
)

type clientRoute struct {
	Method      RMethod
	Path        string
	Name        string
	Description string
	Handler     gin.HandlerFunc
}

func (c *client) AddRoute(method RMethod, path string, handler gin.HandlerFunc) {
	c.AddRouteEntry(&clientRoute{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (c *client) AddRouteEntry(route *clientRoute) {
	c.Handlers = append(c.Handlers, route)
}

func (c *client) SetupRouter() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	for _, route := range c.Handlers {
		switch route.Method {
		case MethodGet:
			r.GET(route.Path, route.Handler)
		case MethodPost:
			r.POST(route.Path, route.Handler)
		case MethodPut:
			r.PUT(route.Path, route.Handler)
		case MethodPatch:
			r.PATCH(route.Path, route.Handler)
		case MethodDelete:
			r.DELETE(route.Path, route.Handler)
		case MethodHead:
			r.HEAD(route.Path, route.Handler)
		case MethodOptions:
			r.OPTIONS(route.Path, route.Handler)
		case MethodAny:
			r.Any(route.Path, route.Handler)
		}
	}

	c.Router = r
}

func (c client) RunRouter(addr string) {
	c.Logger.Info().Msgf("Configured local client [%s]", addr)
	c.SetupRouter()

	c.Router.Run(addr)
}
