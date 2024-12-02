package main

import (
	"net/http"
	"strings"
)

// RouteGroup представляет группу маршрутов с возможностью вложенности и middleware
type RouterGroup struct {
	mux         *http.ServeMux
	prefix      string
	parent      *RouterGroup
	middlewares []Middleware
	logger      *CustomLogger
}

// NewRouteGroup создает новую группу маршрутов
func NewRouterGroup(mux *http.ServeMux, prefix string, logger *CustomLogger) *RouterGroup {
	return &RouterGroup{
		mux:    mux,
		prefix: prefix,
	}
}

// Group создает вложенную группу с общим префиксом
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		mux:    rg.mux,
		prefix: strings.TrimRight(rg.prefix+prefix, "/"),
		parent: rg,
	}
}

// Use добавляет middleware к группе
func (rg *RouterGroup) Use(mw Middleware) {
	rg.middlewares = append(rg.middlewares, mw)
}

func (rg *RouterGroup) GET(prefix string, handler http.HandlerFunc) {
	rg.Handle(http.MethodGet, prefix, handler)
}

func (rg *RouterGroup) POST(prefix string, handler http.HandlerFunc) {
	rg.Handle(http.MethodPost, prefix, handler)
}

func (rg *RouterGroup) PUT(prefix string, handler http.HandlerFunc) {
	rg.Handle(http.MethodPut, prefix, handler)
}

func (rg *RouterGroup) DELETE(prefix string, handler http.HandlerFunc) {
	rg.Handle(http.MethodDelete, prefix, handler)
}

// collectMiddleware собирает middleware от текущей группы и всех родительских
func (rg *RouterGroup) CollectMiddlewares() []Middleware {
	var mws []Middleware
	current := rg
	for current != nil {
		mws = append(current.middlewares, mws...)
		current = current.parent
	}
	return mws
}

// Handle регистрирует обработчик с применением middleware всех уровней
func (rg *RouterGroup) Handle(method, prefix string, handler http.Handler) {
	fullPath := strings.TrimRight(rg.prefix+prefix, "/")

	rg.mux.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			rg.logger.Error("Method Not Allowed: %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		finalHandler := handler

		mws := rg.CollectMiddlewares()

		for i := len(mws) - 1; i >= 0; i-- {
			finalHandler = mws[i](finalHandler)
		}

		finalHandler.ServeHTTP(w, r)
	})
}

func registerHandlers(mux *http.ServeMux, h *handler, logger *CustomLogger) {
	mux.HandleFunc("/test1", h.TestHandler1)

	g1 := NewRouterGroup(mux, "/group1", logger)

	g1.GET("/test", h.TestGroup1Handler)
}
