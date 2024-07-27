package main

import "net/http"

// NewServeMux builds a http.ServeMux that will route requests
// to the given EchoHandler.
// func NewServeMux(echo *EchoHandler) *http.ServeMux {
// 	mux := http.NewServeMux()
// 	mux.Handle("/echo", echo)
// 	return mux
// }
// this version of ServMux is tightly coupled to EchoHandler
// that is unnecessary. We can make it more generic by

type Route interface {
	http.Handler

	// Pattern reports the path at which the route should be registered.
	Pattern() string
}

// NewServeMux builds a ServeMux that will route requests
// to the given routes.
func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}
