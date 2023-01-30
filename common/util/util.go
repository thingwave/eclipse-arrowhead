package util

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	dto "arrowhead.eu/common/datamodels"
)

var systemName string
var systemAddress string
var systemPort int
var srAddress string = "127.0.0.1"
var srPort int = 8443
var srSecure bool = false
var systemAuthenticationInfo = ""

func SetConfig(sysName string, sysAddress string, sysPort int, sr_address string, sr_port int, sr_secure bool) {
	log.Printf("SetConfig()\n")
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
	var client *http.Client
	mode := "http"

	var req dto.ServiceRegistryEntryDTO
	req.ServiceDefinition = service_definition

	var provider dto.SystemRequestDTO
	provider.SystemName = systemname
	provider.Address = address
	provider.Port = port
	req.ProviderSystem = provider
	req.ProviderSystem.AuthenticationInfo = systemAuthenticationInfo
	//req.ProviderSystem.Metadata = systemMetadata
	req.ServiceUri = service_uri
	req.EndOfValidity = "2023-12-31T23:59:59" //XXX dynamic lifetime
	req.Secure = "NOT_SECURE"
	if srSecure {
		req.Secure = "CERTIFICATE" // XXX TOKEN support
		mode = "https"
		caCert, err := ioutil.ReadFile("../certificates/testcloud2/testcloud2.pem")
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		req.ProviderSystem.AuthenticationInfo = ""
		client = &http.Client{}
	}
	req.Version = version
	req.Interfaces = interfaces

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return false, err
	}
	fmt.Printf("%s\n", string(jsonReq))

	srUrl := fmt.Sprintf("%s://%s:%d/serviceregistry/register", mode, srAddress, srPort)
	request, error := http.NewRequest("POST", srUrl, bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

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
	var client *http.Client
	mode := "http"

	if srSecure {
		caCert, err := ioutil.ReadFile("../certificates/testcloud2/testcloud2.pem")
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		mode = "https"
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}

	url := fmt.Sprintf("%s://%s:%d/serviceregistry/unregister?service_definition=%s&system_name=%s&address=%s&port=%d", mode, srAddress, srPort, service_definition, systemname, address, port)
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

	return nil
}

func SetAuthenticationInfo(fileName string) (string, error) {
	fmt.Printf("PEM2AuthenticationInfo(%s)\n", fileName)

	bytes, err := ioutil.ReadFile("../certificates/testcloud2/testcloud2.pem")
	if err != nil {
		return "", err
	}

	cert := string(bytes)
	//fmt.Printf("###\n%s\n###\n", cert)
	end1 := strings.Index(cert, "-----END CERTIFICATE-----")
	cert = cert[0:end1]
	cert = strings.Replace(cert, "-----BEGIN CERTIFICATE-----", "", 1)
	cert = strings.Replace(cert, "-----END CERTIFICATE-----", "", 1)
	cert = strings.ReplaceAll(cert, "\n", "")

	//fmt.Printf("###\n%s\n###\n", cert)
	systemAuthenticationInfo = cert
	return cert, nil
}
