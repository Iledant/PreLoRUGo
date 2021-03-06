package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// Department model
type Department struct {
	ID   int64  `json:"ID"`
	Code int64  `json:"Code"`
	Name string `json:"Name"`
}

// Departments embeddes an array of Department for json export
type Departments struct {
	Departments []Department `json:"Department"`
}

// Validate checks if Department's fields are correctly filled
func (c *Department) Validate() error {
	if c.Code == 0 {
		return errors.New("Champ code incorrect")
	}
	if c.Name == "" {
		return errors.New("Champ name incorrect")
	}
	return nil
}

// Create insert a new Department into database
func (c *Department) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO department (code,name)
 VALUES($1,$2) RETURNING id`, &c.Code, &c.Name).Scan(&c.ID)
}

// Get fetches a Department from database using ID field
func (c *Department) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT code, name FROM department WHERE ID=$1`, c.ID).
		Scan(&c.Code, &c.Name)
}

// Update modifies a department in database
func (c *Department) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE department SET code=$1,name=$2 WHERE id=$3`,
		c.Code, c.Name, c.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Département introuvable")
	}
	return nil
}

// GetAll fetches all Departments from database
func (c *Departments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,code,name FROM department`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row Department
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Code, &row.Name); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		c.Departments = append(c.Departments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(c.Departments) == 0 {
		c.Departments = []Department{}
	}
	return err
}

// Delete removes department whose ID is given from database
func (c *Department) Delete(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}
	if _, err := tx.Exec(`UPDATE community SET department_id=NULL 
		WHERE department_id=$1`, c.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("update %v", err)
	}
	res, err := tx.Exec("DELETE FROM department WHERE id = $1", c.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		tx.Rollback()
		return errors.New("Département introuvable")
	}
	return tx.Commit()
}
