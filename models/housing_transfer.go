package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// HousingTransfer is used for normalizing housing transfer
type HousingTransfer struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// HousingTransfers embeddes an array of housing transfers for json export
type HousingTransfers struct {
	HousingTransfers []HousingTransfer `json:"HousingTransfer"`
}

// Create insert a new housing typology into database
func (r *HousingTransfer) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO housing_transfer(name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
	return err
}

// Valid checks if fields complies with database constraints
func (r *HousingTransfer) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a HousingTransfer from database using ID field
func (r *HousingTransfer) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT name FROM housing_transfer WHERE ID=$1`, r.ID).
		Scan(&r.Name)
	return err
}

// Update modifies a HousingTransfer in database
func (r *HousingTransfer) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing_transfer SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Transfert introuvable")
	}
	return nil
}

// GetAll fetches all housing transfers from database
func (r *HousingTransfers) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name FROM housing_transfer`)
	if err != nil {
		return err
	}
	var row HousingTransfer
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return err
		}
		r.HousingTransfers = append(r.HousingTransfers, row)
	}
	err = rows.Err()
	if len(r.HousingTransfers) == 0 {
		r.HousingTransfers = []HousingTransfer{}
	}
	return err
}

// Delete removes housing transfer whose ID is given from database
func (r *HousingTransfer) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM housing_transfer WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Transfert introuvable")
	}
	return nil
}
