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
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"database/sql"
	"encoding/json"

	dto "arrowhead.eu/common/datamodels"
	_ "github.com/go-sql-driver/mysql"
	//db "arrowhead.eu/common/database"
)

/*
	type Provider struct {
		ID int `json:"id"`
		//    System_id int `json:"system_id"`
		SystemName string `json:"systemName"`
		Address    string `json:"address"`
		//    Service_IP string `json:"service_ip"`
		Port       int    `json:"port"`
		CreatedAt  string `json:"createdAt"`
		UpdatedAt  string `json:"updatedAt"`
		ServiceUri string
	}
*/
type ServiceInterface struct {
	ID             int    `json:"id"`
	Interface_name string `json:"interface_name"`
}

var srdb *sql.DB = nil

///////////////////////////////////////////////////////////////////////////////
//
//
/*func getAllServices(db *sql.DB, serviceType string) bool {
	var query = "SELECT\n" +
		"  service_definition.service_definition,\n" +
		"  service_registry.id as serviceId,\n" +
		"  service_registry.service_uri,\n" +
		"  system_.id as systemId,\n" +
		"  system_.system_name\n" +
		"FROM service_registry\n" +
		"INNER JOIN system_\n" +
		"ON service_registry.system_id = system_.id\n" +
		"INNER JOIN service_definition\n" +
		"ON service_registry.service_id = service_definition.id\n" +
		"ORDER BY system_.system_name;\n"

	fmt.Println(query)

	return false
}*/

func SetSRDB(db *sql.DB) {
	srdb = db
}

func GetSRDB() *sql.DB {
	return srdb
}

// /////////////////////////////////////////////////////////////////////////////
func getInterfacesForService(db *sql.DB, serviceID int64) ([]dto.ServiceInterfaceResponseDTO, error) {
	//fmt.Printf("getInterfacesForService(%v)\n", serviceID)
	ret := make([]dto.ServiceInterfaceResponseDTO, 0)

	results, err := db.Query("SELECT DISTINCT si.id, si.interface_name, UNIX_TIMESTAMP(si.created_at), UNIX_TIMESTAMP(si.updated_at) FROM service_interface si, service_registry_interface_connection sric WHERE si.id=sric.interface_id AND sric.service_registry_id=?", serviceID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return ret, err
	}
	defer results.Close()

	for results.Next() {
		var interfaceEntry dto.ServiceInterfaceResponseDTO
		var created_at, updated_at string
		err = results.Scan(&interfaceEntry.Id, &interfaceEntry.InterfaceName, &created_at, &updated_at)
		if err != nil {
			fmt.Println(err)
			//panic(err.Error()) // proper error handling instead of panic in your app
			continue
		} else {
			interfaceEntry.CreatedAt = timestamp2Arrowhead(created_at)
			interfaceEntry.UpdatedAt = timestamp2Arrowhead(updated_at)
			//fmt.Println("\tInterface name: ", interfaceEntry.InterfaceName)
			fmt.Printf("%+v\n", interfaceEntry)
			ret = append(ret, interfaceEntry)
		}
	}

	return ret, nil
}

func getInterfacesList(db *sql.DB) ([]dto.ServiceInterfaceResponseDTO, error) {
	//fmt.Printf("getInterfacesList()\n")

	ret := make([]dto.ServiceInterfaceResponseDTO, 0)

	results, err := db.Query("SELECT DISTINCT id, interface_name, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_interface")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return ret, err
	}
	defer results.Close()

	for results.Next() {
		var interfaceEntry dto.ServiceInterfaceResponseDTO
		var created_at, updated_at string
		err = results.Scan(&interfaceEntry.Id, &interfaceEntry.InterfaceName, &created_at, &updated_at)
		if err != nil {
			return ret, err
		} else {
			interfaceEntry.CreatedAt = timestamp2Arrowhead(created_at)
			interfaceEntry.UpdatedAt = timestamp2Arrowhead(updated_at)
			ret = append(ret, interfaceEntry)
		}
	}

	return ret, nil
}

func getInterfaceById(db *sql.DB, interfaceId int64) (dto.ServiceInterfaceResponseDTO, error) {
	//fmt.Printf("getInterfaceById('%d')\n", interfaceId)

	var ret dto.ServiceInterfaceResponseDTO
	ret.Id = -1

	results, err := db.Query("SELECT id, interface_name, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_interface WHERE id=?", interfaceId)
	if err != nil {
		//panic(err.Error()) // proper error handling instead of panic in your app
		return ret, err
	}
	defer results.Close()

	if results.Next() {
		var created_at, updated_at string
		err = results.Scan(&ret.Id, &ret.InterfaceName, &created_at, &updated_at)
		ret.CreatedAt = timestamp2Arrowhead(created_at)
		ret.UpdatedAt = timestamp2Arrowhead(updated_at)
		return ret, nil
	}

	return ret, errors.New("No such interface name")
}

func updateInterfaceById(db *sql.DB, interfaceId int64, interfaceName string) error {
	//fmt.Printf("updateInterfaceById('%d')\n", interfaceId)

	res, err := db.Exec("UPDATE service_interface SET interface_name=? WHERE id=?", interfaceName, interfaceId)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	ra, err := res.RowsAffected()
	if ra != 1 || err != nil {
		fmt.Printf("%v\n", err)
		return errors.New("No rows updated")
	}

	return nil
}

func deleteInterfaceById(db *sql.DB, interfaceId int64) error {
	//fmt.Printf("deleteInterfaceById('%d')\n", interfaceId)

	res, err := db.Exec("DELETE FROM service_interface WHERE id=?", interfaceId)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	ra, err := res.RowsAffected()
	if ra != 1 || err != nil {
		fmt.Printf("%v\n", err)
		return errors.New("No rows deleted")
	}

	return nil
}

