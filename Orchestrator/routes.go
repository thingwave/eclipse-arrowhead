package main

import "net/http"

type Route struct {
	Name string
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
		"/orchestrator/echo",
		Echo,
	},
	Route{
		"Orchestrator",
		[]string{"POST"},
		"/orchestrator/orchestration",
		Orchestration,
	},
	Route{
		"StartStoreOrchestration",
		[]string{"GET"},
		"/orchestrator/orchestration/{id:[0-9]+}",
		StartStoreOrchestration,
	},
}

// Private endpoints - /orchestrator/ + <Route below>
var privateRoutes = Routes{
}

// Management endpoints - /orchestrator/mgmt + <Route below>
var mgmtRoutes = Routes{
	Route{
		"HandleAllStoreEntries",
		[]string{"GET", "POST"},
		"/store",
		HandleAllStoreEntries,
	},
	Route{
		"HandleStoreEntryByID",
		[]string{"GET", "DELETE"},
		"/store/{id:[0-9]+}",
		HandleStoreEntryByID,
	},
	Route{
		"GetEntriesbyConsumer",
		[]string{"POST"},
		"/store/all_by_consumer",
		HandleStoreEntrysByConsumer,
	},
	Route{
		"GetTopPriorityEntries",
		[]string{"GET"},
		"/store/all_top_priority",
		HandleStoreEntriesByTopPriority,
	},
	Route{
		"ModifyPriorities",
		[]string{"POST"},
		"/store/modify_priorities",
		HandleStoreModifyPriority,
	},
}
