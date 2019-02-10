package models

import "database/sql"

//BudgetAction model
type BudgetAction struct {
	ID       int       `json:"ID"`
	Code     string    `json:"Code"`
	Name     string    `json:"Name"`
	SectorID NullInt64 `json:"SectorID"`
}

// BudgetActions embeddes an array of BudgetAction for json export
type BudgetActions struct {
	BudgetActions []BudgetAction `json:"BudgetAction"`
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
