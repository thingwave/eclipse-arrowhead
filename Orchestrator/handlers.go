package main

import (
	"encoding/json"
	"fmt"

	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	dto "arrowhead.eu/common/datamodels"
	"github.com/gorilla/mux"
)

// /////////////////////////////////////////////////////////////////////////////
func Echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	fmt.Fprint(w, "Got it!")
}

// /////////////////////////////////////////////////////////////////////////////
func Orchestration(w http.ResponseWriter, r *http.Request) {
	var request dto.ServiceRequestForm
	var response dto.OrchestrationResponseDTO

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10*1024))
	if err != nil {
		panic(err) //XXX
	}
	if err := r.Body.Close(); err != nil {
		panic(err) //XXX
	}
	fmt.Printf("BODY: %s\n", body)

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("REQ: %+v\n", request)

	//validate request before use!
	//validateOrchestrationRequest()

	// get data from database

	//prepare response
	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not get orchestration\n"))
	}

	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr))
}

// /////////////////////////////////////////////////////////////////////////////
func StartStoreOrchestration(w http.ResponseWriter, r *http.Request) {
	var response dto.OrchestrationResponseDTO

	vars := mux.Vars(r)
	//fmt.Printf("ID: %s\n", vars["id"])

	systemId32, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	systemId := int64(systemId32)

	// do the orchestration process below
	_, err = getSystem(GetOrDB(), systemId)
	if err != nil {
		var errMsg dto.ErrorMessageDTO
		errMsg.ErrorMessage = fmt.Sprintf("System with id %d not found.", systemId)
		errMsg.ErrorCode = 400
		errMsg.ExceptionType = "INVALID_PARAMETER"
		jsonRespStr, _ := json.Marshal(errMsg)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(jsonRespStr))
		return
	}

	response.Response, err = GetOrchestrationForSystem(GetOrDB(), systemId)

	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not get orchestration\n"))
	}

	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr))
}

// /////////////////////////////////////////////////////////////////////////////
func HandleAllStoreEntries(w http.ResponseWriter, r *http.Request) {
	var response dto.StoreEntryList

	switch r.Method {
	case http.MethodGet:
		// fill in the response here

	case http.MethodPost:
		// add store rule here

	}

	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not get orchestration\n"))
	}

	//fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr))
}

// /////////////////////////////////////////////////////////////////////////////
func HandleStoreEntryByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if vars["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//fmt.Printf("ID: %s\n", vars["id"])

	switch r.Method {
	case http.MethodGet:
	case http.MethodDelete:
	}
}

// /////////////////////////////////////////////////////////////////////////////
func HandleStoreEntrysByConsumer(w http.ResponseWriter, r *http.Request) {

}

// /////////////////////////////////////////////////////////////////////////////
// GET /orchestrator/mgmt/store/all_top_priority
func HandleStoreEntriesByTopPriority(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("HandleStoreEntriesByTopPriority()\n")
	vars := r.URL.Query()
	/*fmt.Printf("Len(vars): %d\n", len(vars))
	for k,v := range vars {
		fmt.Printf("%s => %s\n", k, v)
	}
	fmt.Printf("%v\n", vars)*/
	var page string = "-1"
	var item_per_page string = "-1"
	var sort_field string = "id"
	var direction string = "ASC"

	pageRaw, ok := vars["page"]
	if ok {
		page = pageRaw[0]
	}

	sortFieldRaw, ok := vars["sort_field"]
	if ok {
		sort_field = sortFieldRaw[0]
	}
	fmt.Printf("vars::sort_field = %s\n", sort_field)

	directionRaw, ok := vars["direction"]
	if ok {
		direction = directionRaw[0]
	}
	fmt.Printf("vars::direction = %s\n", vars["direction"])

	itemPerPageRaw, ok := vars["item_per_page"]
	if ok {
		item_per_page = itemPerPageRaw[0]
	}
	fmt.Printf("page: %s\nitem_per_page: %s\nsoft_field: %s\ndirection: %s\n", page, item_per_page, sort_field, direction)

	var res dto.StoreEntryList
	var err error
	res.Data, err = GetTopPriorityEntries(GetOrDB())
	if err != nil {
		res.Data = make([]dto.StoreEntry, 0)
		res.Count = 0
		//error return
	}
	res.Count = len(res.Data)

	jsonRespStr, _ := json.Marshal(res)
	//fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr)+"\n")
}

///////////////////////////////////////////////////////////////////////////////
//
/*
type ConsumerRule struct {
	ConsumerSystemId      int64  `json:"consumerSystemId"`
	ServiceDefinitionName string `json:"serviceDefinitionName"`
	ServiceInterfaceName  string `json:"ServiceInterfaceName,omitempty"`
}*/

func HandleStoreModifyPriority(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("HandleStoreModifyPriority\n")
	var request dto.PriorityList

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10*1024))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	//fmt.Printf("BODY: %s\n", body)

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Printf("REQ: %+v\nlen():%d\n", request, len(request.PriorityMap))
	for k, v := range request.PriorityMap {
		//fmt.Printf("CHANGE_PRIORITY(%s) => %d\n", k, v)
		UpdatePriorityForSystem(GetOrDB(), k, v)
	}
}
