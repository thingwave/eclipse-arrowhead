package main

import (
	//"fmt"
	"net/http"
	"github.com/gorilla/mux"
	auth "arrowhead.eu/common/auth"
)

func NewRouter(sslEnabled bool) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range clientRoutes {
		var handler http.Handler
		handler = route.HandlerFunc

		//fmt.Printf("got array %v\n", route.Methods)
		router.Methods(route.Methods[:]...).Path(route.Pattern).Name(route.Name).Handler(handler)
		if sslEnabled {
			router.Use(auth.AuthClientMiddleware)
		}

		/*router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)*/

	}

	return router
}
