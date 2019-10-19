package models

import (
	"database/sql"
	"fmt"
)

// CoproDoc model
type CoproDoc struct {
	ID      int64  `json:"ID"`
	CoproID int64  `json:"CoproID"`
	Name    string `json:"Name"`
	Link    string `json:"Link"`
}

// CoproDocs embeddes an array of CoproDoc for json export
type CoproDocs struct {
	Lines []CoproDoc `json:"CoproDoc"`
}

// validate check if fields respect database constraint
func (c *CoproDoc) validate() error {
	if c.Name == "" {
		return fmt.Errorf("nom vide")
	}
	if c.Link == "" {
		return fmt.Errorf("link vide")
	}
	return nil
}

// GetAll fetches all documents linked to a copro whose ID is given
func (c *CoproDocs) GetAll(CoproID int64, db *sql.DB) error {
	rows, err := db.Query(`SELECT id,name,link FROM copro_doc WHERE copro_id=$1`, CoproID)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	defer rows.Close()
	var line CoproDoc
	for rows.Next() {
		if err = rows.Scan(&line.ID, &line.Name, &line.Link); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		c.Lines = append(c.Lines, line)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Err %v", err)
	}
	if len(c.Lines) == 0 {
		c.Lines = []CoproDoc{}
	}
	return nil
}

// Save insert a CoproDoc into database
func (c *CoproDoc) Save(db *sql.DB) error {
	if err := c.validate(); err != nil {
		return err
	}
	if err := db.QueryRow(`INSERT INTO copro_doc (copro_id,name,link) 
	VALUES($1,$2,$3) RETURNING id`, c.CoproID, c.Name, c.Link).Scan(&c.ID); err != nil {
		return fmt.Errorf("insert %v", err)
	}
	return nil
}

// Update modifies the copro doc in the database
func (c *CoproDoc) Update(db *sql.DB) error {
	if err := c.validate(); err != nil {
		return err
	}
	res, err := db.Exec(`UPDATE copro_doc SET name=$1, link=$2 WHERE id=$3`,
		c.Name, c.Link, c.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected %v", err)
	}
	if count != 1 {
		return fmt.Errorf("copro introuvable")
	}
	return nil
}

// Delete removes a copro_doc from database
func (c *CoproDoc) Delete(db *sql.DB) error {
	res, err := db.Exec(`DELETE FROM copro_doc WHERE id=$1`, c.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected %v", err)
	}
	if count != 1 {
		return fmt.Errorf("copro introuvable")
	}
	return nil
}
