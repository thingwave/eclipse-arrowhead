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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	dto "arrowhead.eu/common/datamodels"
	"github.com/gorilla/mux"
)

type LoginMsg struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetSystems(w http.ResponseWriter, r *http.Request) {
	// check JWT!

	w.Header().Set("Content-type", "application/json")

	systems := fetchAllSystems(GetSRDB(), 1, 1000)
	retJson, _ := json.Marshal(systems)
	fmt.Fprint(w, string(retJson))
}

func Login(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	//r.ParseForm()
	//username := r.Form.Get("username")
	//password := r.Form.Get("password")
	var loginMsg LoginMsg
	_ = json.NewDecoder(r.Body).Decode(&loginMsg)
	fmt.Println("Username: " + loginMsg.Username)
	fmt.Println("Password: " + loginMsg.Password)

	if loginMsg.Username == "jench" && loginMsg.Password == "supersecret" {

	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//w.Header().Set("Location", "dashboard.html")
	//w.WriteHeader(http.StatusFound)
	w.Header().Set("Content-type", "application/jwt")
	hdr := JWT_hdr{
		Alg: "HS256",
		Typ: "JWT",
	}
	payload := JWT_payload{
		Iss:   "https://127.0.0.1:8443", //XXX fixme!
		Exp:   time.Now().Unix() + 3600,
		Email: "jens.eliasson@thingwave.eu",
	}
	retJsonHdr, _ := json.Marshal(hdr)
	retJsonPayload, _ := json.Marshal(payload)
	//fmt.Fprint(w, string(retJson))
	fmt.Fprint(w, base64.StdEncoding.EncodeToString(retJsonHdr)+"."+base64.StdEncoding.EncodeToString(retJsonPayload)+".")
}

/// Arrowhead services ////////////////////////////////////////////
func Echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	fmt.Fprint(w, "Got it!")
}

func getBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	fmt.Printf("BODY: %s\n", body)

	return body, err
}

func Query(w http.ResponseWriter, r *http.Request) {
	var queryReq ServiceQueryForm
	var ret dto.ServiceQueryResultDTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := getBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &queryReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {

		if queryReq.ServiceDefinitionRequirement == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//fmt.Printf("################################\nQuery():\n %+v\n################################\n", queryReq)
		var unfilteredHits int = 0
		ret.ServiceQueryData, _ = queryServicesForName(GetSRDB(), queryReq, &unfilteredHits) //XX BUG HERE SOMEWHERE
		ret.UnfilteredHits = unfilteredHits
		retJson, _ := json.Marshal(ret)

		//fmt.Println(string(retJson))
		fmt.Fprint(w, string(retJson))
		return
	}

}

func QueryMulti(w http.ResponseWriter, r *http.Request) {

}

func Register(w http.ResponseWriter, r *http.Request) {
	var regReq dto.ServiceRegistryEntryDTO

	log.Printf(("\nRegister()\n"))

	body, err := getBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n%s\n<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", string(body))
	err = json.Unmarshal(body, &regReq)
	fmt.Printf("\n###\nRegistration request: %+v\n", regReq)

	if checkClientCN(r, regReq.ProviderSystem.SystemName) == false {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, string("Incorrect CN\n"))
		return
	}

	err = registerService(w, r, regReq)
}

func registerService(w http.ResponseWriter, r *http.Request, regReq dto.ServiceRegistryEntryDTO) error {

	err := validateRegistryEntry(regReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return err
	}

	response, err := registerServiceForSystem(GetSRDB(), regReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not register service\n"))
		return err
	}

	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not register service\n"))
		return err
	}

	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr)+"\n")

	return nil
}