func getInterfaceByName(db *sql.DB, interfaceName string) (dto.ServiceInterfaceResponseDTO, error) {
	//fmt.Printf("getInterfaceByName('%s')\n", interfaceName)

	var ret dto.ServiceInterfaceResponseDTO
	ret.Id = -1

	results, err := db.Query("SELECT id, interface_name, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_interface WHERE interface_name=?", interfaceName)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return ret, err
	}
	defer results.Close()

	if results.Next() {
		var created_at, updated_at string
		err = results.Scan(&ret.Id, &ret.InterfaceName, &created_at, &updated_at)
		ret.CreatedAt = timestamp2Arrowhead(created_at)
		ret.UpdatedAt = timestamp2Arrowhead(updated_at)
		return ret, nil
	}

	return ret, errors.New("No such interface name")
}

func addInterfaceByName(db *sql.DB, interfaceName string) (dto.ServiceInterfaceResponseDTO, error) {

	return getInterfaceByName(db, interfaceName)
}

// /////////////////////////////////////////////////////////////////////////////
func getAllServicesBySystem(db *sql.DB, systemID int64) ([]dto.ServiceRegistryResponseDTO, error) {
	//fmt.Printf("getAllServicesBySystem(%v)\n", systemID)

	var ret []dto.ServiceRegistryResponseDTO
	ret = make([]dto.ServiceRegistryResponseDTO, 0)

	results, err := db.Query("SELECT id, service_id, service_uri, end_of_validity, secure, version, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_registry WHERE system_id=?;", systemID)
	if err != nil {
		return ret, err
	}
	defer results.Close()

	for results.Next() {
		var serviceDefID int
		var end_of_validity, metadata sql.NullString
		var created_at, updated_at string
		var serviceEntry dto.ServiceRegistryResponseDTO
		err = results.Scan(&serviceEntry.Id, &serviceDefID, &serviceEntry.ServiceUri, &end_of_validity, &serviceEntry.Secure, &serviceEntry.Version, &created_at, &updated_at)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if end_of_validity.Valid {
			serviceEntry.EndOfValidity = end_of_validity.String
		}
		if metadata.Valid {
			//serviceEntry.Metadata =
		}

		serviceEntry.ServiceDefinition, _ = getServiceDefinitionForService(db, int64(serviceDefID))
		serviceEntry.Provider, _ = getSystem(db, systemID)
		serviceEntry.Interfaces, _ = getInterfacesForService(db, serviceEntry.Id)
		serviceEntry.CreatedAt = timestamp2Arrowhead(created_at)
		serviceEntry.UpdatedAt = timestamp2Arrowhead(updated_at)

		ret = append(ret, serviceEntry)
	}

	return ret, nil
}

// /////////////////////////////////////////////////////////////////////////////
func getServiceDefinitionForService(db *sql.DB, serviceID int64) (dto.ServiceDefinitionResponseDTO, error) {
	//fmt.Printf("getServiceDefinitionForService(%v)\n", serviceID)

	result, err := db.Query("SELECT id, service_definition, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_definition WHERE id=? LIMIT 1;", serviceID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer result.Close()

	ret := dto.ServiceDefinitionResponseDTO{}
	if result.Next() {
		var service_definition sql.NullString
		var created_at, updated_at string

		err = result.Scan(&ret.Id, &service_definition, &created_at, &updated_at)
		if err != nil {
			//fmt.Printf("getServiceDefinitionForService 308: %v\n", err)
			return ret, err
		}

		if service_definition.Valid {
			ret.ServiceDefinition = service_definition.String
		}

		ret.CreatedAt = timestamp2Arrowhead(created_at)
		ret.UpdatedAt = timestamp2Arrowhead(updated_at)

		//fmt.Printf("%+v\n", ret)
		return ret, nil
	} else {
		fmt.Printf("No data for service_id %v!!!\n", serviceID)
	}
	return ret, errors.New("No data")
}

func getServiceDefinitionFromName(db *sql.DB, serviceDefinition string) (dto.ServiceDefinitionResponseDTO, error) {
	//fmt.Printf("getServiceDefinitionFromName(%s)\n", serviceDefinition)
	ret := dto.ServiceDefinitionResponseDTO{}

	result, err := db.Query("SELECT id, service_definition, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_definition WHERE service_definition=? LIMIT 1", serviceDefinition)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return ret, err
	}
	defer result.Close()

	if result.Next() {
		var service_definition sql.NullString
		var created_at, updated_at string

		err = result.Scan(&ret.Id, &service_definition, &created_at, &updated_at)
		if err != nil {
			return ret, err
		}

		if service_definition.Valid {
			ret.ServiceDefinition = service_definition.String
		}

		ret.CreatedAt = timestamp2Arrowhead(created_at)
		ret.UpdatedAt = timestamp2Arrowhead(updated_at)

		fmt.Printf("%+v\n", ret)
		return ret, nil
	} else {
		//fmt.Printf("No data for serviceDefinition %s!\n", serviceDefinition)
	}
	return ret, errors.New("No data")
}

