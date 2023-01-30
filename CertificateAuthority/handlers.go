package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
  "crypto/x509"
	"encoding/pem"
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
	//fmt.Printf("BODY: %s\n", body)

	return body, err
}

// 
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

  // add PEM header and footer
  certData := "-----BEGIN CERTIFICATE-----\n" + string(body) + "\n-----END CERTIFICATE-----\n"
  fmt.Println(certData)

  // decode cert and extract variable for the response
  block, _ := pem.Decode([]byte(certData))
	if block == nil {
		panic("failed to parse certificate PEM") //XXX
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error()) // XXX
	}

  status := "unknown"
  fmt.Printf("Version %s\n", cert.Version)
  fmt.Printf("Serial %s\n", cert.Subject.SerialNumber)
  //producedAt := "==="
  fmt.Printf("CommonName %s\n", cert.Subject.CommonName)
  fmt.Printf("Not before %s\n", cert.NotBefore.String())
  fmt.Printf("Not after %s\n", cert.NotAfter.String())
  fmt.Printf("Status %s\n", status)

}

func PrivSign(w http.ResponseWriter, r *http.Request) {
	var signReq CertificateSigningRequest
	//var signResp CertificateSigningResponse

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := getBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &signReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func PrivCheckTrustedKey(w http.ResponseWriter, r *http.Request) {
	var checkReq TrustedKeyCheckRequest
	//var checkResp TrustedKeyCheckResponse

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
}
