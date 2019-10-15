package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CoproEvent model for storing events linked to a renew project and a event type
type CoproEvent struct {
	ID               int64      `json:"ID"`
	CoproID          int64      `json:"CoproID"`
	CoproEventTypeID int64      `json:"CoproEventTypeID"`
	Name             string     `json:"Name"`
	Date             time.Time  `json:"Date"`
	Comment          NullString `json:"Comment"`
}

// CoproEvents embeddes an array of CoproEvent for json export
type CoproEvents struct {
	CoproEvents []CoproEvent `json:"CoproEvent"`
}

// FullCoproEvent is used to fetch all events linked to a renew project with full
// event type name
type FullCoproEvent struct {
	ID               int64     `json:"ID"`
	CoproEventTypeID int64     `json:"CoproEventTypeID"`
	Name             string    `json:"Name"`
	Date             time.Time `json:"Date"`
	Comment          string    `json:"Comment"`
}

// FullCoproEvents embeddes an array of FullCoproEvent for json export
type FullCoproEvents struct {
	FullCoproEvents []FullCoproEvent `json:"FullCoproEvent"`
}

// Create insert a new CoproEvent into database
func (r *CoproEvent) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO copro_event (copro_id,copro_event_type_id,date,
		comment) VALUES($1,$2,$3,$4) RETURNING id`, r.CoproID, r.CoproEventTypeID,
		r.Date, r.Comment).Scan(&r.ID)
	if err != nil {
		return err
	}
	err = db.QueryRow(`SELECT name FROM copro_event_type WHERE id=$1`,
		r.CoproEventTypeID).Scan(&r.Name)
	return err
}

// Validate check if fields complies with database constraints
func (r *CoproEvent) Validate() error {
	if r.CoproID == 0 {
		return fmt.Errorf("CoproID vide")
	}
	if r.CoproEventTypeID == 0 {
		return fmt.Errorf("CoproEventTypeID vide")
	}
	return nil
}

// Get fetches a CoproEvent from database using ID field
func (r *CoproEvent) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT r.copro_id,r.copro_event_type_id,rt.name,
		r.date,r.comment
	FROM copro_event r
	JOIN copro_event_type rt ON rt.id=r.copro_event_type_id
	WHERE r.id=$1`, r.ID).Scan(&r.CoproID, &r.CoproEventTypeID,
		&r.Name, &r.Date, &r.Comment)
	return err
}

// Update modifies a CoproEvent in database
func (r *CoproEvent) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE copro_event SET copro_id=$1,copro_event_type_id=$2,
	date=$3,comment=$4 WHERE id=$5`, r.CoproID, r.CoproEventTypeID, r.Date,
		r.Comment, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Événement introuvable")
	}
	err = db.QueryRow(`SELECT name FROM copro_event_type WHERE id=$1`,
		r.CoproEventTypeID).Scan(&r.Name)
	return err
}

// GetAll fetches all CoproEvent from database
func (r *CoproEvents) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT r.id,r.copro_id,r.copro_event_type_id,
	 rt.name,r.date,r.comment
	 FROM copro_event r JOIN copro_event_type rt ON rt.id=r.copro_event_type_id`)
	if err != nil {
		return err
	}
	var row CoproEvent
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CoproID, &row.CoproEventTypeID,
			&row.Name, &row.Date, &row.Comment); err != nil {
			return err
		}
		r.CoproEvents = append(r.CoproEvents, row)
	}
	err = rows.Err()
	if len(r.CoproEvents) == 0 {
		r.CoproEvents = []CoproEvent{}
	}
	return err
}

// Delete removes CoproEvenType whose ID is given from database
func (r *CoproEvent) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM copro_event WHERE id=$1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Événement introuvable")
	}
	return nil
}

// GetLinked fetches all events linked to a renew project with their full name
func (f *FullCoproEvents) GetLinked(db *sql.DB, rpID int64) error {
	rows, err := db.Query(`SELECT r.id,r.copro_event_type_id,rt.name,r.date,r.comment
	 FROM copro_event r
	 JOIN copro_event_type rt ON r.copro_event_type_id=rt.id
	  WHERE copro_id=$1`, rpID)
	if err != nil {
		return err
	}
	var row FullCoproEvent
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CoproEventTypeID, &row.Name, &row.Date,
			&row.Comment); err != nil {
			return err
		}
		f.FullCoproEvents = append(f.FullCoproEvents, row)
	}
	err = rows.Err()
	if len(f.FullCoproEvents) == 0 {
		f.FullCoproEvents = []FullCoproEvent{}
	}
	return err
}