// /////////////////////////////////////////////////////////////////////////////
func queryServicesForName(db *sql.DB, request ServiceQueryForm, unfilteredHits *int) ([]dto.ServiceRegistryResponseDTO, error) {
	returnList := []dto.ServiceRegistryResponseDTO{}
	var serviceType string = request.ServiceDefinitionRequirement
	fmt.Printf("queryServicesForName('%s')\n", serviceType)
	*unfilteredHits = 0

	result, err := db.Query("SELECT id, service_definition, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_definition WHERE service_definition=? LIMIT 1", serviceType)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer result.Close()

	if !result.Next() {
		fmt.Printf("No data...")
		return returnList, nil
	}

	var id int64
	var service_definition string
	var created_at, updated_at string

	err = result.Scan(&id, &service_definition, &created_at, &updated_at)
	if err != nil {
		log.Printf("Scan went wrong!\n")
		log.Println(err)
		return returnList, err
	}
	fmt.Printf("id: %d, service_definition: '%s', created_at: '%s', updated_at: '%s'\n", id, service_definition, created_at, updated_at)

	var serviceDef dto.ServiceDefinitionResponseDTO
	serviceDef.Id = id
	serviceDef.ServiceDefinition = service_definition
	serviceDef.CreatedAt = timestamp2Arrowhead(created_at)
	serviceDef.UpdatedAt = timestamp2Arrowhead(updated_at)
	fmt.Printf("serviceDef:\n%v\n", serviceDef)

	results, err2 := db.Query("SELECT id, system_id, service_uri, end_of_validity, secure, metadata, version, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_registry WHERE service_id=?", id)
	if err2 != nil {
		log.Println(err2)
		return returnList, err2
	}
	defer results.Close()
	for results.Next() {
		var srid int
		var entry dto.ServiceRegistryResponseDTO
		entry.ServiceDefinition = serviceDef
		var end_of_validity, metadata sql.NullString
		//var created_at, updated_at string
		err = results.Scan(&srid, &entry.Id, &entry.ServiceUri, &end_of_validity, &entry.Secure, &metadata, &entry.Version, &created_at, &updated_at)
		if err != nil {
			log.Println(err)
			continue
		}
		if end_of_validity.Valid {
			entry.EndOfValidity = end_of_validity.String
		}
		entry.CreatedAt = timestamp2Arrowhead(created_at)
		entry.UpdatedAt = timestamp2Arrowhead(updated_at)

		fmt.Printf("entry:\n%+v\n", entry)
		//fmt.Printf("len(%d)\n", len(request.SecurityRequirements))

		if len(request.SecurityRequirements) >= 1 {
			//fmt.Println("Filtering on SecurityRequirements")

			var validSecureModeFound = false
			for _, sv := range request.SecurityRequirements {
				if entry.Secure == sv {
					//log.Printf("Found matching security requirement: '%s'\n", sv)
					validSecureModeFound = true
				}
			}

			if validSecureModeFound == false {
				//log.Printf("Wrong security mode, discarding...\n")
				*unfilteredHits += 1
				continue
			}
		}

		if request.VersionRequirement != nil {
			//fmt.Println("Filtering on VersionRequirement")
			if *request.VersionRequirement != entry.Version {
				//log.Printf("Wrong version, discarding...\n")
				*unfilteredHits += 1
				continue
			}
		}

		// for each row, scan the result into our composite object
		//    var system_id int
		var sys dto.SystemResponseDTO
		sys.Id = entry.Id
		/*var serviceUri string
		  err = results.Scan(&system_id, &sys.Id, &serviceUri)
		  if err != nil {
		    panic(err.Error()) // proper error handling instead of panic in your app
		  }*/

		sysResult, sysErr := db.Query("SELECT id, system_name, address, port, authentication_info, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM system_ WHERE id=? LIMIT 1", entry.Id)
		if sysErr != nil {
			log.Println(sysErr)
			continue
		}
		defer sysResult.Close()

		if sysResult.Next() {
			var authentication_info sql.NullString
			err = sysResult.Scan(&sys.Id, &sys.SystemName, &sys.Address, &sys.Port, &authentication_info, &created_at, &updated_at)
			if err != nil {
				log.Println(err)
				continue
			}
			if authentication_info.Valid {
				sys.AuthenticationInfo = authentication_info.String
			}
			sys.CreatedAt = timestamp2Arrowhead(created_at)
			sys.UpdatedAt = timestamp2Arrowhead(updated_at)

			entry.Provider = sys
			entry.Interfaces, _ = getInterfacesForService(db, int64(srid))

			// add filter for serviceRequirements, and only add if match!
			var tmpInterfaces = make([]dto.ServiceInterfaceResponseDTO, 0)
			var validInterfaceFound bool = true
			if len(request.InterfaceRequirements) > 0 {
				validInterfaceFound = false
			}

			// filter
			for k, v := range entry.Interfaces {
				//log.Printf("\t[%d]: %s\n", k, v.InterfaceName)
				for _, rv := range request.InterfaceRequirements {
					if v.InterfaceName == rv {
						tmpInterfaces = append(tmpInterfaces, entry.Interfaces[k])
						validInterfaceFound = true
						*unfilteredHits += 1
					}
				}
			}
			if len(request.InterfaceRequirements) > 0 {
				entry.Interfaces = tmpInterfaces
			}
			if validInterfaceFound {
				returnList = append(returnList, entry)
			}
		}
	}

	return returnList, nil
}

// /////////////////////////////////////////////////////////////////////////////
func registerServiceForSystem(db *sql.DB, serviceRegReq dto.ServiceRegistryEntryDTO) (dto.ServiceRegistryResponseDTO, error) {
	var ret dto.ServiceRegistryResponseDTO
	var err error

	//fmt.Printf("\nregisterServiceForSystem('%s', '%s')\n", serviceRegReq.ServiceDefinition, serviceRegReq.ProviderSystem.SystemName)

	// check that a provider system exist
	systemId := checkProvider(db, serviceRegReq.ProviderSystem.SystemName)
	//fmt.Printf("Got systemId=%d\n", systemId)
	if systemId == -1 {
		//return ret, fmt.Errorf("No provider found, adding it...")
		systemId, err = addSystem(db, serviceRegReq.ProviderSystem)
		if err != nil {
			return ret, fmt.Errorf("No provider found, and could not add new provider")
		}
	} else {
		_ = updateProvider(db, serviceRegReq.ProviderSystem, systemId)
	}

	// validate service definition
	serviceDefId, err := insertOrUpdateServiceDefinition(db, serviceRegReq.ServiceDefinition) //checkAndRegisterServiceDefinition(db, serviceRegReq.ServiceDefinition)
	if err != nil {
		return ret, errors.New("Could not save new serviceDefinition type!")
	}
	//fmt.Printf("\tUsing serviceDefinitionID: %d\n", serviceDefId)

	// register the service
	serviceId, err := insertServiceEntry(db, serviceRegReq, systemId, serviceDefId)
	if err != nil && serviceId == -1 {
		log.Printf("%s:%d -> %v\n", "srdatabase.go", 542, err)
		return ret, err
	}
	//fmt.Printf("Got service.ID=%d\n", serviceId)

	service, _ := fetchServiceById(db, serviceId)

	//XXX: fill in ALL fields here
	ret.Id = serviceId
	//ret.ServiceDefinition.Id = serviceDefId
	ret.ServiceDefinition, _ = getServiceDefinitionForService(db, serviceDefId)
	provider, err := getSystem(db, systemId)
	ret.Provider = provider
	ret.Secure = service.Secure
	ret.ServiceUri = service.ServiceUri
	ret.Interfaces, _ = getInterfacesForService(db, serviceId)
	ret.Version = service.Version
	ret.CreatedAt = service.CreatedAt
	ret.UpdatedAt = service.UpdatedAt

	return ret, nil
}

