package main

import (
  "log"
  "fmt"
  "errors"
  "time"
  "strconv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  dto "arrowhead.eu/common/datamodels"
)

type Provider struct {
    ID   int  `json:"id"`
//    System_id int `json:"system_id"`
    SystemName string `json:"systemName"`
    Address string `json:"address"`
//    Service_IP string `json:"service_ip"`
    Port int  `json:"port"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
}

type ServiceInterface struct {
    ID   int    `json:"id"`
    Interface_name string `json:"interface_name"`
}

var audb *sql.DB = nil

///////////////////////////////////////////////////////////////////////////////
//
//
func OpenDatabase(address string, port int, username string, password string, dbname string) (*sql.DB, error) {

  // Open up our database connection. XXX fix login parameters
  db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+":3306)/"+dbname+"?parseTime=true")

  // if there is an error opening the connection, handle it
  if err != nil {
	  fmt.Println("Could not connect to MySQL database")
	  return nil, err
  }

  err = db.Ping()
  if err != nil {
	  fmt.Println("Could not connect to MySQL database")
	  db.Close()
	  return nil, err
  }

  audb = db
  return db, nil
}

func GetAUDB() *sql.DB {
  return audb
}

func getSystem(db *sql.DB, systemId int64) (dto.SystemResponseDTO, error) {
        var ret dto.SystemResponseDTO

        //log.Printf("getSystem(%v)\n", systemId)

        result, err := db.Query("SELECT id, system_name, address, port, authentication_info, metadata, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM system_ WHERE id=? LIMIT 1;", systemId)
        if err != nil {
                fmt.Println(err)
                return ret, err
        }
        defer result.Close()

        result.Next()
        var authentication_info, metadata sql.NullString
        var created_at, updated_at string
        err = result.Scan(&ret.Id, &ret.SystemName, &ret.Address, &ret.Port, &authentication_info, &metadata, &created_at, &updated_at)
        if err != nil {
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

        return ret, err
}


func GetAllIntraCloudRules(db *sql.DB) ([]AuthorizationIntraCloudResponseDTO, error) {
  var ret []AuthorizationIntraCloudResponseDTO = make([]AuthorizationIntraCloudResponseDTO, 0)

  results, err := db.Query("SELECT id, created_at, updated_at, consumer_system_id, provider_system_id, service_id FROM authorization_intra_cloud;")
  if err != nil {
    log.Printf("%v\n", err)
    return ret, err
  }
  defer results.Close()

  for results.Next() {
    var ruleEntry AuthorizationIntraCloudResponseDTO
    var service_id int64
    var consumerSystemId int64
    var providerSystemId int64
    err2 := results.Scan(&ruleEntry.Id, &ruleEntry.CreatedAt, &ruleEntry.UpdatedAt, &consumerSystemId, &providerSystemId, &service_id)
    if err2 != nil {
      fmt.Printf("RESULTS SCAN() error: %v\n", err2)
      continue
    }

    ruleEntry.ConsumerSystem, _ = getSystem(db, consumerSystemId)
    ruleEntry.ProviderSystem, _ = getSystem(db, providerSystemId)
    result2, _ := db.Query("SELECT id, service_definition, created_at, updated_at FROM service_definition WHERE id=?", service_id)
    defer result2.Close()

    result2.Next()
    var srv ServiceDefinitionResponseDTO
    _ = result2.Scan(&srv.Id, &srv.ServiceDefinition, &srv.CreatedAt, &srv.UpdatedAt)
    ruleEntry.ServiceDefinition = srv

    result3, err3 := db.Query("SELECT two.id, two.interface_name, two.created_at, two.updated_at FROM authorization_intra_cloud_interface_connection one INNER JOIN service_interface two WHERE two.id = one.interface_id AND one.id=?;", ruleEntry.Id)
    if err3 != nil {
      fmt.Printf("err3: %v\n", err3)
    }
    defer result3.Close()

    ruleEntry.Interfaces = make([]ServiceInterfaceResponseDTO, 0)
    for result3.Next() {
      var aici ServiceInterfaceResponseDTO
      _ = result3.Scan(&aici.Id, &aici.InterfaceName, &aici.CreatedAt, &aici.UpdatedAt)
      ruleEntry.Interfaces = append(ruleEntry.Interfaces, aici)
    }

    //fmt.Printf("ruleEntry: %v\n", ruleEntry)
    ret = append(ret, ruleEntry)
  }

  return ret, nil
}

func GetOrInsertAuthorizationRule(db *sql.DB, consumerId int64, providerId int64, serviceId int64) (int64, error) {
  log.Printf("GetOrInsertAuthorizationRule(%v, %v, %v):\n", consumerId, providerId, serviceId);

  result, err := db.Query("SELECT id FROM authorization_intra_cloud WHERE consumer_system_id=? AND provider_system_id=? AND service_id=?;", consumerId, providerId, serviceId)
  if err != nil {
    fmt.Printf("%v\n", err)
    return -1, err
  }
  defer result.Close()

  if result.Next() {

    var id int64
    _ = result.Scan(&id)
    fmt.Printf("Found rule id: %v\n", id)
    return id, nil
  } else {
    fmt.Printf("No entry! Inject!!\n")

    //sql1 := fmt.Sprintf("INSERT INTO authorization_intra_cloud(consumer_system_id, provider_system_id, service_id) VALUES(%v, %v, %v);", consumerId, providerId, serviceId)
    //fmt.Printf("%s\n", sql1)
    res, err := db.Exec("INSERT INTO authorization_intra_cloud(consumer_system_id, provider_system_id, service_id) VALUES(?, ?, ?);", consumerId, providerId, serviceId)
    if err != nil {
      fmt.Printf("INJECT err: %s\n", err)
      return -1, err
    }
    id, err := res.LastInsertId()
    if err != nil {
      println("Error:", err.Error())
    } else {
      println("LastInsertId:", id)
      return id, nil
    }

  }

  return -1, errors.New("Something!")

}

func GetIntraCloudRuleById(db *sql.DB, ruleId int64) (AuthorizationIntraCloudResponseDTO, error) {
  var ret AuthorizationIntraCloudResponseDTO

  results, err := db.Query("SELECT id, created_at, updated_at, consumer_system_id, provider_system_id, service_id FROM authorization_intra_cloud WHERE id=? LIMIT 1", ruleId)
  if err != nil {
    log.Printf("%v\n", err)
    return ret, err
  }
  defer results.Close()

  if results.Next() {
    var service_id int64
    var consumerSystemId int64
    var providerSystemId int64
    err2 := results.Scan(&ret.Id, &ret.CreatedAt, &ret.UpdatedAt, &consumerSystemId, &providerSystemId, &service_id)
    if err2 != nil {
      fmt.Printf("RESULTS SCAN() error: %v\n", err2)
      return ret, err2
    }

    ret.ConsumerSystem, _ = getSystem(db, consumerSystemId)
    ret.ProviderSystem, _ = getSystem(db, providerSystemId)
    result2, _ := db.Query("SELECT id, service_definition, created_at, updated_at FROM service_definition WHERE id=?", service_id)
    defer result2.Close()

    result2.Next()
    var srv ServiceDefinitionResponseDTO
    _ = result2.Scan(&srv.Id, &srv.ServiceDefinition, &srv.CreatedAt, &srv.UpdatedAt)
    ret.ServiceDefinition = srv

    result3, err3 := db.Query("SELECT two.id, two.interface_name, two.created_at, two.updated_at FROM authorization_intra_cloud_interface_connection one INNER JOIN service_interface two WHERE two.id = one.interface_id AND one.id=?", ret.Id)
    if err3 != nil {
      fmt.Printf("err3: %v\n", err3)
      return ret, err3
    }
    defer result3.Close()

    ret.Interfaces = make([]ServiceInterfaceResponseDTO, 0)
    for result3.Next() {
      var aici ServiceInterfaceResponseDTO
      _ = result3.Scan(&aici.Id, &aici.InterfaceName, &aici.CreatedAt, &aici.UpdatedAt)
      ret.Interfaces = append(ret.Interfaces, aici)
    }

    return ret, nil
  }

  return ret, errors.New("No such rule")
}

func AddIntraCloudRuleBy(db *sql.DB, rule AuthorizationIntraCloudRequestDTO) (bool, error) { //XXX IMPLEMENT ME!!

  return false, nil
}

func DeleteIntraCloudRuleById(db *sql.DB, ruleId int64) (bool, error) {

  result, err := db.Exec("DELETE FROM authorization_intra_cloud WHERE id=?;", ruleId)
  if err != nil {
    return false, err
  } else {
    nbr, err :=  result.RowsAffected()
    if nbr == 1 && err == nil {
      return true, nil
    }
  }

  return false, err

}

// helpers
func timestamp2Arrowhead(ts string) string {
        //fmt.Printf("timestamp2Arrowhead(%s)\n", ts)
        intTs, _ := strconv.Atoi(ts)
        timestamp := time.Unix(int64(intTs), 0)

        return timestamp.UTC().Format(time.RFC3339)
        //return timestamp.Format(time.RFC3339)
}
