package main

import "net/http"

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Client endpoint description
var clientRoutes = Routes{
	Route{
		"Echo",
		[]string{"GET"},
		"/authorization/echo",
		Echo,
	},
	Route{
		"GetPublicKey",
		[]string{"GET"},
		"/authorization/publickey",
		GetPublicKey,
	},
}

// Private endpoints
var privateRoutes = Routes{
	Route{
		"CheckanIntracloudRule",
		[]string{"POST"},
		"/intracloud/check",
		CheckIntraCloudRule,
	},
}

// Management endpoints
var mgmtRoutes = Routes{
	Route{
		"HandleIntracloudRuleByID",
		[]string{"GET", "DELETE"},
		"/intracloud/{id:[0-9]+}",
		HandleIntraCloudRuleByID,
	},
	Route{
		"HandleAllIntracloudRules",
		[]string{"GET", "POST"},
		"/intracloud",
		HandleIntraCloudRule,
	},
}
