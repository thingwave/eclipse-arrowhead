package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"database/sql"

	dto "arrowhead.eu/common/datamodels"
	_ "github.com/go-sql-driver/mysql"
)

var dmdb *sql.DB = nil

const MAX_ALLOWED_TIME_DIFF = 1000

// /////////////////////////////////////////////////////////////////////////////
func GetDMDB() *sql.DB {
	return dmdb
}

func SetDMDB(db *sql.DB) {
	dmdb = db
}

// /////////////////////////////////////////////////////////////////////////////
func GetDMHistSystems(db *sql.DB) ([]string, error) {
	var ret = []string{}

	result, err := db.Query("SELECT DISTINCT(system_name) FROM dmhist_services;")
	if err != nil {
		//    panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}

	var system_name string
	for result.Next() {
		result.Scan(&system_name)
		//fmt.Printf("Sys: %s\n", system_name)
		ret = append(ret, system_name)
	}

	return ret, nil
}

// /////////////////////////////////////////////////////////////////////////////
func GetDMHistSystemServices(db *sql.DB, sysName string) []string {
	var ret = []string{}

	result, err := db.Query("SELECT DISTINCT(service_name) FROM dmhist_services WHERE system_name=?;", sysName)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil
	}

	var service_name string
	for result.Next() {
		result.Scan(&service_name)
		fmt.Printf("Sys: %s\n", service_name)
		ret = append(ret, service_name)
	}

	return ret
}

func GetDMHistSystemServiceDataSignals(db *sql.DB, sysName string, srvName string, from float64, to float64, count int, signals []SignalProperties) ([]dto.SenMLEntry, error) {
	var response = []dto.SenMLEntry{}

	return response, errors.New("Not implemented!")
}

// /////////////////////////////////////////////////////////////////////////////
func GetDMHistSystemServiceData(db *sql.DB, sysName string, srvName string, from float64, to float64, count int) ([]dto.SenMLEntry, error) {
	var response = []dto.SenMLEntry{}

	serviceId := serviceToID(db, sysName, srvName)
	if serviceId == -1 {
		return nil, errors.New("GetDMHistSystemServiceData: no service ID found")
	}
	fmt.Printf("Using serviceId of: %v\n", serviceId)

	if from < 0.0 {
		from = 0.0 //1970-01-01
	}
	if to < 0.0 {
		to = float64(MAX_ALLOWED_TIME_DIFF + (float64(time.Now().Unix())))
	}

	result, err := db.Query("SELECT id FROM dmhist_messages WHERE sid=? AND bt >=? AND bt <=? ORDER BY bt DESC LIMIT ?", serviceId, from, to, count)
	if err != nil {
		//panic(err.Error())
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var mid int
		result.Scan(&mid)
		fmt.Printf("using mid of: %v\n", mid)

		result2, err := db.Query("SELECT n, u, t, v, vb, CONVERT(vs USING utf8) FROM dmhist_entries WHERE sid=? AND mid=? AND t>=? AND t <=? ORDER BY t DESC", serviceId, mid, from, to)
		if err != nil {
			return nil, err
		}
		defer result2.Close()

		for result2.Next() {
			var e dto.SenMLEntry
			e.N = new(string)
			var u sql.NullString
			e.T = new(float64)

			var v sql.NullFloat64
			var vb sql.NullBool
			var vs sql.NullString
			result2.Scan(e.N, &u, e.T, &v, &vb, &vs)

			if u.Valid {
				e.U = new(string)
				*e.U = u.String
			}
			if v.Valid {
				e.V = new(float64)
				*e.V = v.Float64
			}
			if vb.Valid {
				e.Vb = new(bool)
				*e.Vb = vb.Bool
			}
			if vs.Valid {
				e.Vs = new(string)
				*e.Vs = vs.String
			}
			jsonRespStr, _ := json.Marshal(e)
			fmt.Printf("\te: %s\n", jsonRespStr)

			response = append(response, e)
			jsonRespStr, _ = json.Marshal(response)
			//fmt.Printf("RESP: %s\n", jsonRespStr)
		}

	}

	if len(response) >= 1 {
		response[0].Bn = new(string)
		*response[0].Bn = sysName
	}

	if len(response) == 0 {

	} else if len(response) == 1 {
		response[0].Bt = new(float64)
		*response[0].Bt = *response[0].T
		response[0].T = nil
	} else {
		response[0].Bt = new(float64)
		*response[0].Bt = *response[0].T
		response[0].T = nil
		for i, _ := range response {
			if i == 0 {
				continue
			}
			*response[i].T = *response[i].T - *response[0].Bt
			*response[i].T = math.Round(*response[i].T*100) / 100
		}
	}

	return response, nil
}