/*
	func insertOrUpdateServiceDefinition(db *sql.DB, serviceDefinition string) (int64, error) {
		log.Printf("insertOrUpdateServiceDefinition('%s')\n", serviceDefinition)

		result, err := db.Exec("INSERT IGNORE INTO service_definition(service_definition) VALUES(?) ON DUPLICATE KEY UPDATE updated_at=NOW()", serviceDefinition)
		if err != nil {
			return -1, err
		}
		insertID, err := result.LastInsertId()
		if err != nil {
			log.Println(err)
			return -1, err
		}
		log.Printf("\tserviceDefinitionID: %v\n", insertID)

		return int64(insertID), nil
	}
*/
func insertOrUpdateServiceInterface(db *sql.DB, serviceInterface string) (int64, error) {
	//log.Printf("insertOrUpdateServiceInterface('%s')\n", serviceInterface)

	result, err := db.Exec("INSERT IGNORE INTO service_interface(interface_name) VALUES(?) ON DUPLICATE KEY UPDATE updated_at=NOW()", serviceInterface)
	if err != nil {
		return -1, err
	}
	insertID, err := result.LastInsertId()
	//log.Printf("\tserviceInterfaceID: %v\n", insertID)
	if err != nil {
		return -1, err
	}

	return int64(insertID), nil
}

func insertInterfaceConnection(db *sql.DB, serviceRegistryId int64, interfaceId int64) (int64, error) {
	//log.Printf("insertInterfaceConnection(%v,  %v)\n", serviceRegistryId, interfaceId)

	result, err := db.Exec("INSERT IGNORE INTO service_registry_interface_connection(service_registry_id, interface_id) VALUES(?, ?) ON DUPLICATE KEY UPDATE updated_at=NOW()", serviceRegistryId, interfaceId)
	if err != nil {
		return -1, err
	}
	insertID, err := result.LastInsertId()
	//log.Printf("\tserviceInterfaceConnectionID: %v\n", insertID)
	if err != nil {
		return -1, err
	}

	return int64(insertID), nil
}

///////////////////////////////////////////////////////////////////////////////
//

func updateServiceEntry(db *sql.DB, request dto.ServiceRegistryEntryDTO, ServiceId int64) (bool, error) { //XXX <----------------------------------------------------------------------------------------
	//fmt.Printf("updateServiceEntry()")

	var end_of_validity sql.NullString
	if request.EndOfValidity != "" {
		end_of_validity.String = request.EndOfValidity
		end_of_validity.Valid = true
	}

	result, err := db.Exec("UPDATE service_registry SET service_uri=?, end_of_validity=?, secure=?, version=?, updated_at=NOW() WHERE id=?", request.ServiceUri, end_of_validity, request.Secure, request.Version, ServiceId)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	ra, err := result.RowsAffected()
	//fmt.Printf("ra: %d\n", ra)
	//fmt.Println(err)
	if ra != 1 || err != nil {
		return false, errors.New("No rows updated")
	}

	return true, nil
}

func insertServiceEntry(db *sql.DB, request dto.ServiceRegistryEntryDTO, systemId int64, serviceDefId int64) (int64, error) {
	var id int64 = -1
	//var retId int64 = -1

	//tx, err := db.Begin()
	//defer tx.Rollback()

	result, err := db.Query("SELECT id FROM service_registry WHERE service_id=? AND system_id=? AND service_uri=? LIMIT 1", serviceDefId, systemId, request.ServiceUri) //systemId, serviceDefId, request.ServiceUri)
	if err != nil {
		//panic(err.Error()) // proper error handling instead of panic in your app
		return -1, errors.New("Database error: ")
	}
	defer result.Close()

	if result.Next() {
		result.Scan(&id)
		return id, errors.New("ServiceEntry already exists!")
	}

	//fmt.Printf("No matching service ID found, must insert...\n")

	/*endofValidityPtr := &request.EndOfValidity
	if request.EndOfValidity == "" {
		endofValidityPtr = nil
	}*/
	endofValidityPtr := NewNullString(request.EndOfValidity)
	insertRes, err := db.Exec("INSERT INTO service_registry(service_id, system_id, service_uri, end_of_validity, secure, version) VALUES(?,?,?,?,?,?)", serviceDefId, systemId, request.ServiceUri, endofValidityPtr, request.Secure, request.Version) //metadata must be converted to JSON
	if err != nil {
		//fmt.Printf("%s:%d->%v\n", "srdatabase.go", 671, err.Error())
		return -1, err
	}

	id, err = insertRes.LastInsertId()
	if err != nil {
		//fmt.Printf("%s:%d->%v\n", "srdatabase.go", 677, err.Error())
		return -1, err
	}

	// update/insert into service_interface and service_registry_interface_connection below
	for _, v := range request.Interfaces { //k
		//fmt.Printf("checkAndAdd[%v]: %s\n", k, v)
		interfaceId, _ := insertOrUpdateServiceInterface(db, v)
		insertInterfaceConnection(db, id, interfaceId)
	}

	return id, nil
}

// /////////////////////////////////////////////////////////////////////////////
func checkProvider(db *sql.DB, systemName string) int64 {
	var id int = -1

	//log.Printf("checkProvider(%s)\n", systemName)

	result, err := db.Query("SELECT id FROM system_ WHERE system_name=? LIMIT 1", systemName)
	if err != nil {
		return -1
	}
	defer result.Close()

	result.Next()
	err = result.Scan(&id)
	if err != nil {
		return -1
	}

	return int64(id)
}

