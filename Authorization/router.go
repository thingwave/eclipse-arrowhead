package main

import (
//	"fmt"
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
	}

	subRouterP := router.PathPrefix("/authorization/").Subrouter()
	for _, route := range privateRoutes {
		var handler http.Handler
		handler = route.HandlerFunc

		//fmt.Printf("got array %v\n", route.Methods)
		subRouterP.Methods(route.Methods[:]...).Path(route.Pattern).Name(route.Name).Handler(handler)
		if sslEnabled {
			subRouterP.Use(auth.AuthPrivateMiddleware)
		}
	}

	subRouterM := router.PathPrefix("/authorization/mgmt").Subrouter()
	for _, route := range mgmtRoutes {
		var handler http.Handler
		handler = route.HandlerFunc

		//fmt.Printf("got array %v\n", route.Methods)
		subRouterM.Methods(route.Methods[:]...).Path(route.Pattern).Name(route.Name).Handler(handler)
		if sslEnabled {
			subRouterM.Use(auth.AuthManagementMiddleware)
		}
	}

	fs := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return router
}
