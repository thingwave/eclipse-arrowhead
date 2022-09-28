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

const systemnameStr string = "{sysName:[-_a-zA-Z0-9-]+}"
const servicenameStr string = "{srvName:[-_a-zA-Z0-9-]+}"

var clientRoutes = Routes{
	Route{
		"Echo",
		[]string{"GET"},
		"/datamanager/echo",
		Echo,
	},
	Route{
		"ProxyGetSystems",
		[]string{"GET"},
		"/datamanager/proxy",
		ProxyGetSystems,
	},
	Route{
		"HistorianGetSystems",
		[]string{"GET"},
		"/datamanager/historian",
		HistorianGetSystems,
	},
	Route{
		"ProxyGetServices",
		[]string{"GET"},
		"/datamanager/proxy/" + systemnameStr,
		ProxyGetServices,
	},
	Route{
		"HistorianGetServices",
		[]string{"GET"},
		"/datamanager/historian/" + systemnameStr,
		HistorianGetServices,
	},
	Route{
		"ProxyGetServiceData",
		[]string{"GET"},
		"/datamanager/proxy/" + systemnameStr + "/" + servicenameStr,
		ProxyGetServiceData,
	},
	Route{
		"putPData",
		[]string{"PUT"},
		"/datamanager/proxy/"+ systemnameStr + "/" + servicenameStr,
		ProxyPutServiceData,
	},
	Route{
		"DMProxyWShandler",
		[]string{"GET"},
		"/datamanager/ws/proxy/" + systemnameStr + "/" + servicenameStr,
		DMProxyWShandler,
	},
	Route{
		"getHData",
		[]string{"GET"},
		"/datamanager/historian/" + systemnameStr + "/" + servicenameStr,
		HistorianGetServiceData,
	},
	Route{
		"putHData",
		[]string{"PUT"},
		"/datamanager/historian/" + systemnameStr + "/" + servicenameStr,
		HistorianPutServiceData,
	},
	Route{
		"wsHData",
		[]string{"GET"},
		"/datamanager/ws/historian/" + systemnameStr + "/" + servicenameStr,
		DMHistorianWShandler,
	},
}
