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
	/*Route{
		"Handle all entries",
		[]string{"POST", "GET"},
		"/",
		HandleEntries,
	},
	Route{
		"Handle all entries by id",
		[]string{"GET", "PUT", "PATCH", "DELETE"},
		"/{id:[0-9]+}",
		HandleEntriesId,
	},

	Route{
		"GetGroupedView",
		[]string{"GET"},
		"/grouped",
		HandleGroupedEntries,
	},

	Route{
		"Handle all interfaces",
		[]string{"GET", "POST"},
		"/interfaces",
		HandleAllInterfaces,
	},
	Route{
		"Handle service interfaces by id",
		[]string{"GET", "PUT", "PATCH", "DELETE"},
		"/interfaces/{id:[0-9]+}",
		HandleSInterfaceById,
	},

	Route{
		"GetEntriedByDefinition",
		[]string{"GET"},
		"/servicedef/{serviceDefinition:[-a-zA-Z0-9]+}",
		HandleEntriesByServiceDefinition,
	},

	Route{
		"Handle all services",
		[]string{"GET", "POST"},
		"/services",
		HandleAllServiceDefs,
	},

	Route{
		"Handle service definition by id",
		[]string{"GET", "PUT", "PATCH", "DELETE"},
		"/services/{id:[0-9]+}",
		HandleServiceDefById,
	},

	Route{
		"Manage systems",
		[]string{"GET", "POST"},
		"/systems",
		HandleAllSystems,
	},

	Route{
		"Handle system",
		[]string{"GET", "PUT", "PATCH", "DELETE"},
		"/systems/{id:[0-9]+}",
		HandleSystemById,
	},*/
}
