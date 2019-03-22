package models

import (
	"database/sql"
	"errors"
)

// RenewProjectForecast model
type RenewProjectForecast struct {
	ID             int64      `json:"ID"`
	CommissionID   int64      `json:"CommissionID"`
	Value          int64      `json:"Value"`
	Comment        NullString `json:"Comment"`
	RenewProjectID int64      `json:"RenewProjectID"`
}

// RenewProjectForecasts embeddes an array of RenewProjectForecast for json export
type RenewProjectForecasts struct {
	RenewProjectForecasts []RenewProjectForecast `json:"RenewProjectForecast"`
}

// RenewProjectForecastLine is used to decode a line of RenewProjectForecast batch
type RenewProjectForecastLine struct {
	ID             int64      `json:"ID"`
	CommissionID   int64      `json:"CommissionID"`
	Value          int64      `json:"Value"`
	Comment        NullString `json:"Comment"`
	RenewProjectID int64      `json:"RenewProjectID"`
}

// RenewProjectForecastBatch embeddes an array of RenewProjectForecastLine for json export
type RenewProjectForecastBatch struct {
	Lines []RenewProjectForecastLine `json:"RenewProjectForecast"`
}

// Validate checks if RenewProjectForecast's fields are correctly filled
func (r *RenewProjectForecast) Validate() error {
	if r.CommissionID == 0 || r.Value == 0 || r.RenewProjectID == 0 {
		return errors.New("Champ incorrect")
	}
	return nil
}

// Create insert a new RenewProjectForecast into database
func (r *RenewProjectForecast) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO renew_project_forecast 
	(commission_id,value,comment,renew_project_id)
 VALUES($1,$2,$3,$4) RETURNING id`, &r.CommissionID, &r.Value, &r.Comment,
		&r.RenewProjectID).Scan(&r.ID)
	return err
}

// Get fetches a RenewProjectForecast from database using ID field
func (r *RenewProjectForecast) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT commission_id, value, comment, renew_project_id 
	FROM renew_project_forecast WHERE ID=$1`, r.ID).Scan(&r.CommissionID, &r.Value,
		&r.Comment, &r.RenewProjectID)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a renew_project_forecast in database
func (r *RenewProjectForecast) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE renew_project_forecast SET commission_id=$1,value=$2,
	comment=$3,renew_project_id=$4 WHERE id=$5`,
		r.CommissionID, r.Value, r.Comment, r.RenewProjectID, r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Prévision RU introuvable")
	}
	return err
}

// GetAll fetches all RenewProjectForecasts from database
func (r *RenewProjectForecasts) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,commission_id,value,comment,renew_project_id 
	FROM renew_project_forecast`)
	if err != nil {
		return err
	}
	var row RenewProjectForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommissionID, &row.Value, &row.Comment,
			&row.RenewProjectID); err != nil {
			return err
		}
		r.RenewProjectForecasts = append(r.RenewProjectForecasts, row)
	}
	err = rows.Err()
	return err
}

// Delete removes renew_project_forecast whose ID is given from database
func (r *RenewProjectForecast) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM renew_project_forecast WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Prévision RU introuvable")
	}
	return nil
}

// Save insert a batch of RenewProjectForecastLine into database
func (r *RenewProjectForecastBatch) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_renew_project_forecast 
	(id, commission_id,value,comment,renew_project_id) VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range r.Lines {
		if r.CommissionID == 0 || r.Value == 0 || r.RenewProjectID == 0 {
			tx.Rollback()
			return errors.New("Champs incorrects")
		}
		if _, err = stmt.Exec(r.ID, r.CommissionID, r.Value, r.Comment, r.RenewProjectID); err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`UPDATE renew_project_forecast SET commission_id=t.commission_id,
	value=t.value,comment=t.comment,renew_project_id=t.renew_project_id 
	FROM temp_renew_project_forecast t WHERE t.id = renew_project_forecast.id`)
	if err != nil {
		tx.Rollback()
		return errors.New("UPDATE " + err.Error())
	}
	_, err = tx.Exec(`INSERT INTO renew_project_forecast (commission_id,value,comment,renew_project_id)
	SELECT commission_id,value,comment,renew_project_id from temp_renew_project_forecast 
	  WHERE id NOT IN (SELECT id from renew_project_forecast)`)

	if err != nil {
		tx.Rollback()
		return errors.New("INSERT : " + err.Error())
	}
	tx.Commit()
	return nil
}
