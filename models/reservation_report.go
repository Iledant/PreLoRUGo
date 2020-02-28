package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// ReservationReport model
type ReservationReport struct {
	ID             int64      `json:"ID"`
	BeneficiaryID  int64      `json:"BeneficiaryID"`
	Beneficiary    string     `json:"Beneficiary"`
	Area           float64    `json:"Area"`
	SourceIRISCode string     `json:"SourceIRISCode"`
	DestIRISCode   NullString `json:"DestIRISCode"`
	DestDate       NullTime   `json:"DestDate"`
}

// ReservationReports embeddes an array of ReservationReport for json export and
// queries
type ReservationReports struct {
	Lines []ReservationReport `json:"ReservationReport"`
}

// Valid checks if fields respects database constraints
func (r *ReservationReport) Valid() error {
	if r.BeneficiaryID == 0 {
		return errors.New("BeneficiaryID vide")
	}
	if r.SourceIRISCode == "" {
		return errors.New("SourceIRISCode vide")
	}
	return nil
}

// Create save a reservation report into the database
func (r *ReservationReport) Create(db *sql.DB) error {
	if err := db.QueryRow(`INSERT INTO reservation_report(beneficiary_id,area,
		source_iris_code,dest_iris_code,dest_date) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		r.BeneficiaryID, r.Area, r.SourceIRISCode, r.DestIRISCode, r.DestDate).
		Scan(&r.ID); err != nil {
		return fmt.Errorf("insert %v", err)
	}
	return db.QueryRow(`SELECT name from beneficiary WHERE id=$1`, r.BeneficiaryID).
		Scan(&r.Beneficiary)
}

// Update modifies a reservation report into the database
func (r *ReservationReport) Update(db *sql.DB) error {
	result, err := db.Exec(`UPDATE reservation_report SET beneficiary_id=$1,area=$2,
		source_iris_code=$3,dest_iris_code=$4,dest_date=$5 WHERE id=$6`,
		r.BeneficiaryID, r.Area, r.SourceIRISCode, r.DestIRISCode, r.DestDate, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count == 0 {
		return errors.New("report de réservation introuvable")
	}
	return db.QueryRow(`SELECT name from beneficiary WHERE id=$1`, r.BeneficiaryID).
		Scan(&r.Beneficiary)
}

// Delete removes a reservation report from database
func (r *ReservationReport) Delete(db *sql.DB) error {
	result, err := db.Exec(`DELETE FROM reservation_report WHERE id=$1`, r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count == 0 {
		return errors.New("report de réservation introuvable")
	}
	return nil
}

// GetAll fetches all reservation report from database
func (r *ReservationReports) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT r.id,r.beneficiary_id,r.area,
	r.source_iris_code,r.dest_iris_code,r.dest_date,b.name
	FROM reservation_report r
	JOIN beneficiary b ON b.id=r.beneficiary_id`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var l ReservationReport
	for rows.Next() {
		if err = rows.Scan(&l.ID, &l.BeneficiaryID, &l.Area, &l.SourceIRISCode,
			&l.DestIRISCode, &l.DestDate, &l.Beneficiary); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.Lines = append(r.Lines, l)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.Lines) == 0 {
		r.Lines = []ReservationReport{}
	}
	return nil
}
