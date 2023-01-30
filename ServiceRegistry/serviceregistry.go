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
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"

	db "arrowhead.eu/common/database"
	dto "arrowhead.eu/common/datamodels"
	"arrowhead.eu/common/util"
)

type tomlConfig struct {
	Datasource_address  string
	Datasource_port     int
	Datasource_username string
	Datasource_password string
	Datasource_database string

	Server_address string
	Server_port    int

	Core_system_name string

	Server_ssl_enabled              bool
	Server_ssl_client_auth          string
	Server_ssl_key_store            string
	Server_ssl_key_store_file       string
	Server_ssl_trust_store          string
	Server_ssl_trust_store_password string
}

var config tomlConfig
var mySystemId int64
var authenticationInfo string = ""

func timer() {
	for {
		time.Sleep(1 * time.Second)
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nCtrl-C pressed\nUnregistering services...\n")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

}

func main() {
	log.Println("Eclipse Arrowhead ServiceRegistry in Go, © Lulea University of Technology AB 2022")

	SetupCloseHandler()

	if _, cerr := toml.DecodeFile("application.toml", &config); cerr != nil {
		fmt.Println(cerr)
		return
	}

	/* print the current configuration */
	fmt.Printf("server.address: %s\n", config.Server_address)
	fmt.Printf("server.port: %d\n", config.Server_port)
	fmt.Printf("core.system.name: %s\n", config.Core_system_name)
	fmt.Printf("server.ssl.enabled: %v\n", config.Server_ssl_enabled)
	fmt.Printf("Server.ssl.client.auth: %v\n", config.Server_ssl_client_auth)

	db, err := db.OpenDatabase(config.Datasource_address, config.Datasource_port, config.Datasource_username, config.Datasource_password, config.Datasource_database)
	if err != nil {
		log.Fatal("Could not connect to database!")
	}
	defer db.Close()
	SetSRDB(db)

	mySystemId = checkProvider(db, config.Core_system_name)
	if mySystemId == -1 {
		/*systemId, err = addSystem(db, "serviceregistry")
		if err != nil {
			return ret, err
		}*/
	}
	//fmt.Printf("My systemId is: %d\n", mySystemId)
	unregisterAllServicesForSystem(db, config.Core_system_name)

	var srvRegReq dto.ServiceRegistryEntryDTO
	srvRegReq.ServiceDefinition = "service-register"
	srvRegReq.ProviderSystem.SystemName = config.Core_system_name
	srvRegReq.ProviderSystem.Address = config.Server_address
	srvRegReq.ProviderSystem.Port = config.Server_port
	srvRegReq.ProviderSystem.AuthenticationInfo = ""
	srvRegReq.ProviderSystem.Metadata = nil

	addOrUpdateSystem(db, srvRegReq.ProviderSystem)

	srvRegReq.ServiceUri = "/serviceregistry/register"
	srvRegReq.EndOfValidity = ""
	srvRegReq.Secure = "NOT_SECURE"
	srvRegReq.Metadata = nil
	srvRegReq.Version = 1
	srvRegReq.Interfaces = make([]string, 1)
	srvRegReq.Interfaces[0] = "HTTP-INSECURE-JSON"
	registerServiceForSystem(db, srvRegReq)

	srvRegReq.ServiceDefinition = "service-register"
	srvRegReq.ServiceUri = "/serviceregistry/unregister"
	registerServiceForSystem(db, srvRegReq)

	srvRegReq.ServiceDefinition = "service-unregister"
	srvRegReq.ServiceUri = "/serviceregistry/unregister"
	registerServiceForSystem(db, srvRegReq)

	srvRegReq.ServiceDefinition = "register-system"
	srvRegReq.ServiceUri = "/serviceregistry/register-system"
	registerServiceForSystem(db, srvRegReq)

	srvRegReq.ServiceDefinition = "unregister-system"
	srvRegReq.ServiceUri = "/serviceregistry/unregister-system"
	registerServiceForSystem(db, srvRegReq)

	srvRegReq.ServiceDefinition = "pull-systems"
	srvRegReq.ServiceUri = "/serviceregistry/pull-systems"
	registerServiceForSystem(db, srvRegReq)

	router := NewRouter(config.Server_ssl_enabled)

	if config.Server_ssl_enabled {
		log.Printf("Starting HTTPS server\n")

		authenticationInfo, err = util.SetAuthenticationInfo(config.Server_ssl_key_store)
		if err != nil {
			fmt.Println("Could not load system certificte public key!")
			return
		}

		// Create a CA certificate pool and add cert.pem to it
		caCert, err := ioutil.ReadFile(config.Server_ssl_trust_store)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Create the TLS Config with the CA pool and enable Client certificate validation
		var authType tls.ClientAuthType = tls.NoClientCert
		if config.Server_ssl_client_auth == "need" {
			authType = tls.RequireAndVerifyClientCert
		}
		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: authType,
		}
		tlsConfig.BuildNameToCertificate()

		server := &http.Server{
			Addr:      config.Server_address + ":" + strconv.Itoa(config.Server_port),
			Handler:   router,
			TLSConfig: tlsConfig,
		}
		serr := server.ListenAndServeTLS(config.Server_ssl_key_store, config.Server_ssl_key_store_file)
		if serr != nil {
			log.Fatal("ListenAndServeTLS: ", serr)
		}

	} else {
		log.Printf("Starting HTTP server\n")
		server := &http.Server{
			Addr:    config.Server_address + ":" + strconv.Itoa(config.Server_port),
			Handler: router,
		}
		serr := server.ListenAndServe()
		if serr != nil {
			log.Fatal("ListenAndServe: ", serr)
		}
	}
}
