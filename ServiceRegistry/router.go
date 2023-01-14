/********************************************************************************
 * Copyright (c) 2022 Lulea University of Technology
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0.
 *
 * SPDX-License-Identifier: EPL-2.0
 *
 * Contributors:
 *   ThingWave AB - implementation
 *   Arrowhead Consortia - conceptualization
 ********************************************************************************/

package main

import (
	//	"fmt"
	"net/http"

	auth "arrowhead.eu/common/auth"
	"github.com/gorilla/mux"
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

	subRouterP := router.PathPrefix("/serviceregistry/").Subrouter()
	for _, route := range privateRoutes {
		var handler http.Handler
		handler = route.HandlerFunc

		//fmt.Printf("got array %v\n", route.Methods)
		subRouterP.Methods(route.Methods[:]...).Path(route.Pattern).Name(route.Name).Handler(handler)
		if sslEnabled {
			subRouterP.Use(auth.AuthPrivateMiddleware)
		}
	}

	subRouterM := router.PathPrefix("/serviceregistry/mgmt").Subrouter()
	for _, route := range mgmtRoutes {
		var handler http.Handler
		handler = route.HandlerFunc

		//fmt.Printf("got array %v\n", route.Methods)
		subRouterM.Methods(route.Methods[:]...).Path(route.Pattern).Name(route.Name).Handler(handler)
		if sslEnabled {
			subRouterM.Use(auth.AuthManagementMiddleware)
		}
	}

	//fs := http.FileServer(http.Dir("./public"))
	//router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return router
}
