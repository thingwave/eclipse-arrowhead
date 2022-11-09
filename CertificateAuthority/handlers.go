package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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

func CheckCertificate(w http.ResponseWriter, r *http.Request) {
	var checkReq TrustedKeyCheckRequestDTO
	//var ret TrustedKeyCheckResponseDTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := getBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &checkReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if checkReq.Version != 1 || checkReq.PublicKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*



		fmt.Printf("################################\nQuery():\n %+v\n################################\n", queryReq)
		var unfilteredHits int = 0
		ret.ServiceQueryData, _ = queryServicesForName(GetSRDB(), queryReq, &unfilteredHits) //XX BUG HERE SOMEWHERE
		ret.UnfilteredHits = unfilteredHits
		retJson, _ := json.Marshal(ret)
		fmt.Fprint(w, string(retJson))

	*/
}
