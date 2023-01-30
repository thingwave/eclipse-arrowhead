package main

import (
	"database/sql"
	"errors"
	"fmt"

	dto "arrowhead.eu/common/datamodels"
	util "arrowhead.eu/common/util"
	_ "github.com/go-sql-driver/mysql"
)

var ordb *sql.DB = nil

// /////////////////////////////////////////////////////////////////////////////
func SetORDB(db *sql.DB) {
	ordb = db
}
func GetOrDB() *sql.DB {
	return ordb
}

// /////////////////////////////////////////////////////////////////////////////
func getSystem(db *sql.DB, systemId int64) (dto.SystemResponseDTO, error) {
	var ret dto.SystemResponseDTO
	ret.Id = -1

	//fmt.Printf("getSystem(%v)\n", systemId)

	result, err := db.Query("SELECT id, system_name, address, port, authentication_info, metadata, UNIX_TIMESTAMP(created_at), UNIX_TIMESTAMP(updated_at) FROM system_ WHERE id=? LIMIT 1;", systemId)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	defer result.Close()

	if result.Next() {
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

		ret.CreatedAt = util.Timestamp2Arrowhead(created_at)
		ret.UpdatedAt = util.Timestamp2Arrowhead(updated_at)
		return ret, nil
	}

	return ret, errors.New(fmt.Sprintf("System with id %v not found.", systemId))

}

// XXX: NOT FULLY IMPLEMENTED
func GetTopPriorityEntries(db *sql.DB) ([]dto.StoreEntry, error) {
	res := make([]dto.StoreEntry, 0)

	results, err := db.Query("SELECT * FROM orchestrator_store WHERE priority=1;")
	if err != nil {
		return res, nil
	}
	defer results.Close()

	for results.Next() {
		var storeEntry dto.StoreEntry
		//_ = results.Scan(&storeEntry.Id, &interfaceEntry.InterfaceName, &interfaceEntry.CreatedAt, &interfaceEntry.UpdatedAt) //XXX: BUG HERE
		//fmt.Println("Interface name: ", interfaceEntry.InterfaceName)
		res = append(res, storeEntry)
	}

	return res, nil
}

func UpdatePriorityForSystem(db *sql.DB, systemId string, priority int) error {
	fmt.Printf("UpdatePriorityForSystem(\"%s\"): %v\n", systemId, priority)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
//

// helpers
