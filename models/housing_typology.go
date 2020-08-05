package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// HousingTypology is used for normalizing housing typology
type HousingTypology struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// HousingTypologies embeddes an array of HousingTypology for json export
type HousingTypologies struct {
	HousingTypologies []HousingTypology `json:"HousingTypology"`
}

// Create insert a new housing typology into database
func (r *HousingTypology) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO housing_typology(name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
}

// Valid checks if fields complies with database constraints
func (r *HousingTypology) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a HousingTypology from database using ID field
func (r *HousingTypology) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT name FROM housing_typology WHERE ID=$1`, r.ID).
		Scan(&r.Name)
}

// Update modifies a HousingTypology in database
func (r *HousingTypology) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE housing_typology SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Typologie introuvable")
	}
	return nil
}

// GetAll fetches all HousingTypology from database
func (r *HousingTypologies) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,name FROM housing_typology ORDER BY 2`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row HousingTypology
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.HousingTypologies = append(r.HousingTypologies, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.HousingTypologies) == 0 {
		r.HousingTypologies = []HousingTypology{}
	}
	return err
}

// Delete removes RPEvenType whose ID is given from database
func (r *HousingTypology) Delete(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM housing_typology WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Typologie introuvable")
	}
	return nil
}
