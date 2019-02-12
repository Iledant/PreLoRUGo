package models

import (
	"database/sql"
	"errors"
)

//BudgetAction model
type BudgetAction struct {
	ID       int64     `json:"ID"`
	Code     string    `json:"Code"`
	Name     string    `json:"Name"`
	SectorID NullInt64 `json:"SectorID"`
}

// BudgetActions embeddes an array of BudgetAction for json export
type BudgetActions struct {
	BudgetActions []BudgetAction `json:"BudgetAction"`
}

// Validate checks if the fields of a budget action are correctly filled
func (b *BudgetAction) Validate() error {
	if b.Code == "" || b.Name == "" {
		return errors.New("Champ code ou name incorrect")
	}
	return nil
}

// Create insert a new budget_action into database
func (b *BudgetAction) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO budget_action (code,name,sector_id) 
	VALUES($1,$2,$3) RETURNING id`, b.Code, b.Name, b.SectorID).Scan(&b.ID)
	return err
}

// Update modifies the database entry with sent BudgetAction using it's ID
func (b *BudgetAction) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE budget_action SET code=$1,name=$2,sector_id=$3
	 WHERE id = $4`, b.Code, b.Name, b.SectorID, b.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Action budgétaire introuvable")
	}
	return err
}

// GetAll fetches all BudgetActions from database
func (b *BudgetActions) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,code,name,sector_id FROM budget_action`)
	if err != nil {
		return err
	}
	var r BudgetAction
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ID, &r.Code, &r.Name, &r.SectorID); err != nil {
			return err
		}
		b.BudgetActions = append(b.BudgetActions, r)
	}
	err = rows.Err()
	return err
}

// Delete removes a budget action from database using ID field
func (b *BudgetAction) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM budget_action WHERE id = $1", b.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Action budgétaire introuvable")
	}
	return nil
}
