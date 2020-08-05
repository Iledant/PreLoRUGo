package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// HousingForecast model
type HousingForecast struct {
	ID             int64      `json:"ID"`
	CommissionID   int64      `json:"CommissionID"`
	CommissionDate NullTime   `json:"CommissionDate"`
	CommissionName string     `json:"CommissionName"`
	Value          int64      `json:"Value"`
	Comment        NullString `json:"Comment"`
	ActionID       int64      `json:"ActionID"`
	ActionName     string     `json:"ActionName"`
}

// HousingForecasts embeddes an array of HousingForecast for json export
type HousingForecasts struct {
	HousingForecasts []HousingForecast `json:"HousingForecast"`
}

// HousingForecastLine is used to decode a line of HousingForecast batch
type HousingForecastLine struct {
	ID           int64      `json:"ID"`
	CommissionID int64      `json:"CommissionID"`
	Value        int64      `json:"Value"`
	Comment      NullString `json:"Comment"`
	ActionID     int64      `json:"ActionID"`
}

// HousingForecastBatch embeddes an array of HousingForecastLine for json export
type HousingForecastBatch struct {
	Lines []HousingForecastLine `json:"HousingForecast"`
}

// Validate checks if HousingForecast's fields are correctly filled
func (r *HousingForecast) Validate() error {
	if r.CommissionID == 0 {
		return errors.New("Champ CommissionID incorrect")
	}
	if r.Value == 0 {
		return errors.New("Champ Value incorrect")
	}
	if r.ActionID == 0 {
		return errors.New("Champ ActionID incorrect")
	}
	return nil
}

// Create insert a new HousingForecast into database
func (r *HousingForecast) Create(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO housing_forecast 
	(commission_id,value,comment,action_id)
 VALUES($1,$2,$3,$4) RETURNING id`, &r.CommissionID, &r.Value, &r.Comment,
		&r.ActionID).Scan(&r.ID)
	if err != nil {
		return fmt.Errorf("insert %v", err)
	}
	return db.QueryRow(`SELECT c.name,c.date,b.name FROM housing_forecast h
	JOIN commission c ON c.id=h.commission_id
	JOIN budget_action b ON b.id=h.action_id
	WHERE h.id=$1`, r.ID).Scan(&r.CommissionName, &r.CommissionDate, &r.ActionName)
}

// Get fetches a HousingForecast from database using ID field
func (r *HousingForecast) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT hf.commission_id,c.date,c.name, 
		hf.value,hf.comment,hf.action_id,b.name
	FROM housing_forecast hf
	JOIN commission c ON c.id=hf.commission_id
	JOIN budget_action b ON b.id=hf.action_id
	WHERE hf.ID=$1`, r.ID).Scan(&r.CommissionID, &r.CommissionDate,
		&r.CommissionName, &r.Value, &r.Comment, &r.ActionID, &r.ActionName)
}

// Update modifies a HousingForecast in database
func (r *HousingForecast) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE housing_forecast SET commission_id=$1,value=$2,
	comment=$3,action_id=$4 WHERE id=$5`, r.CommissionID, r.Value, r.Comment,
		r.ActionID, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Prévision housing introuvable")
	}
	return db.QueryRow(`SELECT c.name,c.date,b.name FROM housing_forecast h
	JOIN commission c ON c.id=h.commission_id
	JOIN budget_action b ON b.id=h.action_id
	WHERE h.id=$1`, r.ID).Scan(&r.CommissionName, &r.CommissionDate, &r.ActionName)
}

// GetAll fetches all HousingForecasts from database
func (r *HousingForecasts) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT hf.ID, hf.commission_id,c.date,c.name, 
		hf.value,hf.comment,hf.action_id,b.name 
	FROM housing_forecast hf
	JOIN commission c ON c.id=hf.commission_id
	JOIN budget_action b ON b.id=hf.action_id`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row HousingForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Comment, &row.ActionID,
			&row.ActionName); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.HousingForecasts = append(r.HousingForecasts, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.HousingForecasts) == 0 {
		r.HousingForecasts = []HousingForecast{}
	}
	return err
}

// Get fetches all housing linked HousingForecasts from database
func (r *HousingForecasts) Get(ID int64, db *sql.DB) error {
	rows, err := db.Query(`SELECT hf.ID, hf.commission_id,c.date,c.name, 
		hf.value,hf.comment,hf.action_id,b.name 
	FROM housing_forecast hf
	JOIN commission c ON c.id=hf.commission_id
	JOIN budget_action b ON b.id=hf.action_id
	WHERE hf.action_id=$1`, ID)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row HousingForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Comment, &row.ActionID,
			&row.ActionName); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.HousingForecasts = append(r.HousingForecasts, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.HousingForecasts) == 0 {
		r.HousingForecasts = []HousingForecast{}
	}
	return nil
}

// Delete removes HousingForecast whose ID is given from database
func (r *HousingForecast) Delete(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM housing_forecast WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Prévision housing introuvable")
	}
	return nil
}

// Save insert a batch of HousingForecastLine into database
func (r *HousingForecastBatch) Save(db *sql.DB) error {
	for i, r := range r.Lines {
		if r.CommissionID == 0 {
			return fmt.Errorf("ligne %d, CommissionID nul", i+1)
		}
		if r.Value == 0 {
			return fmt.Errorf("ligne %d, Value nul", i+1)
		}
		if r.ActionID == 0 {
			return fmt.Errorf("ligne %d, ActionID nul", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("début de transaction %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_housing_forecast", "id", "commission_id",
		"value", "comment", "action_id"))
	if err != nil {
		return fmt.Errorf("insert statement %v", err)
	}
	defer stmt.Close()
	for _, r := range r.Lines {
		if _, err = stmt.Exec(r.ID, r.CommissionID, r.Value, r.Comment, r.ActionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement execution %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{`UPDATE housing_forecast SET commission_id=t.commission_id,
	value=t.value,comment=t.comment,action_id=t.action_id 
	FROM temp_housing_forecast t WHERE t.id = housing_forecast.id`,
		`INSERT INTO housing_forecast (commission_id,value,comment,action_id)
	SELECT commission_id,value,comment,action_id from temp_housing_forecast 
		WHERE id NOT IN (SELECT id from housing_forecast)`,
		`DELETE from temp_housing_forecast`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	return tx.Commit()
}