func Unregister(w http.ResponseWriter, r *http.Request) {
	log.Printf("Unregister()\n")
	/*params := r.URL.Query()
	  for i, s := range params {
	    fmt.Printf("%d -> %s\n", i, s)
	  }*/

	service_definition := ""
	if len(r.URL.Query()["service_definition"]) == 1 {
		service_definition = strings.TrimSpace(r.URL.Query()["service_definition"][0])
	}
	system_name := ""
	if len(r.URL.Query()["system_name"]) == 1 {
		system_name = strings.TrimSpace(r.URL.Query()["system_name"][0])
	}
	address := ""
	if len(r.URL.Query()["address"]) == 1 {
		address = strings.TrimSpace(r.URL.Query()["address"][0])
		var a, b, c, d int
		_, serr := fmt.Sscanf(address, "%d.%d.%d.%d", &a, &b, &c, &d)
		if serr != nil {
			address = ""
		}
		if a < 0 || a > 255 || b < 0 || b > 255 || c < 0 || c > 255 || d < 0 || d > 255 {
			address = ""
		}

		// add check for illegal addresses here! (INET_ADDR_ANY, broadcast, multicast, etc...
		// checkValidAddress(address)

	}
	port := ""
	if len(r.URL.Query()["port"]) == 1 {
		port = strings.TrimSpace(r.URL.Query()["port"][0])
	}

	fmt.Println("service_definition: ", service_definition)
	fmt.Println("system_name: ", system_name)
	fmt.Println("address: ", address)
	fmt.Println("port: ", port)

	// check so all parameters are present
	if len(service_definition) == 0 || len(system_name) == 0 || len(address) == 0 || len(port) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Missing parameters\n"))
		return
	}

	porti, erri := strconv.Atoi(port)
	if erri != nil || porti < 0 || porti > 65535 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Port must be an integer between 0 and 65535\n"))
		return
	}

	if checkClientCN(r, system_name) == false {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, string("Incorrect CN\n"))
		return
	}

	ok, err := unregisterServiceForSystem(GetSRDB(), service_definition, system_name, address, porti)
	if ok == false || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Bad request\n"))
		return
	}

}

func RegisterSystem(w http.ResponseWriter, r *http.Request) {
	var regReq dto.SystemRequestDTO

	log.Printf(("\nRegisterSystem()\n"))

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%s\n", string(body))
	err = json.Unmarshal(body, &regReq)
	fmt.Printf("\n###\nRegistrationSystem request: %+v\n", regReq)

	if checkClientCN(r, regReq.SystemName) == false {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, string("Incorrect CN\n"))
		return
	}

	response, err := registerSystem(GetSRDB(), regReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonRespStr, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string("Could not register service\n"))
		return
	}

	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr)+"\n")

}

func UnregisterSystem(w http.ResponseWriter, r *http.Request) {

	system_name := ""
	if len(r.URL.Query()["system_name"]) == 1 {
		system_name = strings.TrimSpace(r.URL.Query()["system_name"][0])
	}
	address := ""
	if len(r.URL.Query()["address"]) == 1 {
		address = strings.TrimSpace(r.URL.Query()["address"][0])
		var a, b, c, d int
		_, serr := fmt.Sscanf(address, "%d.%d.%d.%d", &a, &b, &c, &d)
		if serr != nil {
			address = ""
		}
		if a < 0 || a > 255 || b < 0 || b > 255 || c < 0 || c > 255 || d < 0 || d > 255 {
			address = ""
		}

		// add check for illegal addresses here! (INET_ADDR_ANY, broadcast, multicast, etc...
		// checkValidAddress(address)

	}
	port := ""
	if len(r.URL.Query()["port"]) == 1 {
		port = strings.TrimSpace(r.URL.Query()["port"][0])
	}

	fmt.Println("system_name: ", system_name)
	fmt.Println("address: ", address)
	fmt.Println("port: ", port)

	// check so all parameters are present
	if len(system_name) == 0 || len(address) == 0 || len(port) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Illegal parameters\n"))
		return
	}

	porti, erri := strconv.Atoi(port)
	if erri != nil || porti < 0 || porti > 65535 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Port must be an integer between 0 and 65535\n"))
		return
	}

	if checkClientCN(r, system_name) == false {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, string("Incorrect CN\n"))
		return
	}

	ok, err := unregisterSystem(GetSRDB(), system_name, address, porti)
	if ok == false || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string("Bad request\n"))
		return
	}
}

// checkClientCN checks that the client CN matches with the provided system name
// if the two matches, true is returned. false is returned otherwise
func checkClientCN(r *http.Request, systemName string) bool {
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		cn := strings.ToLower(r.TLS.PeerCertificates[0].Subject.CommonName)

		if systemName == cn && len(systemName) >= 1 {
			return true
		}
	}
	return true //false XXX
}

