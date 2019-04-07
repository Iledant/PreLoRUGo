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
	PLAI      int64      `json:"PLAI"`
	PLUS      int64      `json:"PLUS"`
	PLS       int64      `json:"PLS"`
	ANRU      bool       `json:"ANRU"`
}

// Housings embeddes an array of Housing for json export
type Housings struct {
	Housings []Housing `json:"Housing"`
}

// HousingLine is used to decode a line of Housing batch
type HousingLine struct {
	Reference string     `json:"Reference"`
	Address   NullString `json:"Address"`
	ZipCode   NullInt64  `json:"ZipCode"`
	PLAI      int64      `json:"PLAI"`
	PLUS      int64      `json:"PLUS"`
	PLS       int64      `json:"PLS"`
	ANRU      bool       `json:"ANRU"`
}

// HousingBatch embeddes an array of HousingLine for json export
type HousingBatch struct {
	Lines []HousingLine `json:"Housing"`
}

// Validate checks if Housing's fields are correctly filled
func (h *Housing) Validate() error {
	if h.Reference == "" {
		return errors.New("Champ Reference incorrect")
	}
	return nil
}

// Create insert a new Housing into database
func (h *Housing) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO housing (reference,address,zip_code,plai,plus,pls,anru)
 VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, &h.Reference, &h.Address, &h.ZipCode,
		&h.PLAI, &h.PLUS, &h.PLS, &h.ANRU).Scan(&h.ID)
	return err
}

// Update modifies a housing in database
func (h *Housing) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing SET reference=$1,address=$2,zip_code=$3,
	plai=$4,plus=$5,pls=$6,anru=$7 WHERE id=$8`, h.Reference, h.Address, h.ZipCode,
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

// GetAll fetches all Housings from database
func (h *Housings) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,reference,address,zip_code,plai,plus,pls,anru FROM housing`)
	if err != nil {
		return err
	}
	var row Housing
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Reference, &row.Address, &row.ZipCode,
			&row.PLAI, &row.PLUS, &row.PLS, &row.ANRU); err != nil {
			return err
		}
		h.Housings = append(h.Housings, row)
	}
	err = rows.Err()
	if len(h.Housings) == 0 {
		h.Housings = []Housing{}
	}
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

// Save insert a batch of HousingLine into database
func (h *HousingBatch) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_housing (reference,address,zip_code,
		plai,plus,pls,anru) VALUES ($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range h.Lines {
		if r.Reference == "" {
			tx.Rollback()
			return errors.New("Champ Reference incorrect")
		}
		if _, err = stmt.Exec(r.Reference, r.Address, r.ZipCode, r.PLAI, r.PLUS,
			r.PLS, r.ANRU); err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`UPDATE housing SET address=t.address,zip_code=t.zip_code,
	plai=t.plai,plus=t.plus,pls=t.pls,anru=t.anru FROM temp_housing t 
	WHERE t.reference = housing.reference`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO housing (reference,address,zip_code,plai,plus,pls,anru)
	SELECT reference,address,zip_code,plai,plus,pls,anru from temp_housing 
	  WHERE reference NOT IN (SELECT reference from housing)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE from temp_housing")
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
