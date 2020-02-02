package models

import (
	"database/sql"
	"fmt"
)

// ReservationFee model
type ReservationFee struct {
	ID                   int64       `json:"ID"`
	CurrentBeneficiaryID int64       `json:"CurrentBeneficiaryID"`
	CurrentBeneficiary   string      `json:"CurrentBeneficiary"`
	PastBeneficiaryID    NullInt64   `json:"PastBeneficiaryID"`
	PastBeneficiary      NullString  `json:"PastBeneficiary"`
	CityCode             int64       `json:"CityCode"`
	AddressNumber        NullString  `json:"AddressNumber"`
	AddressStreet        NullString  `json:"AddressStreet"`
	RPLS                 NullString  `json:"RPLS"`
	ConventionID         NullInt64   `json:"ConventionID"`
	Convention           NullString  `json:"Convention"`
	Count                int64       `json:"Count"`
	TransferDate         NullTime    `json:"TransferDate"`
	CommentID            NullInt64   `json:"CommentID"`
	Comment              NullString  `json:"Comment"`
	TransferID           NullInt64   `json:"TransferID"`
	Transfer             NullString  `json:"Transfer"`
	ConventionDate       NullTime    `json:"ConventionDate"`
	EliseRef             NullString  `json:"EliseRef"`
	Area                 NullFloat64 `json:"Area"`
	EndYear              NullInt64   `json:"EndYear"`
	Loan                 NullFloat64 `json:"Loan"`
	Charges              NullFloat64 `json:"Charges"`
}

// ReservationFees embeddes an array of ReservationFee for json export and dedicated
// queries
type ReservationFees struct {
	Lines []ReservationFee `json:"ReservationFee"`
}

// Valid checks if fields complies with database constraints
func (r *ReservationFee) Valid() error {
	if r.CurrentBeneficiaryID == 0 {
		return fmt.Errorf("CurrentBeneficiaryID null")
	}
	if r.CityCode == 0 {
		return fmt.Errorf("CityCode null")
	}
	if r.Count == 0 {
		return fmt.Errorf("Count null")
	}
	return nil
}

// Create insert a new ReservationFee into database
func (r *ReservationFee) Create(db *sql.DB) error {
	if err := db.QueryRow(`INSERT into reservation_fee (current_beneficiary_id,
		past_beneficiary_id,city_code,address_number,address_street,rpls,
		convention_id,count,transfer_date,transfer_id,comment_id,convention_date,
		elise_ref,area,end_year,loan,charges) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,
			$10,$11,$12,$13,$14,$15,$16,$17) RETURNING ID`, &r.CurrentBeneficiaryID,
		&r.PastBeneficiaryID, &r.CityCode, &r.AddressNumber, &r.AddressStreet,
		&r.RPLS, &r.ConventionID, &r.Count, &r.TransferDate, &r.TransferID,
		&r.CommentID, &r.ConventionDate, &r.EliseRef, &r.Area, &r.EndYear, &r.Loan,
		&r.Charges).Scan(&r.ID); err != nil {
		return fmt.Errorf("insert %v", err)
	}
	if err := db.QueryRow(`SELECT b1.name,b2.name,hc.name,ht.name,ho.name
	FROM reservation_fee rf 
	JOIN beneficiary b1 ON rf.current_beneficiary_id=b1.id
	LEFT OUTER JOIN beneficiary b2 ON past_beneficiary_id=b2.id
	LEFT OUTER JOIN housing_convention hc ON hc.id=rf.convention_id
	LEFT OUTER JOIN housing_transfer ht ON ht.id=rf.transfer_id
	LEFT OUTER JOIN housing_comment ho ON ho.id=rf.comment_id
	WHERE rf.id=$1`, r.ID).Scan(&r.CurrentBeneficiary, &r.PastBeneficiary,
		&r.Convention, &r.Transfer, &r.Comment); err != nil {
		return fmt.Errorf("select %v", err)
	}
	return nil
}
