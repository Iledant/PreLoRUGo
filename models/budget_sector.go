package models

import (
	"database/sql"
	"errors"
)

// BudgetSector model
type BudgetSector struct {
	ID       int64      `json:"ID"`
	Name     string     `json:"Name"`
	FullName NullString `json:"FullName"`
}

// BudgetSectors embeddes an array of BudgetSector for json export
type BudgetSectors struct {
	BudgetSectors []BudgetSector `json:"BudgetSector"`
}

// Validate checks if BudgetSector's fields are correctly filled
func (b *BudgetSector) Validate() error {
	if b.Name == "" {
		return errors.New("Champ incorrect")
	}
	return nil
}

// Create insert a new BudgetSector into database
func (b *BudgetSector) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO budget_sector (name,full_name)
 VALUES($1,$2) RETURNING id`, &b.Name, &b.FullName).Scan(&b.ID)
	return err
}

// Get fetches a BudgetSector from database using ID field
func (b *BudgetSector) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT name, full_name FROM budget_sector WHERE ID=$1`, b.ID).
		Scan(&b.Name, &b.FullName)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a budget_sector in database
func (b *BudgetSector) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE budget_sector SET name=$1,full_name=$2 WHERE id=$3`,
		b.Name, b.FullName, b.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("ID introuvable")
	}
	return err
}

// GetAll fetches all BudgetSectors from database
func (b *BudgetSectors) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,name,full_name FROM budget_sector`)
	if err != nil {
		return err
	}
	var row BudgetSector
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Name, &row.FullName); err != nil {
			return err
		}
		b.BudgetSectors = append(b.BudgetSectors, row)
	}
	err = rows.Err()
	if len(b.BudgetSectors) == 0 {
		b.BudgetSectors = []BudgetSector{}
	}
	return err
}

// Delete removes budget_sector whose ID is given from database
func (b *BudgetSector) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM budget_sector WHERE id = $1", b.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("ID introuvable")
	}
	return nil
}
