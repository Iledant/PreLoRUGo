package models

import (
	"database/sql"
	"fmt"
	"time"
)

// AvgPmtTime model
type AvgPmtTime struct {
	Month             time.Time   `json:"Month"`
	AverageTime       float64     `json:"AverageTime"`
	StandardDeviation NullFloat64 `json:"StandardDeviation"`
}

// AvgPmtTimes embeddes an array of AvgPmtTime for json export and dedicated queries
type AvgPmtTimes struct {
	Lines []AvgPmtTime `json:"AveragePaymentTime"`
}

// GetAll fetches the average payments times of the past 12 monthes
func (a *AvgPmtTimes) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT m.d,AVG(p.creation_date-p.receipt_date),
	stddev_samp(p.creation_date-p.receipt_date) 
	FROM payment p,
	(SELECT CURRENT_DATE- i*make_interval(0,1) as d FROM generate_series(11,0,-1) i) m
	WHERE p.creation_date<=m.d AND p.creation_date>=m.d-make_interval(1)
	GROUP BY 1 ORDER BY 1;`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var line AvgPmtTime
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&line.Month, &line.AverageTime, &line.StandardDeviation); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		a.Lines = append(a.Lines, line)
	}
	err = rows.Err()
	if len(a.Lines) == 0 {
		a.Lines = []AvgPmtTime{}
	}
	return err
}
