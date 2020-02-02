package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// HousingConvention is used for normalizing housing convention
type HousingConvention struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// HousingConventions embeddes an array of housing conventions for json export
type HousingConventions struct {
	HousingConventions []HousingConvention `json:"HousingConvention"`
}

// Create insert a new housing typology into database
func (r *HousingConvention) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO housing_convention(name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
	return err
}

// Valid checks if fields complies with database constraints
func (r *HousingConvention) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a HousingConvention from database using ID field
func (r *HousingConvention) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT name FROM housing_convention WHERE ID=$1`, r.ID).
		Scan(&r.Name)
	return err
}

// Update modifies a HousingConvention in database
func (r *HousingConvention) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing_convention SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Convention introuvable")
	}
	return nil
}

// GetAll fetches all housing conventions from database
func (r *HousingConventions) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name FROM housing_convention`)
	if err != nil {
		return err
	}
	var row HousingConvention
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return err
		}
		r.HousingConventions = append(r.HousingConventions, row)
	}
	err = rows.Err()
	if len(r.HousingConventions) == 0 {
		r.HousingConventions = []HousingConvention{}
	}
	return err
}

// Delete removes RPEvenType whose ID is given from database
func (r *HousingConvention) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM housing_convention WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Convention introuvable")
	}
	return nil
}
