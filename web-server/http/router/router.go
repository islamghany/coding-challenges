package router

import (
	"webserver/http/handler"
	"webserver/http/request"
	"webserver/http/response"
)

type Router struct {
	tree *RouteTree
}

func NewRouter() *Router {
	return &Router{
		tree: NewRouteTree(),
	}
}

func (r *Router) addRoute(method, path string, handler handler.Handler) {
	r.tree.AddRoute(method, path, handler)
}

func (r *Router) Get(path string, handler handler.Handler) {
	r.addRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler handler.Handler) {
	r.addRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler handler.Handler) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) Delete(path string, handler handler.Handler) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) ServeHTTP(w response.Response, req *request.Request) {
	r.tree.ServeHTTP(w, req)
}

func (r *Router) SetNotFoundHandler(handler handler.Handler) {
	r.tree.SetNotFoundHandler(handler)
}

func (r *Router) SetMethodNotAllowedHandler(handler handler.Handler) {
	r.tree.SetMethodNotAllowedHandler(handler)
}
