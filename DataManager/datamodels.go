package main

import (
	"errors"
	"fmt"
	"time"
	dto "arrowhead.eu/common/datamodels"
)

type SystemList struct {
	Systems []string `json:"systems"`
}

type SystemServiceList struct {
	Services []string `json:"services"`
}

type SignalProperties struct {
	SigName  string //sig0
	SigKey   string // temperature
	SigCount int    // 7
}
/*
type SenMLEntry struct {
	Bn   *string  `json:"bn,omitempty"`
	Bt   *float64 `json:"bt,omitempty"`
	Bu   *string  `json:"bu,omitempty"`
	Bv   *float64 `json:"bv,omitempty"`
	Bs   *float64 `json:"bs,omitempty"`
	Bver *int     `json:"bver,omitempty"`

	N  *string  `json:"n,omitempty"`
	U  *string  `json:"u,omitempty"`
	V  *float64 `json:"v,omitempty"`
	Vs *string  `json:"vs,omitempty"`
	Vb *bool    `json:"vb,omitempty"`
	Vd *string  `json:"vd,omitempty"`

	S  *float64 `json:"s,omitempty"`
	T  *float64 `json:"t,omitempty"`
	Ut *float64 `json:"ut,omitempty"`
}*/

const RELATIVE_TIMESTAMP_INDICATOR = float64(268435456)

///////////////////////////////////////////////////////////////////////////////
//
func validateSenML(senml []dto.SenMLEntry) error {

	if len(senml) < 1 {
		return errors.New("SenML array must be greater of equal 1")
	}

	if senml[0].Bn == nil || *senml[0].Bn == "" {
		return errors.New("array[] must contain bn tag")
	}

	if senml[0].Bt == nil {
		senml[0].Bt = new(float64)
		*senml[0].Bt = float64(time.Now().UnixNano() / 1e6)
		*senml[0].Bt /= 1000.0
	}

	for i, e := range senml {
		fmt.Println(i, e)
		fmt.Printf("\n")

		if i == 0 {
			continue
		}
		if e.Bn != nil || e.Bt != nil || e.Bu != nil || e.Bv != nil || e.Bs != nil || e.Bver != nil {
			return errors.New("bX tags are only allowed in position [0]")
		}

		if e.N == nil || *e.N == "" {
			return errors.New("An n tag MUST exist in each element at pos >= 1")
		}

		// add more cases below...
		if e.T == nil {
			e.T = senml[0].Bt
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
//
func checkSenMLParameters(senml []dto.SenMLEntry) error {

	return nil
}
