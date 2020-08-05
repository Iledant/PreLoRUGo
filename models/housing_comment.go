package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// HousingComment is used for normalizing housing transfer
type HousingComment struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// HousingComments embeddes an array of housing transfers for json export
type HousingComments struct {
	HousingComments []HousingComment `json:"HousingComment"`
}

// Create insert a new housing typology into database
func (r *HousingComment) Create(db *sql.DB) (err error) {
	return db.QueryRow(`INSERT INTO housing_comment(name) VALUES($1) RETURNING id`,
		&r.Name).Scan(&r.ID)
}

// Valid checks if fields complies with database constraints
func (r *HousingComment) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("Nom vide")
	}
	return nil
}

// Get fetches a HousingComment from database using ID field
func (r *HousingComment) Get(db *sql.DB) (err error) {
	return db.QueryRow(`SELECT name FROM housing_comment WHERE ID=$1`, r.ID).
		Scan(&r.Name)
}

// Update modifies a HousingComment in database
func (r *HousingComment) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing_comment SET name=$1 WHERE id=$2`,
		r.Name, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	if count != 1 {
		return errors.New("Commentaire introuvable")
	}
	return nil
}

// GetAll fetches all housing comments from database
func (r *HousingComments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name FROM housing_comment`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row HousingComment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.HousingComments = append(r.HousingComments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)

	}
	if len(r.HousingComments) == 0 {
		r.HousingComments = []HousingComment{}
	}
	return err
}

// Delete removes housing comment whose ID is given from database
func (r *HousingComment) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM housing_comment WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Commentaire introuvable")
	}
	return nil
}
