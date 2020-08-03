package models

import (
	"database/sql"
	"fmt"
)

// AveragePayment model used to calculate the cumulated average percentage of
// payment of a month to help doing forecasts
type AveragePayment struct {
	Month       int64   `json:"Month"`
	PaymentRate float64 `json:"PaymentRate"`
}

// AveragePayments embeddes an array of AveragePayment for dedicated queries and
// json exports
type AveragePayments struct {
	Lines []AveragePayment `json:"AveragePayment"`
}

// GetAll fetches the cumulated payment rates of each month of the past years
func (a *AveragePayments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`WITH 
		q as (SELECT m,SUM(v) OVER (ORDER BY m) FROM
    	(SELECT EXTRACT(month FROM creation_date)::int m,SUM(value)::bigint v
			FROM payment
			WHERE creation_date<date(EXTRACT(year FROM current_date)||'-01-01')
			GROUP BY 1) q),
  	ma as (SELECT max(sum) FROM q)
	SELECT m,q.sum/ma.max FROM q,ma ORDER BY 1`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var r AveragePayment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Month, &r.PaymentRate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		a.Lines = append(a.Lines, r)
	}
	if len(a.Lines) == 0 {
		a.Lines = []AveragePayment{}
	}
	err = rows.Err()
	return err
}
