package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// RPEventType is used for normalizing types of events attached to a renew project
type RPEventType struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// RPEventTypes embeddes an array of RPEventType for json export
type RPEventTypes struct {
	RPEventTypes []RPEventType `json:"RPEventType"`
}

// Create insert a new RPEventType into database
func (r *RPEventType) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO rp_event_type (name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
}

// Validate check if fields complies with database constraints
func (r *RPEventType) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a RPEventType from database using ID field
func (r *RPEventType) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT name FROM rp_event_type WHERE ID=$1`, r.ID).
		Scan(&r.Name)
}

// Update modifies a RPEventType in database
func (r *RPEventType) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE rp_event_type SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Type d'événement introuvable")
	}
	return nil
}

// GetAll fetches all RPEventType from database
func (r *RPEventTypes) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,name FROM rp_event_type`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row RPEventType
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.RPEventTypes = append(r.RPEventTypes, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.RPEventTypes) == 0 {
		r.RPEventTypes = []RPEventType{}
	}
	return nil
}

// Delete removes RPEvenType whose ID is given from database
func (r *RPEventType) Delete(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM rp_event_type WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Type d'événement introuvable")
	}
	return nil
}
