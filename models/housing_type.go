package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// HousingType is used for normalizing housing type
type HousingType struct {
	ID        int64      `json:"ID"`
	ShortName string     `json:"ShortName"`
	LongName  NullString `json:"LongName"`
}

// HousingTypes embeddes an array of housing types for json export
type HousingTypes struct {
	HousingTypes []HousingType `json:"HousingType"`
}

// Create insert a new housing type into database
func (r *HousingType) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO housing_type(short_name,long_name)
	 VALUES($1,$2) RETURNING id`, r.ShortName, r.LongName).Scan(&r.ID)
	return err
}

// Valid checks if fields complies with database constraints
func (r *HousingType) Valid() error {
	if r.ShortName == "" {
		return fmt.Errorf("Nom court vide")
	}
	return nil
}

// Get fetches a HousingType from database using ID field
func (r *HousingType) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT short_name,long_name FROM housing_type WHERE ID=$1`,
		r.ID).Scan(&r.ShortName, &r.LongName)
	return err
}

// Update modifies a HousingType in database
func (r *HousingType) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing_type SET short_name=$1,long_name=$2
	WHERE id=$3`, r.ShortName, r.LongName, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Type introuvable")
	}
	return nil
}

// GetAll fetches all housing transfers from database
func (r *HousingTypes) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,short_name,long_name FROM housing_type`)
	if err != nil {
		return err
	}
	var row HousingType
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.ShortName, &row.LongName); err != nil {
			return err
		}
		r.HousingTypes = append(r.HousingTypes, row)
	}
	err = rows.Err()
	if len(r.HousingTypes) == 0 {
		r.HousingTypes = []HousingType{}
	}
	return err
}

// Delete removes housing transfer whose ID is given from database
func (r *HousingType) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM housing_type WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Type introuvable")
	}
	return nil
}
