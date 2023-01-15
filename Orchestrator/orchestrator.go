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
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	db "arrowhead.eu/common/database"
	util "arrowhead.eu/common/util"
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

var t = 0

func Timer() {

	for {
		time.Sleep(2 * time.Second)
		t++
		fmt.Printf("t: %d\n", t)
	}
}

//
//
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nCtrl-C pressed")
		util.UnregisterService("orchestrator", config.Server_address, config.Server_port, "orchestration-service", "/orchestration")
		//    time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

}

func main() {
	fmt.Println("MySQL Orchestrator in Go, © ThingWave AB 2022")

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

	db, err := db.OpenDatabase(config.Datasource_address, config.Datasource_port, config.Datasource_username, config.Datasource_password, config.Datasource_database) //"orchestrator", "KbgD2mTr8DQ4vtc", "arrowhead")
	if err != nil {
		log.Fatal("Could not connect to database!")
	}
	defer db.Close()
	SetORDB(db)

	//util.SetConfig(config.Core_system_name, config.Server_address, config.Server_port, config.Sr_address, config.Sr_port, secMode)

	// register all services
	util.RegisterService("orchestrator", config.Server_address, config.Server_port, "orchestration-service", "/orchestration", 1, []string{"HTTP-INSECURE-JSON"})

	router := NewRouter(config.Server_ssl_enabled)

	if config.Server_ssl_enabled {
		log.Printf("\nStarting HTTPS server\n")

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

		// Create a Server instance to listen on a port with the TLS config
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

		// Create a Server instance to listen on a port without TLS
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
