package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Client endpoints
var clientRoutes = Routes{
	Route{
		"Echo",
		[]string{"GET"},
		"/certificate-authority/echo",
		Echo,
	},
	Route{
		"CheckCertificate",
		[]string{"POST"},
		"/certificate-authority/checkCertificate",
		CheckCertificate,
	},
}

// Private endpoints - /certificate-authority/ + <Route below>
var privateRoutes = Routes{
	Route{
		"SignCSR",
		[]string{"POST"},
		"/sign",
		PrivSign,
	},
	Route{
		"CheckTrustedKey",
		[]string{"POST"},
		"/checkTrustedKey",
		PrivCheckTrustedKey,
	},
}

// Management endpoints - /certificate-authority/mgmt/ + <Route below>
var mgmtRoutes = Routes{
}
