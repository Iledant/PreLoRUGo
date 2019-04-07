package models

import (
	"database/sql"
	"errors"
)

// RenewProject model
type RenewProject struct {
	ID             int64     `json:"ID"`
	Reference      string    `json:"Reference"`
	Name           string    `json:"Name"`
	Budget         int64     `json:"Budget"`
	Population     NullInt64 `json:"Population"`
	CompositeIndex NullInt64 `json:"CompositeIndex"`
}

// RenewProjects embeddes an array of RenewProject for json export
type RenewProjects struct {
	RenewProjects []RenewProject `json:"RenewProject"`
}

// Validate checks if the fields of a renew project are correctly filled
func (r *RenewProject) Validate() error {
	if r.Reference == "" || r.Name == "" || r.Budget == 0 {
		return errors.New("Champ reference, name ou budget incorrect")
	}
	return nil
}

// Create insert a renew project into database returning it's ID
func (r *RenewProject) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO renew_project (reference,name,budget,population,
		composite_index) VALUES($1,$2,$3,$4,$5) RETURNING id`, r.Reference, r.Name,
		r.Budget, r.Population, r.CompositeIndex).Scan(&r.ID)
	return err
}

// Update modifies a renew program into database
func (r *RenewProject) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE renew_project SET reference=$1, name=$2, budget=$3,
	population=$4, composite_index=$5 WHERE id = $6`, r.Reference, r.Name, r.Budget,
		r.Population, r.CompositeIndex, r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Projet de renouvellement introuvable")
	}
	return err
}

// GetAll fetches all renew projects from database
func (r *RenewProjects) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,reference,name,budget,population,composite_index
	 FROM renew_project`)
	if err != nil {
		return err
	}
	var row RenewProject
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Reference, &row.Name, &row.Budget, &row.Population, &row.CompositeIndex); err != nil {
			return err
		}
		r.RenewProjects = append(r.RenewProjects, row)
	}
	err = rows.Err()
	if len(r.RenewProjects) == 0 {
		r.RenewProjects = []RenewProject{}
	}
	return err
}

// Delete removes a renew project whose ID is given from database
func (r *RenewProject) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM renew_project WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Projet de renouvellement introuvable")
	}
	return nil
}
