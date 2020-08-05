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
func (r *HousingConvention) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO housing_convention(name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
}

// Valid checks if fields complies with database constraints
func (r *HousingConvention) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a HousingConvention from database using ID field
func (r *HousingConvention) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT name FROM housing_convention WHERE ID=$1`, r.ID).
		Scan(&r.Name)
}

// Update modifies a HousingConvention in database
func (r *HousingConvention) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE housing_convention SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	if count != 1 {
		return errors.New("Convention introuvable")
	}
	return nil
}

// GetAll fetches all housing conventions from database
func (r *HousingConventions) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,name FROM housing_convention`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row HousingConvention
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.HousingConventions = append(r.HousingConventions, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.HousingConventions) == 0 {
		r.HousingConventions = []HousingConvention{}
	}
	return nil
}

// Delete removes RPEvenType whose ID is given from database
func (r *HousingConvention) Delete(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM housing_convention WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Convention introuvable")
	}
	return nil
}
