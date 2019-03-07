package models

import (
	"database/sql"
)

// Beneficiary model
type Beneficiary struct {
	ID   int64  `json:"ID"`
	Code int64  `json:"Code"`
	Name string `json:"Name"`
}

// Beneficiaries embeddes an array of Beneficiary for json export
type Beneficiaries struct {
	Beneficiaries []Beneficiary `json:"Beneficiary"`
}

// GetAll fetches all Beneficiaries from database
func (b *Beneficiaries) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,code,name FROM beneficiary`)
	if err != nil {
		return err
	}
	var row Beneficiary
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Code, &row.Name); err != nil {
			return err
		}
		b.Beneficiaries = append(b.Beneficiaries, row)
	}
	err = rows.Err()
	return err
}
