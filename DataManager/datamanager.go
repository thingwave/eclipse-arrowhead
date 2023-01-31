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
	//acl "accesscontrol"
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

	Sr_address string `toml:"sr_address"`
	Sr_port    int    `toml:"sr_port"`

	Websockets_enabled bool
	Acl_file           string

	Server_ssl_enabled              bool
	Server_ssl_client_auth          string
	Server_ssl_key_store            string
	Server_ssl_key_store_file       string
	Server_ssl_trust_store          string
	Server_ssl_trust_store_password string
}

var config tomlConfig

var t = 0
var authenticationInfo string = ""

func Timer() {

	for {
		time.Sleep(2 * time.Second)
		t++
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nCtrl-C pressed")
		util.UnregisterService(config.Core_system_name, config.Server_address, config.Server_port, "historian", "/datamanager/historian")
		util.UnregisterService(config.Core_system_name, config.Server_address, config.Server_port, "proxy", "/datamanager/proxy")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

}

func main() {
	fmt.Println("MySQL DataManager in Go, © ThingWave AB 2022")

	SetupCloseHandler()

	if _, cerr := toml.DecodeFile("application.toml", &config); cerr != nil {
		fmt.Println(cerr)
		return
	}

	fmt.Printf("server.address: %s\n", config.Server_address)
	fmt.Printf("server.port: %d\n", config.Server_port)
	fmt.Printf("core.system.name: %s\n", config.Core_system_name)
	fmt.Printf("sr_address: %s\n", config.Sr_address)
	fmt.Printf("sr_port: %d\n", config.Sr_port)
	fmt.Printf("server.ssl.enabled: %v\n", config.Server_ssl_enabled)

	db, err := db.OpenDatabase(config.Datasource_address, config.Datasource_port, config.Datasource_username, config.Datasource_password, config.Datasource_database)
	if err != nil {
		log.Fatal("Could not connect to database!")
	}
	defer db.Close()
	SetDMDB(db)

	aclErr := aclInit(config.Acl_file)
	if aclErr != nil {
		log.Printf("Could not load ACL file, aborting....\n")
		os.Exit(-1)
	}

	/*var secMode int = 0
	if config.Server_ssl_enabled {
		secMode = 1
	}
	*/
	util.SetConfig(config.Core_system_name, config.Server_address, config.Server_port, config.Sr_address, config.Sr_port, config.Server_ssl_enabled)

	//util.UnregisterService(config.Core_system_name, config.Server_address, config.Server_port, "historian", "/datamanager/historian")
	//util.UnregisterService(config.Core_system_name, config.Server_address, config.Server_port, "proxy", "/datamanager/proxy")

	interfaces := make([]string, 1)
	if config.Server_ssl_enabled {
		util.SetTLSConfig(config.Server_ssl_key_store, config.Server_ssl_key_store_file)
		authenticationInfo, err = util.SetAuthenticationInfo(config.Server_ssl_key_store)
		if err != nil {
			fmt.Println("Could not load system certificte public key!")
			return
		}
		interfaces[0] = "HTTP-SECURE-JSON"
	} else {
		interfaces[0] = "HTTP-INSECURE-JSON"
	}
	util.RegisterService(config.Core_system_name, config.Server_address, config.Server_port, "historian", "/datamanager/historian", 1, interfaces)
	util.RegisterService(config.Core_system_name, config.Server_address, config.Server_port, "proxy", "/datamanager/proxy", 1, interfaces)

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
		log.Println("Starting HTTP server")

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
