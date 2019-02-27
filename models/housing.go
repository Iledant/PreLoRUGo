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
