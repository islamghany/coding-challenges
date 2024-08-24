package router

import (
	"context"
	"strings"
	"sync"
	"webserver/http/handler"
	"webserver/http/request"
	"webserver/http/response"
)

// route tree structure
// a tree structure is used to store the routes hierarchically
// each node in the tree represents a segment of the path
// the root node represents the first segment of the path
// the leaves of the tree contain the handler functions
type RouteTree struct {
	root             *Node
	notFound         handler.Handler
	methodNotAllowed handler.Handler
	pool             sync.Pool
	mutex            sync.RWMutex
}

type Node struct {
	Path       string
	handlers   map[string]handler.Handler
	childrens  map[string]*Node
	parameter  string
	wildcard   *Node
	middleware []func(handler.Handler) handler.Handler
}

// NewRouter creates a new router.
func NewRouteTree() *RouteTree {
	return &RouteTree{
		root: &Node{
			Path:      "/",
			childrens: make(map[string]*Node),
			handlers:  make(map[string]handler.Handler),
		},
		pool: sync.Pool{
			New: func() interface{} {
				return make(map[string]string)
			},
		},
	}
}

// routing algorithm
// the routing algorithm is a recursive function that traverses the tree
// starting from the root node
// for each segment of the path:
// a) Match against the static children of the current node
// b) If no match is found, check for dynamic segments (parameters)
// c) check for wildcard segments
// continue until the end of the path is reached

// AddRoute adds a route to the route tree.
func (rt *RouteTree) AddRoute(method, path string, handlerFunc handler.Handler, middleware ...func(handler.Handler) handler.Handler) {
	// lock the route tree to prevent concurrent writes
	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	segments := splitPath(path)
	// begin at the root node
	current := rt.root
	for i, segment := range segments {
		// check if the segment is a parameter (dynamic segment)
		if len(segment) > 0 && segment[0] == ':' {
			//  If the segment starts with a :, it is considered a parameter. If the current node's parameter field is empty,
			//  it is set to the segment without the :. If the segment is the last one, the handler and middleware are assigned
			//  to the current node. Otherwise, a new node is created if it does not exist, and the function moves to the child
			//  node corresponding to the parameter.
			if current.parameter == "" {
				current.parameter = segment[1:]
			}
			// leaf node
			if i == len(segments)-1 {
				current.handlers[method] = handlerFunc
				current.middleware = middleware
			} else {
				// for non-leaf nodes, create a new node if it does not exist
				if current.childrens == nil {
					current.childrens = make(map[string]*Node)
				}
				current = current.childrens[":"+current.parameter]

			}
			// check if the segment is a wildcard
		} else if segment == "*" {
			// If the segment is *, it is considered a wildcard. If the current node's wildcard field is nil,
			// a new wildcard node is created. The handler and middleware are assigned to the wildcard node, and the loop breaks
			// as wildcards match any remaining path segments.
			if current.wildcard == nil {
				current.wildcard = &Node{
					Path:     "*",
					handlers: make(map[string]handler.Handler),
				}
			}
			current.wildcard.handlers[method] = handlerFunc
			current.wildcard.middleware = middleware
			break

		} else { // check if the segment is a static segment
			// For static segments, the function checks if the current node's childrens map is nil and initializes it if necessary.
			// It then checks if a child node for the segment exists. If not, a new node is created and added to the childrens map.
			// The function then moves to the child node.
			if current.childrens == nil {
				current.childrens = make(map[string]*Node)
			}
			child, ok := current.childrens[segment]
			if !ok {
				child = &Node{
					Path:      segment,
					handlers:  make(map[string]handler.Handler),
					childrens: make(map[string]*Node),
				}
				current.childrens[segment] = child
			}
			current = child
		}

	}
	if current.handlers == nil {
		current.handlers = make(map[string]handler.Handler)
	}
	current.handlers[method] = handlerFunc
	current.middleware = middleware

}

// FindRoute finds a route in the route tree.
func (rt *RouteTree) FindRoute(method, path string) (handler.Handler, map[string]string) {
	// lock the route tree to prevent concurrent reads
	rt.mutex.RLock()
	defer rt.mutex.RUnlock()

	segments := splitPath(path)
	current := rt.root
	// parameters are stored in a map
	params := rt.pool.Get().(map[string]string)
	for k := range params {
		delete(params, k)
	}

	for idx, segment := range segments {

		if child, ok := current.childrens[segment]; ok {
			current = child
		} else if current.parameter != "" {
			params[current.parameter] = segment

			if child, ok := current.childrens[":"+current.parameter]; ok {
				current = child
			}
			// if the child node does not exist and it's not the last segment, return not found
			if !ok && idx != len(segments)-1 {
				return rt.notFound, params
			}
		} else if current.wildcard != nil {
			return rt.applyMiddleware(current.wildcard.handlers[method], current.wildcard.middleware), params
		} else {
			return rt.notFound, params
		}
	}
	if handler, ok := current.handlers[method]; ok {
		return rt.applyMiddleware(handler, current.middleware), params
	}

	if len(current.handlers) > 0 {
		return rt.methodNotAllowed, params
	}
	return rt.notFound, params
}

// apply middleware
func (rt *RouteTree) applyMiddleware(handler handler.Handler, middleware []func(handler.Handler) handler.Handler) handler.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

// SetNotFoundHandler sets the not found handler.
func (rt *RouteTree) SetNotFoundHandler(handler handler.Handler) {
	rt.notFound = handler
}

// SetMethodNotAllowedHandler sets the method not allowed handler.
func (rt *RouteTree) SetMethodNotAllowedHandler(handler handler.Handler) {
	rt.methodNotAllowed = handler
}

// ServeHTTP serves an HTTP request.
func (rt *RouteTree) ServeHTTP(res response.Response, req *request.Request) {
	handler, params := rt.FindRoute(req.Method, req.Path)
	if handler != nil {
		ctx := context.WithValue(req.Context, "params", params)
		handler.ServeHTTP(res, req.WithContext(ctx))
	} else {
		res.Write([]byte("404 Not Found"))
	}
	rt.pool.Put(params) // release the params map
}

func splitPath(path string) []string {
	return strings.Split(path, "/")[1:]
}
