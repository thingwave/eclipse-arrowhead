package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	//"strings"
	//"strconv"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	dto "arrowhead.eu/common/datamodels"
)

type ProxyElement struct {
	SysName string
	SrvName string
	Message string
}

var proxyElements = []ProxyElement{}

///////////////////////////////////////////////////////////////////////////////
//
func Echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	fmt.Fprint(w, "Got it!")
}

///////////////////////////////////////////////////////////////////////////////
//
func ProxyGetSystems(w http.ResponseWriter, r *http.Request) {
	log.Println("ProxyGetSystems:")
	var result SystemList
	result.Systems = make([]string, 0)

	/*vars := mux.Vars(r)
	if vars["sysName"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	//fmt.Fprintf(w, "SysName %v\n", vars["sysName"])

	for _, e := range proxyElements {
		//if e.SysName == vars["sysName"] {
		//fmt.Printf("found service %s for system %s\n", e.SrvName, vars["sysName"])
		if isSystemInList(e.SysName, result.Systems) == false {
			result.Systems = append(result.Systems, e.SysName)
		}
		//}
	}

	jsonRespStr, _ := json.Marshal(result)
	//fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonRespStr))
}

var upgrader = websocket.Upgrader{ReadBufferSize: 10 * 1024, WriteBufferSize: 10 * 1024}

func DMHistorianWShandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DMHistorianWShandler()\n")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WS upgrade failed!\n")
		return
	}
	defer ws.Close()

	vars := mux.Vars(r)
	sysName := vars["sysName"]
	srvName := vars["srvName"]
	fmt.Printf("SysName %v\nSrvName: %v\n", sysName, srvName)

	// handle connection
	for {
		var request []dto.SenMLEntry
		mt, payload, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Printf("WS packet input:\n%s\n", string(payload))

		err = json.Unmarshal(payload, &request)
		if err != nil {
			fmt.Printf("Error decoding JSON\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = validateSenML(request)
		if err != nil {
			log.Printf("Invalid SenML message: %v\n", err)
			return
		}

		body, err := json.Marshal(request)
		if err != nil {
			log.Println("Could not Marshal")
			return
		}

		fmt.Printf("body is: %s\n", string(body[:]))
		err = PutDMHistSystemServiceData(GetDMDB(), vars["sysName"], vars["srvName"], string(body[:]), request)

		// echo message back to client for now!
		ws.WriteMessage(mt, payload)
	}

}

///////////////////////////////////////////////////////////////////////////////
//
func HistorianGetSystems(w http.ResponseWriter, r *http.Request) {
	var response SystemList
	var err error

	response.Systems, err = GetDMHistSystems(GetDMDB())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not get system list\n"))
		return
	}

	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not get system list\n"))
		return
	}

	//fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonRespStr))
}

///////////////////////////////////////////////////////////////////////////////
//
func ProxyGetServices(w http.ResponseWriter, r *http.Request) {
	var response SystemServiceList
	response.Services = []string{}

	vars := mux.Vars(r)
	if vars["sysName"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Missing system name\n"))
		return
	}

	log.Printf("ProxyGetServices(%s):\n", vars["sysName"])

	for _, e := range proxyElements {
		if vars["sysName"] == e.SysName {
			response.Services = append(response.Services, e.SrvName)
		}
	}

	jsonRespStr, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonRespStr))
}

///////////////////////////////////////////////////////////////////////////////
//
func HistorianGetServices(w http.ResponseWriter, r *http.Request) {
	var response SystemServiceList

	log.Printf("HistorianGetServices:\n")

	vars := mux.Vars(r)
	if vars["sysName"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Missing system name\n"))
		return
	}

	log.Printf("HistorianGetServices(%s):\n", vars["sysName"])
	//fmt.Fprintf(w, "SysName %v\n", vars["sysName"])
	response.Services = GetDMHistSystemServices(GetDMDB(), vars["sysName"])
	//fmt.Fprintf(w, "system %v\n", systems[0])

	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not generate JSON\n"))
		return
	}

	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonRespStr))
}

///////////////////////////////////////////////////////////////////////////////
//
func ProxyGetServiceData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["sysName"] == "" || vars["srvName"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("ProxyGetServiceData('%s', '%s')\n", vars["sysName"], vars["srvName"])

	for _, e := range proxyElements {
		if e.SysName == vars["sysName"] && e.SrvName == vars["srvName"] {
			w.Header().Add("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", e.Message)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

///////////////////////////////////////////////////////////////////////////////
//
func HistorianGetServiceData(w http.ResponseWriter, r *http.Request) {
	var response = []dto.SenMLEntry{}

	vars := mux.Vars(r)

	var count int = 1
	countStr := r.URL.Query().Get("count")
	if countStr != "" {
		fmt.Printf("count: %v\n", countStr)
		count2, err2 := strconv.Atoi(countStr)
		if err2 != nil {
			return
		}
		count = count2
	}
	/*var from float64 = -1
	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		fmt.Printf("from: %v\n", fromStr)
		from2, err2 := strconv.ParseFloat(fromStr, 64)
		if err2 != nil {
			return
		}
		from = from2
	}*/
	/*var to float64 = -1
	toStr := r.URL.Query().Get("to")
	if countStr != "" {
		fmt.Printf("to: %v\n", toStr)
		to2, err2 := strconv.ParseFloat(toStr, 64)
		if err2 != nil {
			return
		}
		to = to2
	}*/

	signals := make([]SignalProperties, 0)
	var sigXcnt = 0
	params := r.URL.Query()
	//fmt.Printf("%v", params)
	for {
		sigXname := fmt.Sprintf("sig%d", sigXcnt)
		//fmt.Println(sigXname)

		found := false
		for k, v := range params {
			if k == sigXname {
				fmt.Printf("Found: %s => %s\n", sigXname, v[0])

				if v[0] == "" {
					return
				}
				sigProp := SignalProperties{SigName: k, SigKey: v[0], SigCount: 1}
				signals = append(signals, sigProp)
				found = true
			}
		}
		sigXcnt += 1

		if !found {
			break
		}
	}

	sigXcnt = 0
	for {
		sigXname := fmt.Sprintf("sig%d", sigXcnt)
		sigXvalue := fmt.Sprintf("sig%dcount", sigXcnt)
		//fmt.Println(sigXname)

		for k, v := range params {
			if k == sigXvalue {
				//fmt.Printf("Found: %s => %s\n", sigXname, v[0])

				// update the signal list
				for i, _ := range signals {
					if signals[i].SigName == sigXname {
						//fmt.Printf("Update %s with count: %s\n", sigXname, v[0])
						count2, err2 := strconv.Atoi(v[0])
						if err2 != nil {
							fmt.Printf("Error decoding signal count\n")
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						signals[i].SigCount = count2
						break
					}
				}

				break
			}
		}
		sigXcnt += 1

		if sigXcnt > len(signals) {
			break
		}
	}

	for _, sig := range signals {
		fmt.Printf("%s :: %v\n", sig.SigName, sig.SigCount)
	}
	/*for k, v := range params {

		sigXname := fmt.Sprintf("\nsig%d", sigXcnt)
		fmt.Println(sigXname)
		fmt.Println(k)
		if k == sigXname {
			fmt.Printf("Found: %s\n", sigXname)
		}
		sigXcnt += 1
		fmt.Println(k, " => ", v)
	}*/

	if len(signals) > 0 {
		log.Println("User requested signals, call alt. db method!")

		response, _ = GetDMHistSystemServiceDataSignals(GetDMDB(), vars["sysName"], vars["srvName"], -1, -1, count, signals)
	} else {

		body, _ := GetDMHistSystemServiceData(GetDMDB(), vars["sysName"], vars["srvName"], -1, -1, count)
		for i, e := range body {
			fmt.Printf("%v: %v\n", i, e)
			response = append(response, e)
		}
	}

	
	//fmt.Fprintf(w, "SysName %v\nSrvName: %v\n", vars["sysName"], vars["srvName"])
	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//fmt.Fprint(w, string("Could n\n"))
		return
	}

	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonRespStr))
}