func updateProvider(db *sql.DB, system dto.SystemRequestDTO, systemId int64) error {
	//var systemName string = system.SystemName
	var address string = system.Address
	var port int = system.Port
	//var authenticationInfo = newNullString(system.AuthenticationInfo)

	//fmt.Printf("updateProvider('%s', '%s', %d):\n", systemName, address, port)

	_, err := db.Exec("UPDATE system_ (address, port, updated_at) VALUES(?, ?, NOW()) WHERE id=?", address, port, systemId) //XX metadata is NULL?
	if err != nil {
		return err
	}

	return nil
}

func unregisterAllServicesForSystem(db *sql.DB, systemName string) {

	systemId := checkProvider(db, systemName)
	if systemId == -1 {
		//fmt.Printf("unregisterAllServicesForSystem: no such system %v!\n", systemId)
		return
	}

	_, err := db.Exec("DELETE FROM service_registry WHERE system_id = ?", systemId)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
//
/*
func getProvider(db *sql.DB, systemID int) (SystemResponseDTO, error) {
	var ret SystemResponseDTO
	ret.Id = systemID

	fmt.Printf("getProvider(%v): :\n", systemID)

	//result, err := db.Query("SELECT id, system_name, address, port, authentication_info, metadata, UNIX_TIMESTAMP(created_at) FROM system_ WHERE id=? LIMIT 1", systemID)
	result, err := db.Query("SELECT id, system_name, address, port, authentication_info, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM system_ WHERE id=? LIMIT 1;", systemID)
	if err != nil {
		return ret, err
	}
	defer result.Close()

	if result.Next() {
		var authentication_info sql.NullString
		var createdAt, updatedAt string
		//err = result.Scan(&ret.Id, &ret.SystemName, &ret.Address, &ret.Port, &ret.AuthenticationInfo, &metadata, &createdAt)
		err = result.Scan(&ret.Id, &ret.SystemName, &ret.Address, &ret.Port, &authentication_info, &createdAt, &updatedAt)

		fmt.Printf("createdAt: %v\n", createdAt)
		intCa, _ := strconv.Atoi(createdAt)
		ca := time.Unix(int64(intCa), 0)
		fmt.Printf("updatedAt: %v\n", updatedAt)
		intUa, _ := strconv.Atoi(updatedAt)
		ua := time.Unix(int64(intUa), 0)
		//fmt.Printf("ca: %v, ua:%v\n", ca, ua)
		//ret.CreatedAt = createdAt.In(time.UTC).Format(time.RFC3339) //createdAt.Format(time.RFC3339)
		ret.CreatedAt = ca.UTC().Format(time.RFC3339)
		ret.UpdatedAt = ua.UTC().Format(time.RFC3339) //"2006-01-02 15:04:05")
	} else {
		fmt.Println("Could not fetch entry")
		return ret, errors.New("Database error")
	}

	fmt.Printf("%+v\n", ret)

	return ret, nil
}
*/

// /////////////////////////////////////////////////////////////////////////////
func insertOrUpdateServiceDefinition(db *sql.DB, serviceDefiniton string) (int64, error) {

	var id int = -1
	result, err := db.Query("SELECT id FROM service_definition WHERE service_definition=? LIMIT 1", serviceDefiniton)
	if err != nil {
		return -1, err
	}
	defer result.Close()

	if !result.Next() {

		db.Exec("INSERT INTO service_definition(service_definition) VALUES(?)", serviceDefiniton)
		result, err = db.Query("SELECT id FROM service_definition WHERE service_definition=? LIMIT 1", serviceDefiniton)
		defer result.Close()
		result.Next()
		err = result.Scan(&id)
	} else {
		err = result.Scan(&id)
	}

	return int64(id), err
}

// /////////////////////////////////////////////////////////////////////////////
func getAllServiceDefinitions(db *sql.DB) ([]dto.ServiceDefinitionResponseDTO, error) {
	var ret []dto.ServiceDefinitionResponseDTO = make([]dto.ServiceDefinitionResponseDTO, 0)
	result, err := db.Query("SELECT id, service_definition, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_definition")
	if err != nil {
		return ret, err
	}
	defer result.Close()

	for result.Next() {
		var serviceDef dto.ServiceDefinitionResponseDTO
		var created_at, updated_at string
		result.Scan(&serviceDef.Id, &serviceDef.ServiceDefinition, &created_at, &updated_at)
		serviceDef.CreatedAt = timestamp2Arrowhead(created_at)
		serviceDef.UpdatedAt = timestamp2Arrowhead(updated_at)

		ret = append(ret, serviceDef)
	}

	return ret, nil
}

// /////////////////////////////////////////////////////////////////////////////
func addOrUpdateSystem(db *sql.DB, system dto.SystemRequestDTO) (int64, error) {
	var systemName string = system.SystemName
	var address string = system.Address
	var port int = system.Port
	var authenticationInfo = newNullString(system.AuthenticationInfo)

	//log.Printf("addOrUpdateSystem('%s', '%s', %d):\n", systemName, address, port)

	result, err := db.Exec("INSERT INTO system_ (system_name, address, port, authentication_info, metadata) VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE updated_at=NOW()", systemName, address, port, authenticationInfo, nil) //XX metadata is NULL
	if err != nil {
		return -1, err
	}
	insertID, err := result.LastInsertId()
	//log.Printf("insertID: %v\n", insertID)
	if err != nil {
		return -1, err
	}

	return insertID, nil
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// /////////////////////////////////////////////////////////////////////////////
func getIDforSystem(db *sql.DB, systemName string, address string, port int) int64 {
	result, err := db.Query("SELECT id FROM system_ WHERE system_name=? AND address=? AND port=? LIMIT 1", systemName, address, port)
	if err != nil {
		return -1
	}
	defer result.Close()

	var id int64 = -1
	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return -1
		}
		//fmt.Printf("System name (name, id): %s/%d\n", systemName, id)
	}

	return id
}