// /////////////////////////////////////////////////////////////////////////////
func PutDMHistSystemServiceData(db *sql.DB, sysName string, srvName string, body string, message []dto.SenMLEntry) error {
	var ret error = nil

	log.Printf("PutDMHistSystemServiceData:\n%s", body)

	if len(message) == 0 {
		return errors.New("Empty message")
	}

	if message[0].Bt == nil {
		fmt.Printf("bt is nil, creating!")
		var bt64 float64
		message[0].Bt = &bt64
		*message[0].Bt = float64(time.Now().UnixNano() / 1e6)
	}
	bt := *message[0].Bt

	maxTs := getLargestTimestamp(message)
	minTs := getSmallestTimestamp(message)

	var bu *string = nil
	if message[0].Bu != nil {
		bu = message[0].Bu
	}

	serviceId := serviceToID(db, sysName, srvName)
	if serviceId == -1 {
		//return errors.New("PutDMHistSystemServiceData: no service ID found")
		serviceId, ret = createEndpoint(db, sysName, srvName)
		if ret != nil {
			return errors.New("PutDMHistSystemServiceData: cannot create system")
		}
	}
	//fmt.Printf("PUT Using serviceId of: %v for %s\n", serviceId, body)

	//fmt.Printf("Storing BODY:\n%s\n", body)

	stmt, err := db.Prepare("INSERT INTO dmhist_messages(sid, bt, mint, maxt, msg) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("PREPARE error: %v\n", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(serviceId, bt, minTs, maxTs, body)
	if err != nil {
		fmt.Printf("EXEC error: %v\n", err)
	}
	mid, _ := res.LastInsertId()
	log.Printf("MID: %v:\n", mid)
	for i, e := range message {
		fmt.Printf("%d: %+v:\n", i, e)
		var t float64 = bt
		if e.T != nil {
			t += *e.T
		}

		var u sql.NullString
		u.Valid = false
		if e.U != nil {
			u.String = *e.U
			u.Valid = true
		} else if bu != nil {
			u.String = *bu
			u.Valid = true
		}

		if e.V != nil {
			sql := "INSERT INTO dmhist_entries(sid, mid, n, t, u, v) VALUES(?,?,?,?,?,?)"
			res, err = db.Exec(sql, serviceId, mid, *e.N, t, u, *e.V)
		} else if e.Vs != nil {
			sql := fmt.Sprintf("INSERT INTO dmhist_entries(sid, mid, n, t, u, %s) VALUES(?, ?, ?, ?, ?, ?)", "vs")
			res, err = db.Exec(sql, serviceId, mid, *e.N, t, u, *e.Vs)
		} else if e.Vb != nil {
			sql := fmt.Sprintf("INSERT INTO dmhist_entries(sid, mid, n, t, u, %s) VALUES(?, ?, ?, ?, ?, ?)", "vb")
			res, err = db.Exec(sql, serviceId, mid, *e.N, t, u, *e.Vb)
		} else if e.Vd != nil {
			sql := fmt.Sprintf("INSERT INTO dmhist_entries(sid, mid, n, t, u, %s) VALUES(?, ?, ?, ?, ?, ?)", "vd")
			res, err = db.Exec(sql, serviceId, mid, *e.N, t, u, *e.Vd)
		}

		if err != nil {
			fmt.Printf("INSERT error: %v\n", err)
		}
	}

	return err
}

//=================================================================================================
// assistant methods

func createEndpoint(db *sql.DB, systemName string, serviceName string) (int64, error) {

	serviceId := serviceToID(db, systemName, serviceName)
	if serviceId != -1 {
		return serviceId, nil
	}

	res, err := db.Exec("INSERT INTO dmhist_services(system_name, service_name) "+"VALUES(?,?);", systemName, serviceName)
	if err != nil {
		return -1, err
	}
	//defer res.Close()
	serviceId, err = res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return serviceId, nil
}

// =================================================================================================
// returns largest (newest) timestamp value
func getLargestTimestamp(message []dto.SenMLEntry) float64 {
	var ret float64 = *message[0].Bt
	bt := *message[0].Bt

	for _, e := range message {

		if e.T == nil {
			continue
		}

		if *e.T > RELATIVE_TIMESTAMP_INDICATOR { // absolute
			if *e.T > ret {
				ret = *e.T
			}
		} else {
			if *e.T+bt > ret {
				ret = *e.T + bt
			}
		}
		//fmt.Printf("%v: %v\n", i, e)
	}

	return ret
}

// =================================================================================================
// returns largest (newest) timestamp value
func getSmallestTimestamp(message []dto.SenMLEntry) float64 {
	var ret float64 = *message[0].Bt
	bt := *message[0].Bt

	for _, e := range message {

		if e.T == nil {
			continue
		}

		if *e.T > RELATIVE_TIMESTAMP_INDICATOR { // absolute
			if *e.T < ret {
				ret = *e.T
			}
		} else {
			if *e.T+bt < ret {
				ret = *e.T + bt
			}
		}

	}
	return ret
}

// /////////////////////////////////////////////////////////////////////////////
func serviceToID(db *sql.DB, sysName string, srvName string) int64 {
	var ret int64 = -1

	result, err := db.Query("SELECT id FROM dmhist_services WHERE system_name=? AND service_name=? LIMIT 1;", sysName, srvName)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return -1
	}
	defer result.Close()

	result.Next()
	result.Scan(&ret)

	fmt.Printf("found serviceId: %v for %s\n", ret, srvName)
	return ret
}
