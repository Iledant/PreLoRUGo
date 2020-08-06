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
	Name           string     `json:"Name"`
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
	ID            int64     `json:"ID"`
	RPEventTypeID int64     `json:"RPEventTypeID"`
	Name          string    `json:"Name"`
	Date          time.Time `json:"Date"`
	Comment       string    `json:"Comment"`
}

// FullRPEvents embeddes an array of FullRPEvent for json export
type FullRPEvents struct {
	FullRPEvents []FullRPEvent `json:"FullRPEvent"`
}

// Create insert a new RPEvent into database
func (r *RPEvent) Create(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO rp_event (renew_project_id,rp_event_type_id,date,
		comment) VALUES($1,$2,$3,$4) RETURNING id`, r.RenewProjectID, r.RPEventTypeID,
		r.Date, r.Comment).Scan(&r.ID)
	if err != nil {
		return fmt.Errorf("insert %v", err)
	}
	return db.QueryRow(`SELECT name FROM rp_event_type WHERE id=$1`,
		r.RPEventTypeID).Scan(&r.Name)
}

// Validate check if fields complies with database constraints
func (r *RPEvent) Validate() error {
	if r.RenewProjectID == 0 {
		return fmt.Errorf("Champ RenewProjectID vide")
	}
	if r.RPEventTypeID == 0 {
		return fmt.Errorf("Champ RPEventTypeID vide")
	}
	return nil
}

// Get fetches a RPEvent from database using ID field
func (r *RPEvent) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT r.renew_project_id,r.rp_event_type_id,rt.name,
		r.date,r.comment
	FROM rp_event r
	JOIN rp_event_type rt ON rt.id=r.rp_event_type_id
	WHERE r.id=$1`, r.ID).Scan(&r.RenewProjectID, &r.RPEventTypeID,
		&r.Name, &r.Date, &r.Comment)
}

// Update modifies a RPEvent in database
func (r *RPEvent) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE rp_event SET renew_project_id=$1,rp_event_type_id=$2,
	date=$3,comment=$4 WHERE id=$5`, r.RenewProjectID, r.RPEventTypeID, r.Date,
		r.Comment, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Événement introuvable")
	}
	return db.QueryRow(`SELECT name FROM rp_event_type WHERE id=$1`,
		r.RPEventTypeID).Scan(&r.Name)
}

// GetAll fetches all RPEvent from database
func (r *RPEvents) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT r.id,r.renew_project_id,r.rp_event_type_id,
	 rt.name,r.date,r.comment
	 FROM rp_event r JOIN rp_event_type rt ON rt.id=r.rp_event_type_id`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row RPEvent
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.RenewProjectID, &row.RPEventTypeID,
			&row.Name, &row.Date, &row.Comment); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.RPEvents = append(r.RPEvents, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.RPEvents) == 0 {
		r.RPEvents = []RPEvent{}
	}
	return nil
}

// Delete removes RPEvenType whose ID is given from database
func (r *RPEvent) Delete(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM rp_event WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
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
		return fmt.Errorf("select %v", err)
	}
	var row FullRPEvent
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.RPEventTypeID, &row.Name, &row.Date,
			&row.Comment); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		f.FullRPEvents = append(f.FullRPEvents, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(f.FullRPEvents) == 0 {
		f.FullRPEvents = []FullRPEvent{}
	}
	return nil
}
