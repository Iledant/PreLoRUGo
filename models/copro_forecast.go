package models

import (
	"database/sql"
	"errors"
	"fmt"
)

// CoproForecast model
type CoproForecast struct {
	ID           int64      `json:"ID"`
	CommissionID int64      `json:"CommissionID"`
	Value        int64      `json:"Value"`
	Comment      NullString `json:"Comment"`
	CoproID      int64      `json:"CoproID"`
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
}

// CoproForecastBatch embeddes an array of CoproForecastLine for json export
type CoproForecastBatch struct {
	Lines []CoproForecastLine `json:"CoproForecast"`
}

// Validate checks if CoproForecast's fields are correctly filled
func (r *CoproForecast) Validate() error {
	if r.CommissionID == 0 || r.Value == 0 || r.CoproID == 0 {
		return errors.New("Champ incorrect")
	}
	return nil
}

// Create insert a new CoproForecast into database
func (r *CoproForecast) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO copro_forecast 
	(commission_id,value,comment,copro_id)
 VALUES($1,$2,$3,$4) RETURNING id`, &r.CommissionID, &r.Value, &r.Comment,
		&r.CoproID).Scan(&r.ID)
	return err
}

// Get fetches a CoproForecast from database using ID field
func (r *CoproForecast) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT commission_id, value, comment, copro_id 
	FROM copro_forecast WHERE ID=$1`, r.ID).Scan(&r.CommissionID, &r.Value,
		&r.Comment, &r.CoproID)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a CoproForecast in database
func (r *CoproForecast) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE copro_forecast SET commission_id=$1,value=$2,
	comment=$3,copro_id=$4 WHERE id=$5`,
		r.CommissionID, r.Value, r.Comment, r.CoproID, r.ID)
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
	return err
}

// GetAll fetches all CoproForecasts from database
func (r *CoproForecasts) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,commission_id,value,comment,copro_id 
	FROM copro_forecast`)
	if err != nil {
		return err
	}
	var row CoproForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.Value, &row.Comment,
			&row.CoproID); err != nil {
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
	rows, err := db.Query(`SELECT id,commission_id,value,comment,copro_id 
	FROM copro_forecast WHERE copro_id=$1`, ID)
	if err != nil {
		return err
	}
	var row CoproForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.Value, &row.Comment,
			&row.CoproID); err != nil {
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
		if r.CommissionID == 0 || r.Value == 0 || r.CoproID == 0 {
			return fmt.Errorf("ligne %d, champs incorrects", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_copro_forecast 
	(id, commission_id,value,comment,copro_id) VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range r.Lines {
		if _, err = stmt.Exec(r.ID, r.CommissionID, r.Value, r.Comment, r.CoproID); err != nil {
			tx.Rollback()
			return err
		}
	}
	queries := []string{`UPDATE copro_forecast SET commission_id=t.commission_id,
	value=t.value,comment=t.comment,copro_id=t.copro_id 
	FROM temp_copro_forecast t WHERE t.id = copro_forecast.id`,
		`INSERT INTO copro_forecast (commission_id,value,comment,copro_id)
	SELECT commission_id,value,comment,copro_id from temp_copro_forecast 
		WHERE id NOT IN (SELECT id from copro_forecast)`,
		`DELETE from temp_copro_forecast`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	tx.Commit()
	return nil
}
