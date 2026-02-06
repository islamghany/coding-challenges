package handler

import (
	"webserver/http/request"
	"webserver/http/response"
)

type Handler interface {
	ServeHTTP(response.Response, *request.Request)
}

type HandlerFunc func(response.Response, *request.Request)

func (f HandlerFunc) ServeHTTP(w response.Response, r *request.Request) {
	f(w, r)
}
