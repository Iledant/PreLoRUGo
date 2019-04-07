package models

import (
	"database/sql"
	"errors"
)

// Commission model
type Commission struct {
	ID   int64    `json:"ID"`
	Name string   `json:"Name"`
	Date NullTime `json:"Date"`
}

// Commissions embeddes an array of Commission for json export
type Commissions struct {
	Commissions []Commission `json:"Commission"`
}

// Validate checks if Commission's fields are correctly filled
func (c *Commission) Validate() error {
	if c.Name == "" {
		return errors.New("Champ name vide")
	}
	return nil
}

// Create insert a new Commission into database
func (c *Commission) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO commission (name,date)
 VALUES($1,$2) RETURNING id`, &c.Name, &c.Date).Scan(&c.ID)
	return err
}

// Get fetches a Commission from database using ID field
func (c *Commission) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT name, date FROM commission WHERE ID=$1`, c.ID).
		Scan(&c.Name, &c.Date)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a commission in database
func (c *Commission) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE commission SET name=$1,date=$2 WHERE id=$3`,
		c.Name, c.Date, c.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Commission introuvable")
	}
	return err
}

// GetAll fetches all Commissions from database
func (c *Commissions) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name,date FROM commission`)
	if err != nil {
		return err
	}
	var row Commission
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name, &row.Date); err != nil {
			return err
		}
		c.Commissions = append(c.Commissions, row)
	}
	err = rows.Err()
	if len(c.Commissions) == 0 {
		c.Commissions = []Commission{}
	}
	return err
}

// Delete removes commission whose ID is given from database
func (c *Commission) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM commission WHERE id = $1", c.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Commission introuvable")
	}
	return nil
}
