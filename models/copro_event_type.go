package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// CoproEventType is used for normalizing types of events attached to a copro project
type CoproEventType struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// CoproEventTypes embeddes an array of CoproEventType for json export
type CoproEventTypes struct {
	CoproEventTypes []CoproEventType `json:"CoproEventType"`
}

// Create insert a new CoproEventType into database
func (r *CoproEventType) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO copro_event_type (name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
	return err
}

// Validate check if fields complies with database constraints
func (r *CoproEventType) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name vide")
	}
	return nil
}

// Get fetches a CoproEventType from database using ID field
func (r *CoproEventType) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT name FROM copro_event_type WHERE ID=$1`, r.ID).
		Scan(&r.Name)
	return err
}

// Update modifies a CoproEventType in database
func (r *CoproEventType) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE copro_event_type SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Type d'événement introuvable")
	}
	return nil
}

// GetAll fetches all CoproEventType from database
func (r *CoproEventTypes) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name FROM copro_event_type`)
	if err != nil {
		return err
	}
	var row CoproEventType
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return err
		}
		r.CoproEventTypes = append(r.CoproEventTypes, row)
	}
	err = rows.Err()
	if len(r.CoproEventTypes) == 0 {
		r.CoproEventTypes = []CoproEventType{}
	}
	return err
}

// Delete removes CoproEvenType whose ID is given from database
func (r *CoproEventType) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM copro_event_type WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Type d'événement introuvable")
	}
	return nil
}
