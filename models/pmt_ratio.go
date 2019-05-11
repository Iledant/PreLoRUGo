package models

import (
	"database/sql"
	"fmt"
)

// PmtRatio is used to calculate one line of payment transformation ratio of
// commitments for all budget actions
type PmtRatio struct {
	Index int     `json:"Index"`
	Ratio float64 `json:"Ratio"`
}

// PmtRatios embeddes an array of PmtRatio for json export
type PmtRatios struct {
	PmtRatios []PmtRatio `json:"PmtRatio"`
}

// PmtRatioBatch is used to decode the payload of a post request fo save ratios
// of a given year
type PmtRatioBatch struct {
	Year   int        `json:"Year"`
	Ratios []PmtRatio `json:"Ratios"`
}

// PmtRatiosYears is used to fetch years with ratios payments from database
type PmtRatiosYears struct {
	Years []int `json:"PmtRatiosYear"`
}

// Get fetches the payment transformation ratios of commitments for the given
// year
func (p *PmtRatios) Get(db *sql.DB, year int) error {
	rows, err := db.Query(`WITH cmt AS (SELECT SUM(value) AS value FROM commitment 
  WHERE extract(year FROM creation_date)=$1 AND value>0),
  pmt AS (SELECT SUM(p.value) AS value, 
    extract(year from p.creation_date)::integer-$1 AS index
    FROM payment p, commitment c WHERE p.commitment_id=c.id AND c.value >0 
      AND extract(year FROM c.creation_date)=$1 GROUP BY 2)
SELECT pmt.index, pmt.value/cmt.value AS ratio FROM cmt, pmt`, year)
	if err != nil {
		return fmt.Errorf("get request %v", err)
	}
	var r PmtRatio
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Index, &r.Ratio); err != nil {
			return err
		}
		p.PmtRatios = append(p.PmtRatios, r)
	}
	err = rows.Err()
	if len(p.PmtRatios) == 0 {
		p.PmtRatios = []PmtRatio{}
	}
	return err
}

// Save updates or inserts the payment ratios of a given year
func (p *PmtRatioBatch) Save(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM ratio WHERE year=$1`, p.Year); err != nil {
		tx.Rollback()
		return fmt.Errorf("DELETE %v", err)
	}
	stmt, err := tx.Prepare(`INSERT INTO ratio (year, index, ratio) VALUES ($1,$2,$3)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("statement creation %v", err)
	}
	defer stmt.Close()
	for _, r := range p.Ratios {
		if _, err = stmt.Exec(p.Year, r.Index, r.Ratio); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement execution %v", err)
		}
	}
	tx.Commit()
	return nil
}

// Get fetches all years from ratio table in the database
func (p *PmtRatiosYears) Get(db *sql.DB) error {
	rows, err := db.Query(`SELECT DISTINCT year FROM ratio`)
	if err != nil {
		return fmt.Errorf("get request %v", err)
	}
	var y int
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&y); err != nil {
			return err
		}
		p.Years = append(p.Years, y)
	}
	err = rows.Err()
	if len(p.Years) == 0 {
		p.Years = []int{}
	}
	return err
}
