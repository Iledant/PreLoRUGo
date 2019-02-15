package models

import (
	"database/sql"
	"errors"
)

// Housing model
type Housing struct {
	ID        int64      `json:"ID"`
	Reference string     `json:"Reference"`
	Address   NullString `json:"Address"`
	ZipCode   NullInt64  `json:"ZipCode"`
	PLAI      int        `json:"PLAI"`
	PLUS      int        `json:"PLUS"`
	PLS       int        `json:"PLS"`
	ANRU      bool       `json:"ANRU"`
}

// Housings embeddes an array of Housing for json export
type Housings struct {
	Housings []Housing `json:"Housing"`
}

// Validate checks if Housing's fields are correctly filled
func (h *Housing) Validate() error {
	if h.Reference == "" {
		return errors.New("Champ reference incorrect")
	}
	return nil
}

// Create insert a new Housing into database
func (h *Housing) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO housing (reference,address,zip_code,plai,plus,pls,anru) 
	VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, &h.Reference, &h.Address,
		&h.ZipCode, &h.PLAI, &h.PLUS, &h.PLS, &h.ANRU).Scan(&h.ID)
	return err
}

// Update modifies a housing in database
func (h *Housing) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing SET reference=$1,address=$2,zip_code=$3,
	plai=$4,plus=$5,pls=$6,anru=$7 WHERE id = $8`, h.Reference, h.Address, h.ZipCode,
		h.PLAI, h.PLUS, h.PLS, h.ANRU, h.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Logement introuvable")
	}
	return err
}

// GetAll fetches all housings from database
func (h *Housings) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,reference,address,zip_code,plai,plus,pls,anru
	FROM housing`)
	if err != nil {
		return err
	}
	var r Housing
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ID, &r.Reference, &r.Address, &r.ZipCode, &r.PLAI,
			&r.PLUS, &r.PLS, &r.ANRU); err != nil {
			return err
		}
		h.Housings = append(h.Housings, r)
	}
	err = rows.Err()
	return err
}

// Delete removes housing whose ID is given from database
func (h *Housing) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM housing WHERE id = $1", h.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Logement introuvable")
	}
	return nil
}
