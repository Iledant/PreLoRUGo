package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// CoproForecast model
type CoproForecast struct {
	ID             int64      `json:"ID"`
	CommissionID   int64      `json:"CommissionID"`
	CommissionDate NullTime   `json:"CommissionDate"`
	CommissionName string     `json:"CommissionName"`
	Value          int64      `json:"Value"`
	Comment        NullString `json:"Comment"`
	CoproID        int64      `json:"CoproID"`
	ActionID       int64      `json:"ActionID"`
	ActionCode     int64      `json:"ActionCode"`
	ActionName     string     `json:"ActionName"`
}

// CoproForecasts embeddes an array of CoproForecast for json export
type CoproForecasts struct {
	CoproForecasts []CoproForecast `json:"CoproForecast"`
}

// CoproForecastLine is used to decode a line of CoproForecast batch
type CoproForecastLine struct {
	ID           int64      `json:"ID"`
	CommissionID int64      `json:"CommissionID"`
	Value        int64      `json:"Value"`
	Comment      NullString `json:"Comment"`
	CoproID      int64      `json:"CoproID"`
	ActionCode   int64      `json:"ActionCode"`
}

// CoproForecastBatch embeddes an array of CoproForecastLine for json export
type CoproForecastBatch struct {
	Lines []CoproForecastLine `json:"CoproForecast"`
}

// Validate checks if CoproForecast's fields are correctly filled
func (r *CoproForecast) Validate() error {
	if r.CommissionID == 0 || r.Value == 0 || r.CoproID == 0 || r.ActionID == 0 {
		return errors.New("Champ incorrect")
	}
	return nil
}

// Create insert a new CoproForecast into database
func (r *CoproForecast) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO copro_forecast 
	(commission_id,value,comment,copro_id,action_id)
 VALUES($1,$2,$3,$4,$5) RETURNING id`, &r.CommissionID, &r.Value, &r.Comment,
		&r.CoproID, &r.ActionID).Scan(&r.ID)
	if err != nil {
		return err
	}
	err = db.QueryRow(`SELECT c.name, c.date, b.code, b.name 
	FROM commission c, budget_action b WHERE c.id=$1 AND b.id=$2`,
		r.CommissionID, r.ActionID).Scan(&r.CommissionName, &r.CommissionDate,
		&r.ActionCode, &r.ActionName)
	return err
}

// Get fetches a CoproForecast from database using ID field
func (r *CoproForecast) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT cf.commission_id,c.date,c.name, 
		cf.value,cf.comment,cf.copro_id, b.code, b.name
	FROM copro_forecast cf
	JOIN commission c ON c.id=cf.commission_id
	JOIN budget_action b ON b.id = cf.action_id WHERE cf.ID=$1`, r.ID).
		Scan(&r.CommissionID, &r.CommissionDate, &r.CommissionName,
			&r.Value, &r.Comment, &r.CoproID, &r.ActionCode, &r.ActionName)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a CoproForecast in database
func (r *CoproForecast) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE copro_forecast SET commission_id=$1,value=$2,
	comment=$3,copro_id=$4,action_id=$5 WHERE id=$6`,
		r.CommissionID, r.Value, r.Comment, r.CoproID, r.ActionID, r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Prévision copro introuvable")
	}
	err = db.QueryRow(`SELECT c.name, c.date, b.code, b.name 
	FROM commission c, budget_action b WHERE c.id=$1 AND b.id=$2`, r.CommissionID,
		r.ActionID).Scan(&r.CommissionName, &r.CommissionDate, &r.ActionCode, &r.ActionName)
	return err
}

// GetAll fetches all CoproForecasts from database
func (r *CoproForecasts) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT cf.id, cf.commission_id,c.date,c.name, 
		cf.value,cf.comment,cf.copro_id, b.code, b.name
	FROM copro_forecast cf
	JOIN commission c ON c.id=cf.commission_id
	JOIN budget_action b ON b.id = cf.action_id`)
	if err != nil {
		return err
	}
	var row CoproForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Comment, &row.CoproID,
			&row.ActionCode, &row.ActionName); err != nil {
			return err
		}
		r.CoproForecasts = append(r.CoproForecasts, row)
	}
	err = rows.Err()
	if len(r.CoproForecasts) == 0 {
		r.CoproForecasts = []CoproForecast{}
	}
	return err
}

// Get fetches all copro linked CoproForecasts from database
func (r *CoproForecasts) Get(ID int64, db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT cf.id, cf.commission_id,c.date,c.name, 
		cf.value,cf.comment,cf.copro_id, b.id, b.code, b.name
	FROM copro_forecast cf
	JOIN commission c ON c.id=cf.commission_id
	JOIN budget_action b ON b.id = cf.action_id
	WHERE cf.copro_id=$1`, ID)
	if err != nil {
		return err
	}
	var row CoproForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Comment, &row.CoproID, &row.ActionID,
			&row.ActionCode, &row.ActionName); err != nil {
			return err
		}
		r.CoproForecasts = append(r.CoproForecasts, row)
	}
	err = rows.Err()
	if len(r.CoproForecasts) == 0 {
		r.CoproForecasts = []CoproForecast{}
	}
	return err
}

// Delete removes CoproForecast whose ID is given from database
func (r *CoproForecast) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM copro_forecast WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Prévision copro introuvable")
	}
	return nil
}

// Save insert a batch of CoproForecastLine into database
func (r *CoproForecastBatch) Save(db *sql.DB) (err error) {
	for i, r := range r.Lines {
		if r.CommissionID == 0 || r.Value == 0 || r.CoproID == 0 || r.ActionCode == 0 {
			return fmt.Errorf("ligne %d, champs incorrects", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_copro_forecast", "id", "commission_id",
		"value", "comment", "copro_id", "action_code"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range r.Lines {
		if _, err = stmt.Exec(r.ID, r.CommissionID, r.Value, r.Comment, r.CoproID,
			r.ActionCode); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{`UPDATE copro_forecast SET commission_id=t.commission_id,
	value=t.value,comment=t.comment,copro_id=t.copro_id, action_id=b.id
	FROM temp_copro_forecast t JOIN budget_action b ON t.action_code = b.code
	WHERE t.id = copro_forecast.id`,
		`INSERT INTO copro_forecast (commission_id,value,comment,copro_id,action_id)
	SELECT t.commission_id,t.value,t.comment,t.copro_id,b.id from temp_copro_forecast t
		JOIN budget_action b ON t.action_code = b.code
		WHERE t.id NOT IN (SELECT id from copro_forecast)`,
		`DELETE from temp_copro_forecast`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d, %v", i, err)
		}
	}
	tx.Commit()
	return nil
}