// /////////////////////////////////////////////////////////////////////////////
func getIDforServiceDefinition(db *sql.DB, serviceDefinition string) int64 {
	result, err := db.Query("SELECT id FROM service_definition WHERE service_definition=? LIMIT 1", serviceDefinition)
	if err != nil {
		return -1
	}
	defer result.Close()

	result.Next()
	var id int64 = -1
	var id32 int
	err = result.Scan(&id32)
	if err != nil {
		//fmt.Printf("getIDforServiceDefinition 404 - %v\n", err)
		return -1
	}
	id = int64(id32)
	//fmt.Printf("System name (serviceDefinition,id): %s, %d\n", serviceDefinition, id)

	return id
}

// /////////////////////////////////////////////////////////////////////////////
func registerSystem(db *sql.DB, system dto.SystemRequestDTO) (dto.SystemResponseDTO, error) {
	var ret dto.SystemResponseDTO

	systemID, err := addOrUpdateSystem(db, system)
	//systemID := getIDforSystem(db, system.SystemName, system.Address, system.Port)
	//fmt.Printf("SystemID: %d\n", systemID)

	if err != nil || systemID == -1 {
		return ret, fmt.Errorf("Could not register or update system!")
	}

	// get updated information
	ret, err = getSystem(db, systemID)
	return ret, err
}

// /////////////////////////////////////////////////////////////////////////////
func getSystem(db *sql.DB, systemId int64) (dto.SystemResponseDTO, error) {
	var ret dto.SystemResponseDTO

	//log.Printf("getSystem(%d)\n", systemId)

	result, err := db.Query("SELECT id, system_name, address, port, authentication_info, metadata, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM system_ WHERE id=? LIMIT 1;", systemId)
	if err != nil {
		//log.Printf("## A")
		log.Println(err)
		return ret, err
	}
	defer result.Close()

	result.Next()
	var authentication_info, metadata sql.NullString
	var created_at, updated_at string
	err = result.Scan(&ret.Id, &ret.SystemName, &ret.Address, &ret.Port, &authentication_info, &metadata, &created_at, &updated_at)
	if err != nil {
		//log.Printf("## B")
		fmt.Println(err)
	}
	if authentication_info.Valid {
		ret.AuthenticationInfo = authentication_info.String
	}
	if metadata.Valid {
		ret.Metadata = nil //XXX func to convert metadata to map
	}

	ret.CreatedAt = timestamp2Arrowhead(created_at)
	ret.UpdatedAt = timestamp2Arrowhead(updated_at)

	//log.Printf("## Z")
	return ret, err
}

// /////////////////////////////////////////////////////////////////////////////
func replaceSystem(db *sql.DB, id int64, system dto.SystemRequestDTO) (dto.SystemResponseDTO, error) {
	ret, err := getSystem(db, id)

	//fmt.Printf("Replace DATA: %+v\n", system)

	if system.SystemName != "" {
		_, err := db.Exec("UPDATE system_ SET system_name=? WHERE id =?;", system.SystemName, id)
		if err != nil {
			return ret, err
		}
	}
	if system.Address != "" {
		_, err := db.Exec("UPDATE system_ SET address=? WHERE id =?;", system.Address, id)
		if err != nil {
			return ret, err
		}
	}
	if system.Port != 0 {
		_, err := db.Exec("UPDATE system_ SET port=? WHERE id =?;", system.Port, id)
		if err != nil {
			return ret, err
		}
	}
	if system.AuthenticationInfo != "" {
		_, err := db.Exec("UPDATE system_ SET authentication_info=? WHERE id =?;", system.AuthenticationInfo, id)
		if err != nil {
			return ret, err
		}
	}

	ret, err = getSystem(db, id)
	return ret, err
}

