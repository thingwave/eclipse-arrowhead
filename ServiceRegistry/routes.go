package main

import(
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
		"/serviceregistry/echo",
		Echo,
	},
	Route{
		"Query",
		[]string{"POST"},
		"/serviceregistry/query",
		Query,
	},
	Route{
		"QueryMulti",
		[]string{"POST"},
		"/serviceregistry/query/multi",
		QueryMulti,
	},
	Route{
		"Register",
		[]string{"POST"},
		"/serviceregistry/register",
		Register,
	},
	Route{
		"Unregister",
		[]string{"DELETE"},
		"/serviceregistry/unregister",
		Unregister,
	},
	Route{
		"RegisterSystem",
		[]string{"POST"},
		"/serviceregistry/register-system",
		RegisterSystem,
	},
	Route{
		"Unregister",
		[]string{"DELETE"},
		"/serviceregistry/uunregister-system",
		UnregisterSystem,
	},
}

// Private endpoints - /serviceregistry/ + <Route below>
var privateRoutes = Routes{
	/*Route{
		"Query",
		[]string{"POST"},
		"/",
		PrivQuery,
	},*/
	Route{
		"PullSystems",
		[]string{"GET"},
		"/pull-systems",
		PrivPullSystems,
	},
	Route{
		"QueryAll",
		[]string{"GET"},
		"/query/all",
		PrivQueryAll,
	},
	Route{
		"QuerySystem",
		[]string{"POST"},
		"/query/system",
		PrivQuerySystem,
	},
	Route{
		"QuerySystemByID",
		[]string{"GET"},
		"/query/system/{id:[0-9]+}",
		PrivQuerySystemById,
	},
}

// Management endpoints - /serviceregistry/mgmt/ + <Route below>
var mgmtRoutes = Routes{
	Route{
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
	},
}
