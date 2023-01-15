package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	dto "arrowhead.eu/common/datamodels"
)

var systemName string
var systemAddress string
var systemPort int
var srAddress string
var srPort int = 8443
var srSecure int = 0

func SetConfig(sysName string, sysAddress string, sysPort int, sr_address string, sr_port int, sr_secure int) {
	systemName = sysName
	systemAddress = sysAddress
	systemPort = sysPort
	srAddress = sr_address
	srPort = sr_port
	srSecure = sr_secure
}

func TestSRAvailability() (dto.ServiceRegistryResponseDTO, error) {
	log.Printf("TestSRAvailability()\n")

	var ret dto.ServiceRegistryResponseDTO

	return ret, nil
}

func RegisterService(systemname string, address string, port int, service_definition string, service_uri string, version int, interfaces []string) (bool, error) {
	log.Printf("RegisterService('%s', '%s', %d, '%s', '%s')\n", systemname, address, port, service_definition, service_uri)

	var req dto.ServiceRegistryEntryDTO
	req.ServiceDefinition = service_definition

	var provider dto.SystemRequestDTO
	provider.SystemName = systemname
	provider.Address = address
	provider.Port = port
	req.ProviderSystem = provider
	//req.ProviderSystem.AuthenticationInfo = systemAuthenticationInfo
	//req.ProviderSystem.Metadata = systemMetadata
	req.ServiceUri = service_uri
	req.EndOfValidity = "2023-12-31T23:59:59"
	req.Secure = "NOT_SECURE"
	if srSecure == 1 {
		req.Secure = "CERTIFICATE" // XXX TOKEN support
	}
	req.Version = version
	req.Interfaces = interfaces

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return false, err
	}
	fmt.Printf("%s\n", string(jsonReq))

	srUrl := "http://192.168.11.22:8443/serviceregistry/register" // fmt.Sprintf("%s://%s:%d/serviceregistry/register", mode, sr_address, sr_port)

	request, error := http.NewRequest("POST", srUrl, bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	return true, nil
}

func UnregisterService(systemname string, address string, port int, service_definition string, service_uri string) error {
	log.Printf("UnregisterService('%s', '%s', %d, '%s', '%s')\n", systemname, address, port, service_definition, service_uri)

	client := &http.Client{}
	url := fmt.Sprintf("http://192.168.11.22:8443/serviceregistry/unregister?service_definition=%s&system_name=%s&address=%s&port=%d", service_definition, systemname, address, port)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))

	return nil //errors.New("Not implemented")
}
