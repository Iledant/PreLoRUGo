package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// ConventionType is used for normalizing convention type
type ConventionType struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// ConventionTypes embeddes an array of convetion types for json export
type ConventionTypes struct {
	ConventionTypes []ConventionType `json:"ConventionType"`
}

// Create insert a new convention type into database
func (r *ConventionType) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO convention_type(name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
	return err
}

// Valid checks if fields complies with database constraints
func (r *ConventionType) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a ConventionType from database using ID field
func (r *ConventionType) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT name FROM convention_type WHERE ID=$1`, r.ID).
		Scan(&r.Name)
	return err
}

// Update modifies a ConventionType in database
func (r *ConventionType) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE convention_type SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Type de convention introuvable")
	}
	return nil
}

// GetAll fetches all housing transfers from database
func (r *ConventionTypes) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name FROM convention_type`)
	if err != nil {
		return err
	}
	var row ConventionType
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return err
		}
		r.ConventionTypes = append(r.ConventionTypes, row)
	}
	err = rows.Err()
	if len(r.ConventionTypes) == 0 {
		r.ConventionTypes = []ConventionType{}
	}
	return err
}

// Delete removes housing transfer whose ID is given from database
func (r *ConventionType) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM convention_type WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Type de convention introuvable")
	}
	return nil
}
