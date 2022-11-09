package main

import (

	//"strings"
	//"strconv"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	//db "arrowhead.eu/common/database"
)

var cadb *sql.DB = nil

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

///////////////////////////////////////////////////////////////////////////////
//
//
/*func OpenDatabase(address string, port int, username string, password string, dbname string) (*sql.DB, error) {

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

	srdb = db
	return db, nil
}
*/
///////////////////////////////////////////////////////////////////////////////
//
//
func SetCADB(db *sql.DB) {
	cadb = db
}

func GetCADB() *sql.DB {
	return cadb
}