///////////////////////////////////////////////////////////////////////////////
//
func ProxyPutServiceData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyPutServiceData()")
	var request []dto.SenMLEntry

	vars := mux.Vars(r)
	sysName := vars["sysName"]
	srvName := vars["srvName"]
	fmt.Printf("\tsysName %v\nsrvName: %v\n", sysName, srvName)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Printf("BODY: %s\n", body)

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("Error decoding JSON\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateSenML(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err = json.Marshal(request)
	updateProxyData(sysName, srvName, body)

	fmt.Printf("ProxyPutServiceData REQ: %s\n", body)
	w.WriteHeader(http.StatusOK)
}

func updateProxyData(sysName string, srvName string, senml []byte) {
	if isSystemServiceInList(sysName, srvName) == false {
		fmt.Printf("Adding new element!\n")
		var newElement ProxyElement
		newElement.SysName = sysName
		newElement.SrvName = srvName
		newElement.Message = string(senml)
		proxyElements = append(proxyElements, newElement)
	} else {
		for i, e := range proxyElements {
			if e.SysName == sysName && e.SrvName == srvName {
				//fmt.Printf("Updating element\n")
				proxyElements[i].Message = string(senml)
			}
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
//
func DMProxyWShandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DMProxyWShandler()\n")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade failed!\n")
		return
	}
	defer ws.Close()

	vars := mux.Vars(r)
	sysName := vars["sysName"]
	srvName := vars["srvName"]
	fmt.Printf("SysName %v\nSrvName: %v\n", sysName, srvName)

	// handle connection
	for {
		var request []dto.SenMLEntry
		mt, payload, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Printf("%s\n", string(payload))

		err = json.Unmarshal(payload, &request)
		if err != nil {
			fmt.Printf("Error decoding JSON\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = validateSenML(request)
		if err != nil {
			log.Printf("Invalid SenML message: %v\n", err)
			return
		}

		body, err := json.Marshal(request)
		updateProxyData(sysName, srvName, body)

		// echo message back to client for now!
		ws.WriteMessage(mt, payload)
	}
}

///////////////////////////////////////////////////////////////////////////////
//
func HistorianPutServiceData(w http.ResponseWriter, r *http.Request) {
	var request []dto.SenMLEntry

	vars := mux.Vars(r)
	fmt.Printf("SysName %v\nSrvName: %v\n", vars["sysName"], vars["srvName"])

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Printf("BODY: %s\n", body)

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("Error decoding JSON\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateSenML(request)
	if err != nil {
		log.Println("Malformed SenML")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, e := range request {
		fmt.Printf("%v: %+v\n", i, e)
	}

	log.Printf("%s\n", string(body))
	err = PutDMHistSystemServiceData(GetDMDB(), vars["sysName"], vars["srvName"], string(body), request)

	fmt.Printf("HistorianPutServiceData REQ: %+v\n", request)
	w.WriteHeader(http.StatusOK)
}

func isSystemInList(sysName string, sysList []string) bool {

	for _, e := range sysList {
		if e == sysName {
			return true
		}
	}

	return false
}

func isSystemServiceInList(sysName string, srvName string) bool {

	for _, e := range proxyElements {
		if e.SysName == sysName && e.SrvName == srvName {
			return true
		}
	}

	return false
}