// /////////////////////////////////////////////////////////////////////////////
func modifySystem(db *sql.DB, id int64, system dto.SystemRequestDTO) (dto.SystemResponseDTO, error) {
	var ret dto.SystemResponseDTO
	//  var err

	/*systemID := checkProvider(db, system.SystemName)
	if systemID == -1 {
		fmt.Printf("No such system\n")
		return ret, fmt.Errorf("No such system")
	}*/

	//fmt.Printf("DATA: %+v\n", system)

	if system.SystemName != "" {
		_, err := db.Exec("UPDATE system_ SET system_name=? WHERE id =?;", system.SystemName, id)
		//fmt.Printf("UPDATE system_ SET system_name=%v WHERE id =%v;\n", system.SystemName, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if system.Address != "" {
		_, err := db.Exec("UPDATE system_ SET address=? WHERE id =?;", system.Address, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if system.Port != 0 {
		_, err := db.Exec("UPDATE system_ SET port=? WHERE id = ?;", system.Port, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if system.AuthenticationInfo != "" {
		_, err := db.Exec("UPDATE system_ SET authentication_info=? WHERE id =?;", system.AuthenticationInfo, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if system.Metadata != nil {

		metadata, _ := json.Marshal(*system.Metadata)
		//fmt.Println(string(metadata))
		_, err := db.Exec("UPDATE system_ SET metadata=? WHERE id =?;", metadata, id)
		if err != nil {
			panic(err.Error())
		}
	}

	ret, err := getSystem(db, id)

	return ret, err
}

// /////////////////////////////////////////////////////////////////////////////
func deleteSystem(db *sql.DB, id int64) error {

	_, err := db.Exec("DELETE FROM system_ WHERE id = ?", id)
	_, err = db.Exec("DELETE FROM service_registry WHERE system_id = ?", id)

	return err
}

// /////////////////////////////////////////////////////////////////////////////
func getAllSystems(db *sql.DB, direction string) ([]dto.SystemResponseDTO, error) {
	var response []dto.SystemResponseDTO = make([]dto.SystemResponseDTO, 0)

	results, err := db.Query("SELECT id, system_name, address, port, authentication_info, created_at, updated_at FROM system_ ORDER BY id " + direction)
	if err != nil {
		return response, nil
	}
	defer results.Close()

	for results.Next() {
		var provider dto.SystemResponseDTO
		var authentication_info sql.NullString
		err = results.Scan(&provider.Id, &provider.SystemName, &provider.Address, &provider.Port, &authentication_info, &provider.CreatedAt, &provider.UpdatedAt)
		if authentication_info.Valid {
			provider.AuthenticationInfo = authentication_info.String
		}
		//fmt.Printf("> %s\n", provider.SystemName)
		response = append(response, provider)
	}

	return response, nil
}

func getAllInterfaces(db *sql.DB) ([]dto.InterfaceDTO, error) {
	var response []dto.InterfaceDTO = make([]dto.InterfaceDTO, 0)

	results, err := db.Query("SELECT id, interface_name FROM service_interface ORDER BY id ASC")
	if err != nil {
		return response, nil
	}
	defer results.Close()

	for results.Next() {
		var sinterface dto.InterfaceDTO

		err = results.Scan(&sinterface.Id, &sinterface.Value)
		response = append(response, sinterface)
	}

	return response, nil
}

func getAllServiceDefinitionsSimple(db *sql.DB) ([]dto.ServiceDTO, error) {
	var ret []dto.ServiceDTO = make([]dto.ServiceDTO, 0)
	result, err := db.Query("SELECT id, service_definition FROM service_definition ORDER BY id ASC")
	if err != nil {
		return ret, err
	}
	defer result.Close()

	for result.Next() {
		var serviceDef dto.ServiceDTO
		result.Scan(&serviceDef.Id, &serviceDef.Value)

		ret = append(ret, serviceDef)
	}

	return ret, nil
}

// /////////////////////////////////////////////////////////////////////////////
func unregisterServiceForSystem(db *sql.DB, service_definition string, system_name string, address string, port int) (bool, error) {
	//log.Printf("unregisterServiceForSystem\n")

	systemID := getIDforSystem(db, system_name, address, port)
	//fmt.Printf("SystemID: %d\n", systemID)
	if systemID == -1 {
		return false, errors.New("No such system:" + system_name)
	}

	serviceDefID := getIDforServiceDefinition(db, service_definition)
	//fmt.Printf("ServiceDefID: %d\n", serviceDefID)

	_, err := db.Exec("DELETE FROM service_registry WHERE service_id=? AND system_id=?", serviceDefID, systemID)
	if err != nil {
		//log.Printf("%v\n", err)
		return false, err
	}

	return true, nil
}

// /////////////////////////////////////////////////////////////////////////////
func unregisterSystem(db *sql.DB, system_name string, address string, port int) (bool, error) {
	//log.Printf("unregisterSystem\n")

	systemID := getIDforSystem(db, system_name, address, port)
	//fmt.Printf("SystemID: %d\n", systemID)
	if systemID == -1 {
		return false, errors.New("No such system:" + system_name)
	}

	_, err := db.Exec("DELETE FROM system_ WHERE id=?", systemID)
	if err != nil {
		//log.Printf("%v\n", err)
		return false, err
	}

	return true, nil
}

// /////////////////////////////////////////////////////////////////////////////
func fetchAllSystems(db *sql.DB, page int, items_per_page int) []dto.SystemResponseDTO {
	var response = []dto.SystemResponseDTO{}

	//log.Print("fetchAllSystems()")

	//results, err := db.Query("SELECT id, system_name, address, port, metadata, created_at, updated_at FROM system_ LIMIT ?;", (page * items_per_page))

	results, err := db.Query("SELECT id, system_name, address, port, authentication_info, metadata, created_at, updated_at FROM system_;")
	if err != nil {
		return nil
	}
	defer results.Close()

	for results.Next() {
		var provider dto.SystemResponseDTO
		var authentication_info, metadata sql.NullString
		err = results.Scan(&provider.Id, &provider.SystemName, &provider.Address, &provider.Port, &authentication_info, &metadata, &provider.CreatedAt, &provider.UpdatedAt)
		if authentication_info.Valid {
			provider.AuthenticationInfo = authentication_info.String
		}
		if metadata.Valid {
			//XXX
		}
		//	fmt.Printf("> %s\n", provider.SystemName)
		response = append(response, provider)
	}

	// XXX: now check if the number of results are greater than total=page*items_per_page

	return response
}

// /////////////////////////////////////////////////////////////////////////////
func fetchServiceById(db *sql.DB, serviceId int64) (dto.ServiceRegistryEntryDTO, error) {
	var ret dto.ServiceRegistryEntryDTO

	//fmt.Printf("fetchServiceById(%d)\n", serviceId)

	result, err := db.Query("SELECT service_id, system_id, service_Uri, end_of_validity, secure, version, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at)  FROM service_registry WHERE id=? LIMIT 1", serviceId)
	if err != nil {
		return ret, err
	}
	defer result.Close()

	var serviceType int64
	var systemID int64
	var end_of_validity sql.NullString
	var created_at, updated_at string
	if result.Next() {
		result.Scan(&serviceType, &systemID, &ret.ServiceUri, &end_of_validity, &ret.Secure, &ret.Version, &created_at, &updated_at)
		if end_of_validity.Valid {
			ret.EndOfValidity = timestamp2Arrowhead(end_of_validity.String)
		}
	} else {
		log.Println(err)
		return ret, errors.New("No such service")
	}

	results2, err2 := db.Query("SELECT service_definition FROM service_definition WHERE id=?", serviceType)
	if err2 != nil {
		return ret, err2
	}
	defer results2.Close()
	results2.Next()
	var serviceDefinition string
	results2.Scan(&serviceDefinition)

	provider, _ := getSystem(db, systemID)
	ret.ServiceDefinition = serviceDefinition
	ret.ProviderSystem.SystemName = provider.SystemName
	ret.ProviderSystem.Address = provider.Address
	ret.ProviderSystem.Port = provider.Port
	ret.ProviderSystem.AuthenticationInfo = provider.AuthenticationInfo
	ret.CreatedAt = timestamp2Arrowhead(created_at)
	ret.UpdatedAt = timestamp2Arrowhead(updated_at)

	results3, err3 := db.Query("SELECT interface_id FROM service_registry_interface_connection WHERE service_registry_id=? LIMIT 1", serviceId)
	if err3 != nil {
		return ret, err3
	}
	defer results3.Close()
	results3.Next()
	var interfaceId int
	results3.Scan(&interfaceId)
	results4, err4 := db.Query("SELECT interface_name FROM service_interface WHERE id=? LIMIT 1", interfaceId)
	if err4 != nil {
		return ret, err4
	}
	defer results4.Close()
	results4.Next()
	var interfaceTypes []string = make([]string, 1)
	results4.Scan(&interfaceTypes[0])
	ret.Interfaces = interfaceTypes

	return ret, nil
}

func deleteServiceById(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM service_registry WHERE id=?;", id)
	return err
}

func getAllServices(db *sql.DB, page *int, item_per_page *int) ([]dto.ServiceRegistryResponseDTO, error) {
	//log.Printf("getAllServices()\n")

	services := make([]dto.ServiceRegistryResponseDTO, 0)
	results, err := db.Query("SELECT id, service_id, system_id, service_uri, UNIX_TIMESTAMP(end_of_validity), secure, version, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM service_registry ORDER BY id;")
	if err != nil {
		fmt.Println(err)
		return services, err
	}
	defer results.Close()

	for results.Next() {
		var serviceID, systemID int64
		var entry dto.ServiceRegistryResponseDTO

		var end_of_validity sql.NullString
		var created_at, updated_at string
		err = results.Scan(&entry.Id, &serviceID, &systemID, &entry.ServiceUri, &end_of_validity, &entry.Secure, &entry.Version, &created_at, &updated_at)
		if end_of_validity.Valid {
			entry.EndOfValidity = timestamp2Arrowhead(end_of_validity.String)
		}
		//fmt.Printf("service_id: %v system_id: %v\n", serviceID, systemID)
		entry.Provider, _ = getSystem(db, systemID)
		entry.Interfaces, _ = getInterfacesForService(db, entry.Id)
		entry.ServiceDefinition, _ = getServiceDefinitionForService(db, serviceID)
		entry.CreatedAt = timestamp2Arrowhead(created_at)
		entry.UpdatedAt = timestamp2Arrowhead(updated_at)

		services = append(services, entry)
	}

	return services, nil
}

func getAllServicesFromServiceDefinition(db *sql.DB, serviceDefinition string, page *int, item_per_page *int) ([]dto.ServiceRegistryResponseDTO, error) {
	//log.Printf("getAllServicesFromServiceDefinition('%s')\n", serviceDefinition)
	services := make([]dto.ServiceRegistryResponseDTO, 0)

	serviceDef, err := getServiceDefinitionFromName(GetSRDB(), serviceDefinition)
	if err != nil {
		//log.Printf("No such service definition: %s\n", serviceDefinition)
		return services, errors.New("No such service definition")
	}

	serviceDefId := serviceDef.Id
	//log.Printf("\tUsing ServiceDefinitionID: %v\n", serviceDefId)

	results, err := db.Query("SELECT id, service_id, system_id, service_uri, end_of_validity, secure, version, created_at, updated_at FROM service_registry WHERE service_id=?", serviceDefId)
	if err != nil {
		fmt.Println(err)
		return services, err
	}
	defer results.Close()

	for results.Next() {
		var serviceID, systemID int64
		var entry dto.ServiceRegistryResponseDTO

		var end_of_validity sql.NullString
		var created_at, updated_at string
		err = results.Scan(&entry.Id, &serviceID, &systemID, &entry.ServiceUri, &end_of_validity, &entry.Secure, &entry.Version, &created_at, &updated_at)
		if end_of_validity.Valid {
			entry.EndOfValidity = timestamp2Arrowhead(end_of_validity.String)
		}

		entry.Provider, _ = getSystem(db, systemID)
		entry.Interfaces, _ = getInterfacesForService(db, entry.Id)
		entry.ServiceDefinition, _ = getServiceDefinitionForService(db, serviceID)
		entry.CreatedAt = timestamp2Arrowhead(created_at)
		entry.UpdatedAt = timestamp2Arrowhead(updated_at)

		fmt.Printf("ENTRY: %+v\n", entry)
		services = append(services, entry)
	}

	fmt.Printf("SERVICES: %+v\n", services)
	return services, nil
}

func addSystem(db *sql.DB, system dto.SystemRequestDTO) (int64, error) {
	stmt, err := db.Prepare("INSERT IGNORE INTO system_(system_name, address, port, authentication_info) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("%v\n", err)
		return -1, err
	}

	res, err := stmt.Exec(system.SystemName, system.Address, system.Port, system.AuthenticationInfo)
	if err != nil {
		return -1, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int64(lid), err
}

// /////////////////////////////////////////////////////////////////////////////
func addServiceDefinition(db *sql.DB, definition string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO service_definition(service_definition) VALUES(?) ON DUPLICATE KEY UPDATE updated_at=NOW()")
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(definition)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int64(lid), err
}

func updateServiceDefinitionById(db *sql.DB, id int64, newDefinition string) error {
	stmt, err := db.Prepare("UPDATE service_definition SET service_definition = ? WHERE id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(newDefinition, id)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if ra != 1 || err != nil {
		return errors.New("No rows updated")
	}

	return nil
}

func deleteServiceDefinitionById(db *sql.DB, id int64) error {
	res, err := db.Exec("DELETE FROM service_definition WHERE id=?;", id)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if ra != 1 || err != nil {
		return errors.New("No rows deleted")
	}

	return nil
}

// helpers
func timestamp2Arrowhead(ts string) string {
	//fmt.Printf("timestamp2Arrowhead(%s)\n", ts)
	intTs, _ := strconv.Atoi(ts)
	timestamp := time.Unix(int64(intTs), 0)

	return timestamp.UTC().Format(time.RFC3339)
	//return timestamp.Format(time.RFC3339)
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
