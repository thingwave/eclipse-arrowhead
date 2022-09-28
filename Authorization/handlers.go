package main

import (
  "log"
  "fmt"
  "errors"
  "io"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "strconv"
  "strings"

  "github.com/gorilla/mux"

  dto "arrowhead.eu/common/datamodels"
)

///////////////////////////////////////////////////////////////////////////////
// Client endpoint description

func Echo(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
  fmt.Fprint(w, "Got it!")
}

func GetPublicKey(w http.ResponseWriter, r *http.Request) {
  log.Println("GetPublicKey()")

  if true {
    cert := getCert()
    cert = strings.Replace(cert, "-----BEGIN PUBLIC KEY-----\n", "", 1)
    cert = strings.Replace(cert, "-----END PUBLIC KEY-----\n", "", 1)
    cert = strings.Replace(cert, "\n", "", 100)
    cert = "\"" + cert + "\""
    //w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, cert)
  } else {
	var errMsg dto.ErrorMessageDTO
	errMsg.ErrorMessage = "Authorization core service runs in insecure mode."
	errMsg.ErrorCode = 500
	errMsg.ExceptionType = "ARROWHEAD"
	errMsg.Origin = "/authorization/publickey"
	jsonRespStr, _ := json.Marshal(errMsg)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, string(jsonRespStr))
  }
}

///////////////////////////////////////////////////////////////////////////////
// Private endpoint description

func CheckIntraCloudRule(w http.ResponseWriter, r *http.Request) {

}

///////////////////////////////////////////////////////////////////////////////
// Management Endpoint Description

func HandleIntraCloudRule(w http.ResponseWriter, r *http.Request) {
  //var err error

  if r.Method == http.MethodGet {
    log.Println("Get all Intracloud rules")

    var ret AuthorizationIntraCloudListResponseDTO
    ret.Data, _ = GetAllIntraCloudRules(GetAUDB())
    ret.Count = len(ret.Data)

    jsonRespStr, _ := json.Marshal(ret)
    w.Header().Add("Content-Type", "application/json")
    fmt.Fprint(w, string(jsonRespStr))
    return
  } else if r.Method == http.MethodPost {
    log.Println("Add Intracloud rule")

    body, err := getBody(r)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    var request AuthorizationIntraCloudRequestDTO
    err = json.Unmarshal(body, &request)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    // validate and print
    fmt.Printf("%+v\n", request)
    rules, err := createBulkAuthorizationIntraCloudResponse(request.ConsumerId, request.ProviderIds, request.ServiceDefinitionIds, request.InterfaceIds)

    var ret AuthorizationIntraCloudListResponseDTO
    ret.Data = make([]AuthorizationIntraCloudResponseDTO, 0)
    for _, ruleId := range rules {
      ruleData, _ := GetIntraCloudRuleById(GetAUDB(), ruleId)
      ret.Data = append(ret.Data, ruleData)
    }

    jsonRespStr, _ := json.Marshal(ret)
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, string(jsonRespStr))

  }

}

func createBulkAuthorizationIntraCloudResponse(consumerId int64, providerIds []int64, serviceDefinitionIds []int64, interfaceIds []int64) ([]int64, error) {
  log.Println("createBulkAuthorizationIntraCloudResponse started...")
  ruleIds := make([]int64, 0)

  consumer, err := getSystem(GetAUDB(), consumerId)
  if err != nil {
    return ruleIds, errors.New(fmt.Sprintf("Consumer system with id of %d not exists"))
  }


  fmt.Printf("Creating Auth rule for ConsumerId: %v\n", consumer.Id)
  fmt.Printf("with Providers:\n")
  for _, pv := range(providerIds) {
    //fmt.Printf("  [%v]: %v\n", pk, pv)

    fmt.Printf("and ServiceDefinitionIds:\n")
    for _, sv := range(serviceDefinitionIds) {
      //fmt.Printf("    [%v]: %v\n", sk, sv)

      id, errInj := GetOrInsertAuthorizationRule(GetAUDB(), consumer.Id, pv, sv)
      if errInj != nil {
	      log.Println(errInj)
	      continue
      }

      if id != -1 {
	      fmt.Printf("RuleID to use: %v\n", id) //XXX IMPLEMENT RESPOSE RESULT
        ruleIds = append(ruleIds, id)
      } else {
	
      }
      /*
      
      fmt.Printf("using InterfaceIds:\n")
      for ik, iv := range(interfaceIds) {
	fmt.Printf("      [%v]: %v\n", ik, iv)

      }*/
    }
  }

  return ruleIds, nil
}

///////////

func HandleIntraCloudRuleByID(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)
  var ruleIDstr = vars["id"]
  if ruleIDstr == "" {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  ruleId32, err := strconv.Atoi(ruleIDstr)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  ruleId := int64(ruleId32)

  if ruleId == 0 {
    var errMsg dto.ErrorMessageDTO
    errMsg.ErrorMessage = "Id must be greater than 0."
    errMsg.ErrorCode = 400
    errMsg.ExceptionType = "BAD_PAYLOAD"
    errMsg.Origin = "/authorization/mgmt/intracloud/{id}"
    jsonRespStr, _ := json.Marshal(errMsg)
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, string(jsonRespStr))
    return
  }

  if r.Method == http.MethodGet {
    log.Printf("GetIntraCloudRuleByID(%v)\n", ruleId)
    
    ret, err := GetIntraCloudRuleById(GetAUDB(), ruleId)
    if err == nil {
      jsonRespStr, _ := json.Marshal(ret)
      w.Header().Add("Content-Type", "application/json")
      fmt.Fprint(w, string(jsonRespStr))
    } else {
      var errMsg ErrorMessageDTO
      errMsg.ErrorMessage = fmt.Sprintf("AuthorizationIntraCloud with id of '%v' not exists", ruleId)
      errMsg.ErrorCode = 400
      errMsg.ExceptionType = "INVALID_PARAMETER"
      errMsg.Origin = "/authorization/mgmt/intracloud/{id}"
      jsonRespStr, _ := json.Marshal(errMsg)
      w.Header().Add("Content-Type", "application/json")
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprint(w, string(jsonRespStr))
      return
    }
  } else if r.Method == http.MethodDelete {
    log.Printf("DeleteIntraCloudRuleByID(%v)\n", ruleId)
    ok, err := DeleteIntraCloudRuleById(GetAUDB(), ruleId)
    if ok != true || err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
  }

}


///////////////////////////////////////////////////////////////////////////////
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
