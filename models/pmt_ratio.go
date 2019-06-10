package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// PmtRatio is used to calculate one line of payment transformation ratio of
// commitments for all budget actions
type PmtRatio struct {
	Index      int     `json:"Index"`
	SectorID   int     `json:"SectorID"`
	SectorName string  `json:"SectorName"`
	Ratio      float64 `json:"Ratio"`
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
	rows, err := db.Query(`WITH cmt AS (SELECT SUM(c.value) AS value, a.sector_id 
		FROM commitment c, budget_action a
		WHERE extract(year FROM c.creation_date)=$1 AND c.value>0 AND c.action_id=a.id
		GROUP BY 2),
  pmt AS (SELECT SUM(p.value) AS value, 
    extract(year from p.creation_date)::integer-$1 AS index, a.sector_id
    FROM payment p, commitment c, budget_action a WHERE p.commitment_id=c.id 
      AND c.value >0 AND c.action_id=a.id 
      AND extract(year FROM c.creation_date)=$1 GROUP BY 2,3)
	SELECT pmt.index, pmt.sector_id, s.name, pmt.value/cmt.value AS ratio 
	FROM cmt, pmt, budget_sector s 
	WHERE cmt.sector_id=pmt.sector_id AND pmt.sector_id=s.id`, year)
	if err != nil {
		return fmt.Errorf("get request %v", err)
	}
	var r PmtRatio
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Index, &r.SectorID, &r.SectorName, &r.Ratio); err != nil {
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
	stmt, err := tx.Prepare(pq.CopyIn("ratio", "year", "index", "sector_id",
		"ratio"))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("statement creation %v", err)
	}
	defer stmt.Close()
	for _, r := range p.Ratios {
		if _, err = stmt.Exec(p.Year, r.Index, r.SectorID, r.Ratio); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement execution %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
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
