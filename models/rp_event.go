package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// RPEvent model for storing events linked to a renew project and a event type
type RPEvent struct {
	ID             int64      `json:"ID"`
	RenewProjectID int64      `json:"RenewProjectID"`
	RPEventTypeID  int64      `json:"RPEventTypeID"`
	Date           time.Time  `json:"Date"`
	Comment        NullString `json:"Comment"`
}

// RPEvents embeddes an array of RPEvent for json export
type RPEvents struct {
	RPEvents []RPEvent `json:"RPEvent"`
}

// FullRPEvent is used to fetch all events linked to a renew project with full
// event type name
type FullRPEvent struct {
	ID              int64     `json:"ID"`
	RPEventTypeID   int64     `json:"RPEventTypeID"`
	RPEventTypeName string    `json:"RPEventTypeName"`
	Date            time.Time `json:"Date"`
	Comment         string    `json:"Comment"`
}

// FullRPEvents embeddes an array of FullRPEvent for json export
type FullRPEvents struct {
	FullRPEvents []FullRPEvent `json:"FullRPEvent"`
}

// Create insert a new RPEvent into database
func (r *RPEvent) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO rp_event (renew_project_id,rp_event_type_id,date,
		comment) VALUES($1,$2,$3,$4) RETURNING id`, r.RenewProjectID, r.RPEventTypeID,
		r.Date, r.Comment).Scan(&r.ID)
	return err
}

// Validate check if fields complies with database constraints
func (r *RPEvent) Validate() error {
	if r.RenewProjectID == 0 || r.RPEventTypeID == 0 {
		return fmt.Errorf("Champ RenewProjectID ou RPEventTypeID vide")
	}
	return nil
}

// Get fetches a RPEvent from database using ID field
func (r *RPEvent) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT renew_project_id,rp_event_type_id,date,comment
	FROM rp_event WHERE ID=$1`, r.ID).
		Scan(&r.RenewProjectID, &r.RPEventTypeID, &r.Date, &r.Comment)
	return err
}

// Update modifies a RPEvent in database
func (r *RPEvent) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE rp_event SET renew_project_id=$1,rp_event_type_id=$2,
	date=$3,comment=$4 WHERE id=$5`, r.RenewProjectID, r.RPEventTypeID, r.Date,
		r.Comment, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Événement introuvable")
	}
	return nil
}

// GetAll fetches all RPEvent from database
func (r *RPEvents) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,renew_project_id,rp_event_type_id,date,comment
	 FROM rp_event`)
	if err != nil {
		return err
	}
	var row RPEvent
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.RenewProjectID, &row.RPEventTypeID, &row.Date,
			&row.Comment); err != nil {
			return err
		}
		r.RPEvents = append(r.RPEvents, row)
	}
	err = rows.Err()
	if len(r.RPEvents) == 0 {
		r.RPEvents = []RPEvent{}
	}
	return err
}

// Delete removes RPEvenType whose ID is given from database
func (r *RPEvent) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM rp_event WHERE id = $1", r.ID)
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
func (f *FullRPEvents) GetLinked(db *sql.DB, rpID int64) error {
	rows, err := db.Query(`SELECT r.id,r.rp_event_type_id,rt.name,r.date,r.comment
	 FROM rp_event r
	 JOIN rp_event_type rt ON r.rp_event_type_id=rt.id
	  WHERE renew_project_id=$1`, rpID)
	if err != nil {
		return err
	}
	var row FullRPEvent
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.RPEventTypeID, &row.RPEventTypeName, &row.Date,
			&row.Comment); err != nil {
			return err
		}
		f.FullRPEvents = append(f.FullRPEvents, row)
	}
	err = rows.Err()
	if len(f.FullRPEvents) == 0 {
		f.FullRPEvents = []FullRPEvent{}
	}
	return err
}