///////////////////////////////////////////////////////////////////////////////
//                            PRIVATE endpoints

// GET /serviceregistry/pull-systems
func PrivPullSystems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Service Priv::Pull-systems request")
	var response dto.SystemListDTO

	systems, err := getAllSystems(GetSRDB(), "ASC")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Data = systems
	response.Count = len(response.Data)

	jsonRespStr, _ := json.Marshal(response)
	//fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr))
}

// GET /serviceregistry/query/all
func PrivQueryAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Service Priv::QueryAll request")
	var response dto.ServiceRegistryListResponseDTO

	services, err := getAllServices(GetSRDB(), nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Data = services
	response.Count = len(response.Data)

	jsonRespStr, _ := json.Marshal(response)
	//fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr))
}

// /serviceregistry/query/system
func PrivQuerySystem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Service Priv::QuerySystem request")

	var req dto.SystemRequestDTO

	body, err := getBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%s\n", string(body))
	err = json.Unmarshal(body, &req)
	fmt.Printf("QuerySystemReq: %+v\n", req)

	systemID := getIDforSystem(GetSRDB(), req.SystemName, req.Address, req.Port)
	fmt.Printf("SystemID: %d\n", systemID)

	if err != nil || systemID == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// get updated information
	response, err := getSystem(GetSRDB(), systemID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonRespStr, _ := json.Marshal(response)
	fmt.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr)+"\n")
}

// GET /serviceregistry/query/system/{id}
func PrivQuerySystemById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var systemIDstr = vars["id"]
	if systemIDstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	systemId32, err := strconv.Atoi(systemIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	systemId := int64(systemId32)
	fmt.Printf("\nQuerySystemById(%v)\n", systemId)

	ret, err := getSystem(GetSRDB(), systemId)
	if err != nil {
		fmt.Println("No such system")
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

	jsonRespStr, _ := json.Marshal(ret)
	log.Println(string(jsonRespStr))
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonRespStr))
}

///////////////////////////////////////////////////////////////////////////////
//                            MANAGEMENT endpoints

///////////////////////////////////////////////////////////////////////////////
//

