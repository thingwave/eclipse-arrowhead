package util

import (
	"log"
	"fmt"
	"io/ioutil"
//	"errors"
//	"strings"
	"net/http"
	dto "arrowhead.eu/common/datamodels"
)

func SetConfig(systemName string, sr_address string, sr_port int) {

}

func TestSRAvailability() (dto.ServiceRegistryResponseDTO, error) {
	log.Printf("TestSRAvailability()\n")

	var ret dto.ServiceRegistryResponseDTO

	return ret, nil
}

func RegisterService(systemname string, address string, port int, service_definition string, service_uri string) (bool, error) {
	log.Printf("RegisterService()\n")

	return true, nil
}

func UnregisterService(systemname string, address string, port int, service_definition string, service_uri string) error {
    log.Printf("UnregisterService(''%s'', ''%s'', %d, '%s'', '%s'')\n", systemname, address, port, service_definition, service_uri)

	client := &http.Client{}
    req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8443/serviceregistry/unregister?system_name=datamanager", nil)
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

