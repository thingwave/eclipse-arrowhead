package database

import (
	"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var ahdb *sql.DB = nil

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

	ahdb = db
	return db, nil
}

///////////////////////////////////////////////////////////////////////////////
//
//
func GetDB() *sql.DB {
	return ahdb
}