type ServiceRegistryEntry struct {
	Id        int64  `json:"data"`
	Version   int32  `json:"version"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type ServiceRegistryEntryList struct {
	Data  []ServiceRegistryEntry `json:"data"`
	Count int32                  `json:"count"`
}

// /serviceregistry/mgmt
func HandleEntries(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleEntries()")

	if r.Method == http.MethodGet {
		var response dto.ServiceRegistryListResponseDTO

		response.Data, _ = getAllServices(GetSRDB(), nil, nil)
		response.Count = len(response.Data)
		jsonRespStr, _ := json.Marshal(response)

		fmt.Println(string(jsonRespStr))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonRespStr)+"\n")
		return
	} else if r.Method == http.MethodPost {
		var regReq dto.ServiceRegistryEntryDTO

		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("%s\n", string(body))
		err = json.Unmarshal(body, &regReq)
		fmt.Printf("REGREQ: %+v\n", regReq)

		if checkClientCN(r, "SYSOP") == false {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, string("Incorrect CN\n"))
			return
		}

		err = registerService(w, r, regReq)
	}
}

// /serviceregistry/mgmt/{id}
func HandleEntriesId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var serviceIDstr = vars["id"]
	if serviceIDstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceId32, err := strconv.Atoi(serviceIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceId := int64(serviceId32)
	fmt.Printf("HandleEntriesId(%d)\n", serviceId)
	//fmt.Printf("Service ID: %d\n", serviceId)

	switch r.Method {
	case http.MethodGet:
		respDTO, err := fetchServiceById(GetSRDB(), serviceId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(respDTO)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
		return
	case http.MethodPut:
		request := dto.ServiceRegistryEntryDTO{}
		system := dto.ServiceRegistryEntryDTOIncomplete{}

		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("%s\n", string(body))
		err = json.Unmarshal(body, &system)
		fmt.Printf("\tMGT::REG-REQ: %+v\n", system)

		// convert incomplete to dto
		if system.ServiceDefinition != "" {
			request.ServiceDefinition = system.ServiceDefinition
		}
		if system.ProviderSystem != nil {
			request.ProviderSystem = *system.ProviderSystem
		}
		if system.ServiceUri != "" {
			request.ServiceUri = system.ServiceUri
		}
		if system.EndOfValidity != "" {
			request.EndOfValidity = system.EndOfValidity
		}
		if system.Secure != "" {
			request.Secure = system.Secure
		}
		if system.Metadata != nil {
			request.Metadata = system.Metadata
		}
		if system.Version != nil {
			request.Version = *system.Version
		}
		if len(system.Interfaces) > 0 {
			request.Interfaces = system.Interfaces
		}

		ok, err := updateServiceEntry(GetSRDB(), request, serviceId)
		if ok == false || err != nil {
			return
		}

	case http.MethodPatch:
		//request := dto.ServiceRegistryEntryDTO{}
		request, err := fetchServiceById(GetSRDB(), serviceId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		system := dto.ServiceRegistryEntryDTOIncomplete{}

		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//fmt.Printf("%s\n", string(body))
		err = json.Unmarshal(body, &system)
		fmt.Printf("\tMGT::REG-REQ: %+v\n", system)
		fmt.Printf("\tMGT::REG-CUR: %+v\n", request)

		// convert incomplete to dto
		if system.ServiceDefinition != "" {
			request.ServiceDefinition = system.ServiceDefinition
		}
		if system.ProviderSystem != nil {
			request.ProviderSystem = *system.ProviderSystem
		}
		if system.ServiceUri != "" {
			request.ServiceUri = system.ServiceUri
		}
		if system.EndOfValidity != "" {
			request.EndOfValidity = system.EndOfValidity
		}
		if system.Secure != "" {
			request.Secure = system.Secure
		}
		if system.Metadata != nil {
			request.Metadata = system.Metadata
		}
		if system.Version != nil {
			request.Version = *system.Version
		}
		if len(system.Interfaces) > 0 {
			request.Interfaces = system.Interfaces
		}

		ok, err := updateServiceEntry(GetSRDB(), request, serviceId)
		fmt.Printf("ok: %v\n", ok)
		fmt.Println(err)
		if ok == false || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		err = deleteServiceById(GetSRDB(), serviceId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

type ServiceRegistryGroupedDTO struct {
	ServicesGroupedBySystems           []ServicesGroupedBySystemDTO            `json:"servicesGroupedBySystems"`
	ServicesGroupedByServiceDefinition []ServicesGroupedByServiceDefinitionDTO `json:"servicesGroupedByServiceDefinition"`
	AutoCompleteData                   AutoCompleteDataDTO                     `json:"autoCompleteData"`
}

type ServicesGroupedBySystemDTO struct {
	SystemId   int64                            `json:"systemId"`
	SystemName string                           `json:"systemName"`
	Address    string                           `json:"address"`
	Port       int                              `json:"port"`
	Services   []dto.ServiceRegistryResponseDTO `json:"services"`
}

type ServicesGroupedByServiceDefinitionDTO struct {
	ServiceDefinitionId int64  `json:"serviceDefinitionId"`
	ServiceDefinition   string `json:"serviceDefinition"`
}

type AutoCompleteDataDTO struct {
	ServiceList   []dto.ServiceDTO        `json:"serviceList"`
	SystemList    []dto.SystemResponseDTO `json:"systemList"`
	InterfaceList []dto.InterfaceDTO      `json:"interfaceList"`
}

// /serviceregistry/mgmt/grouped
func HandleGroupedEntries(w http.ResponseWriter, r *http.Request) {
	var ret ServiceRegistryGroupedDTO

	fmt.Println("HandleGroupedEntries()")

	ret.ServicesGroupedBySystems = make([]ServicesGroupedBySystemDTO, 0)
	systems, err := getAllSystems(GetSRDB(), "ASC")
	if err != nil {

	}
	for _, system := range systems {
		var sys ServicesGroupedBySystemDTO
		sys.SystemId = system.Id
		sys.SystemName = system.SystemName
		sys.Address = system.Address
		sys.Port = system.Port
		sys.Services, _ = getAllServicesBySystem(GetSRDB(), system.Id)

		ret.ServicesGroupedBySystems = append(ret.ServicesGroupedBySystems, sys)
	}

	ret.ServicesGroupedByServiceDefinition = make([]ServicesGroupedByServiceDefinitionDTO, 0)

	ret.AutoCompleteData.InterfaceList, _ = getAllInterfaces(GetSRDB())
	ret.AutoCompleteData.ServiceList, _ = getAllServiceDefinitionsSimple(GetSRDB())
	ret.AutoCompleteData.SystemList, _ = getAllSystems(GetSRDB(), "ASC")

	w.Header().Add("Content-Type", "application/json")
	jsonRespStr, _ := json.Marshal(ret)
	fmt.Println(string(jsonRespStr))
	fmt.Fprint(w, string(jsonRespStr))
}

type InterfaceNameB struct {
	InterfaceName string `json:"interfaceName"`
}

// /serviceregistry/mgmt/interfaces
func HandleAllInterfaces(w http.ResponseWriter, r *http.Request) {
	var err error

	fmt.Println("HandleAllInterfaces()")

	switch r.Method {
	case http.MethodGet:
		var response dto.ServiceInterfaceListResponseDTO
		response.Data, err = getInterfacesList(GetSRDB())
		if err != nil {

			return
		}
		response.Count = len(response.Data)

		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(response)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
		return
	case http.MethodPost:
		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var request InterfaceNameB
		err = json.Unmarshal(body, &request)
		fmt.Printf("%+v\n", request)

		ret, err := getInterfaceByName(GetSRDB(), request.InterfaceName)
		if ret.Id == -1 {
			// n
			//addInterfaceByName(GetSRDB(), request.InterfaceName)
			interfaceId, _ := insertOrUpdateServiceInterface(GetSRDB(), request.InterfaceName)
			fmt.Printf("interfaceId: %v, err: %s\n", interfaceId, err)
			ret, _ := getInterfaceByName(GetSRDB(), request.InterfaceName)

			w.Header().Add("Content-Type", "application/json")
			jsonRespStr, _ := json.Marshal(ret)
			fmt.Println(string(jsonRespStr))
			fmt.Fprint(w, string(jsonRespStr))
		} else {
			w.WriteHeader(http.StatusBadRequest) //XXA Ah error system alreayd exists
			return
		}
	}
}

// //serviceregistry/mgmt/interfaces/{id}
func HandleSInterfaceById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var idS = vars["id"]
	var serviceDefId int64 = -1
	var err error

	if idS != "" {
		var tmp int
		tmp, err = strconv.Atoi(idS)
		serviceDefId = int64(tmp)
		if err != nil {
			return
		}
	}
	fmt.Printf("HandleSInterfaceById(%v)\n", serviceDefId) //XXX validate id

	switch r.Method {
	case http.MethodGet:
		ret, err := getInterfaceById(GetSRDB(), serviceDefId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(ret)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
	case http.MethodPut:
		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(body)
		var request InterfaceNameB
		err = json.Unmarshal(body, &request)
		fmt.Printf("%+v\n", request)
		err = updateInterfaceById(GetSRDB(), serviceDefId, request.InterfaceName)
		ret, err := getInterfaceById(GetSRDB(), serviceDefId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(ret)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))

	case http.MethodDelete:
		err = deleteInterfaceById(GetSRDB(), serviceDefId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case http.MethodPatch:
		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(body)
		var request InterfaceNameB
		err = json.Unmarshal(body, &request)
		fmt.Printf("%+v\n", request)
		err = updateInterfaceById(GetSRDB(), serviceDefId, request.InterfaceName)
		ret, err := getInterfaceById(GetSRDB(), serviceDefId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(ret)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
	}

}

// /mgmt/servicedef/{serviceDefinition}
func HandleEntriesByServiceDefinition(w http.ResponseWriter, r *http.Request) { //XXX MUST IMPLEMENT!
	vars := mux.Vars(r)
	var serviceDefinition = vars["serviceDefinition"]
	//var err error

	var pageS = vars["page"]
	var page int = -1 //XXX
	var pageP *int = nil
	var items_per_pageS = vars["items_per_page"]
	var items_per_page int = 50 //XXX
	var items_per_pageP *int = nil
	var err error
	if pageS != "" {
		page, err = strconv.Atoi(pageS)
		if err != nil {
			return
		}
		pageP = &page
	}

	if items_per_pageS != "" {
		items_per_page, err = strconv.Atoi(items_per_pageS)
		if err != nil {
			return
		}
		items_per_pageP = &items_per_page
	}

	fmt.Printf("HandleEntriesByServiceDefinition('%s')\n", serviceDefinition)
	if len(serviceDefinition) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ret dto.ServiceRegistryListResponseDTO
	ret.Data, _ = getAllServicesFromServiceDefinition(GetSRDB(), serviceDefinition, pageP, items_per_pageP)
	ret.Count = len(ret.Data)

	w.Header().Add("Content-Type", "application/json")
	jsonRespStr, _ := json.Marshal(ret)
	fmt.Println(string(jsonRespStr))
	fmt.Fprint(w, string(jsonRespStr))
}

// /mgmt/services
func HandleAllServiceDefs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("HandleAllServiceDefs()\n")

	switch r.Method {
	case http.MethodGet:
		var err error
		var ret dto.ServiceDefinitionResponseListDTO
		ret.Data, err = getAllServiceDefinitions(GetSRDB())
		if err != nil {

		}
		ret.Count = int64(len(ret.Data))
		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(ret)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
		return

	case http.MethodPost:
		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var update serviceDefUpdate
		err = json.Unmarshal(body, &update)
		fmt.Printf("%+v\n", update)
		serviceDefId, err := addServiceDefinition(GetSRDB(), update.ServiceDefinition)
		if err != nil {
			//return error msg
			log.Println(err)
		} else {
			resp, err := getServiceDefinitionForService(GetSRDB(), serviceDefId)
			if err != nil {

			}

			w.Header().Add("Content-Type", "application/json")
			jsonRespStr, _ := json.Marshal(resp)
			fmt.Println(string(jsonRespStr))
			fmt.Fprint(w, string(jsonRespStr))
			return
		}
	}

}

type serviceDefUpdate struct {
	ServiceDefinition string `json:"serviceDefinition"`
}

// /mgmt/services/{id}
func HandleServiceDefById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var idS = vars["id"]
	var serviceDefId int64 = -1
	var err error

	if idS != "" {
		var tmp int
		tmp, err = strconv.Atoi(idS)
		serviceDefId = int64(tmp)
		if err != nil {
			return
		}
	}
	fmt.Printf("HandleServiceDefById(%v)\n", serviceDefId)

	switch r.Method {
	case http.MethodGet:
		resp, err := getServiceDefinitionForService(GetSRDB(), serviceDefId)
		if err != nil {

		}

		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(resp)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
		return

	case http.MethodPut:
		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var update serviceDefUpdate
		err = json.Unmarshal(body, &update)
		if update.ServiceDefinition != "" {
			err = updateServiceDefinitionById(GetSRDB(), serviceDefId, update.ServiceDefinition)
			fmt.Println(err)
			w.WriteHeader(http.StatusOK) //XXX check resturn code!
		}

		resp, err := getServiceDefinitionForService(GetSRDB(), serviceDefId)
		if err != nil {

		}

		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(resp)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
		return
	case http.MethodPatch:
		body, err := getBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var update serviceDefUpdate
		err = json.Unmarshal(body, &update)
		if update.ServiceDefinition != "" {
			updateServiceDefinitionById(GetSRDB(), serviceDefId, update.ServiceDefinition)
		}

		resp, err := getServiceDefinitionForService(GetSRDB(), serviceDefId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		jsonRespStr, _ := json.Marshal(resp)
		fmt.Println(string(jsonRespStr))
		fmt.Fprint(w, string(jsonRespStr))
		return
	case http.MethodDelete:
		err = deleteServiceDefinitionById(GetSRDB(), serviceDefId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusInternalServerError)
}

/*
//
func HandleService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var systemIDstr = vars["id"]
	if systemIDstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:
	}

	w.WriteHeader(http.StatusInternalServerError)
}
*/

///////////////////////////////////////////////////////////////////////////////
//
func HandleAllSystems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleAllSystems()")

	switch r.Method {
	case http.MethodGet:
		var response dto.SystemListDTO

		vars := mux.Vars(r)
		var pageS = vars["page"]
		var page int = -1 //XXX
		var items_per_pageS = vars["items_per_page"]
		var items_per_page int = 500 //XXX
		var err error
		if pageS != "" {
			page, err = strconv.Atoi(pageS)
			if err != nil {
				return
			}
		}

		if items_per_pageS != "" {
			items_per_page, err = strconv.Atoi(items_per_pageS)
			if err != nil {
				return
			}
		}

		/*var sort_field = vars["sort_field"]
		  var direction = vars["direction"]*/

		response.Data = fetchAllSystems(GetSRDB(), page, items_per_page)
		response.Count = len(response.Data)
		if response.Count == 0 {
			response.Data = []dto.SystemResponseDTO{}
		}

		//w.WriteHeader(http.StatusInternalServerError)
		jsonRespStr, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(string(jsonRespStr))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonRespStr)+"\n")

	case http.MethodPost:
		var system dto.SystemRequestDTO
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := r.Body.Close(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("%s\n", string(body))
		err = json.Unmarshal(body, &system)
		fmt.Printf("SystemCreate: %+v\n", system)
		if system.SystemName == "" || system.Port < 0 || system.Port > 65535 || system.Address == "" {
			fmt.Printf("\tIllegal request\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response, err := registerSystem(GetSRDB(), system)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, string("Could not register system\n"))
			return
		}
		jsonRespStr, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, string("Could not register system\n"))
		}

		fmt.Println(string(jsonRespStr))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonRespStr)+"\n")
	}
}

///////////////////////////////////////////////////////////////////////////////
type metadata struct {
	Metadata map[string]string `json:"metadata"`
}

//
func HandleSystemById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var systemIDstr = vars["id"]
	if systemIDstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id32, err := strconv.Atoi(systemIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := int64(id32)
	fmt.Printf("\nHandleSystem(%d)\n", id)

	switch r.Method {
	case http.MethodGet:
		ret, err := getSystem(GetSRDB(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string("No such system\n"))
			return
		}

		/*if ret.SystemName == "" && ret.Id <= 0 { //xxx
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string("No such system\n"))
			return
		}*/

		jsonRespStr, err := json.Marshal(ret)

		fmt.Println(string(jsonRespStr))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonRespStr)+"\n")
		return

	case http.MethodPut:
		var system dto.SystemRequestDTO

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := r.Body.Close(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("%s\n", string(body))
		err = json.Unmarshal(body, &system)
		fmt.Printf("REPLSYS: %+v\n", system)

		ret, err := replaceSystem(GetSRDB(), id, system)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, string("Could not replace system\n"))
			return
		}

		jsonRespStr, err := json.Marshal(ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, string("Could not register system\n"))
			return
		}
		fmt.Println(string(jsonRespStr))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonRespStr)+"\n")
		return

	case http.MethodPatch:
		var system dto.SystemRequestDTO

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := r.Body.Close(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("%s\n", string(body))

		var msg map[string]interface{}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", msg)

		fmt.Printf("%+v\n", system)
		for k, x := range msg {
			if k == "metadata" {
				fmt.Printf("Found METADATA:\n%v\n", x)
				var systemMD dto.SystemRequestDTO

				md := make(map[string]string, 0)
				json.Unmarshal(body, &systemMD)
				if systemMD.Metadata != nil {
					for mdkey, mdval := range *systemMD.Metadata {
						fmt.Printf("[%s]: %s\n", mdkey, mdval)
						md[mdkey] = mdval
					}
				}

				system.Metadata = &md
				continue
			}
			switch v := x.(type) {
			case int:
				//fmt.Println("int:", v)
			case float64:
				if k == "port" {
					system.Port = int(v)
				}
			case string:
				if k == "systemName" {
					system.SystemName = v
				} else if k == "address" {
					system.Address = v
				} else if k == "authenticationInfo" {
					system.AuthenticationInfo = v
				}
			default:
				//fmt.Println("unknown")
			}
		}
		fmt.Printf("%+v\n", system)

		ret, err := modifySystem(GetSRDB(), id, system)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, string("Could not modify system\n"))
			return
		}

		jsonRespStr, err := json.Marshal(ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, string("Could not register system\n"))
			return
		}

		fmt.Println(string(jsonRespStr))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonRespStr)+"\n")
		return

	case http.MethodDelete:
		fmt.Printf("DELETE %v\n", id)
		err := deleteSystem(GetSRDB(), id)
		if err == nil {
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
}
