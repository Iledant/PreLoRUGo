package models

import (
	"database/sql"
	"fmt"
)

// CityReportLine is used to decode a line of the query that fetches commitments
// and payments per policy and per year in a city
type CityReportLine struct {
	Kind       int64 `Json:"Kind"`
	Year       int64 `Json:"Year"`
	Commitment int64 `Json:"Commitment"`
	Payment    int64 `Json:"Payment"`
}

// CityReport embeddes an array of CityReportLine for json export
type CityReport struct {
	Lines []CityReportLine `json:"CityReport"`
}

// GetAll fetches commitments and payments per policy and year in a city
func (c *CityReport) GetAll(db *sql.DB, inseeCode, firstYear, lastYear int64) (err error) {
	qry := fmt.Sprintf(`WITH
	housingCmt AS (SELECT cmt.year,SUM(cmt.value) cmt
		FROM cumulated_commitment cmt
		JOIN housing h ON cmt.housing_id=h.id
		WHERE h.zip_code=$1
		GROUP BY 1),
  housingPmt AS (SELECT p.year,SUM(p.value) pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN housing h ON c.housing_id=h.id
		WHERE h.zip_code=$1
		GROUP BY 1),
	housingFig AS (SELECT housingCmt.year,housingCmt.cmt,housingPmt.pmt
		FROM housingCmt
  	FULL OUTER JOIN housingPmt ON housingCmt.year=housingPmt.year),
	coproCmt AS (SELECT cmt.year,SUM(cmt.value) cmt 
		FROM cumulated_commitment cmt
		JOIN copro co ON cmt.copro_id=co.id
		WHERE co.zip_code=$1
		GROUP BY 1),
  coproPmt AS (SELECT p.year,SUM(p.value) pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN copro co ON c.copro_id=co.id
		WHERE co.zip_code=$1
		GROUP BY 1),
	coproFig AS (SELECT coproCmt.year,coproCmt.cmt,coproPmt.pmt
		FROM coproCmt
		FULL OUTER JOIN coproPmt ON coproCmt.year=coproPmt.year),
	renewProjectCmt AS (SELECT cmt.year,SUM(cmt.value) cmt 
		FROM cumulated_commitment cmt
		JOIN rp_cmt_city_join rp ON cmt.id=rp.commitment_id
		WHERE rp.city_code=$1
		GROUP BY 1),
  renewProjectPmt AS (SELECT p.year,SUM(p.value) pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN rp_cmt_city_join rp ON c.id=rp.commitment_id
		WHERE rp.city_code=$1
		GROUP BY 1),
	renewProjectFig AS (SELECT renewProjectCmt.year,renewProjectCmt.cmt,renewProjectPmt.pmt 
		FROM renewProjectCmt
    FULL OUTER JOIN renewProjectPmt ON renewProjectPmt.year=renewProjectCmt.year)
	SELECT y,k,COALESCE(q.cmt,0),COALESCE(q.pmt,0)
		FROM generate_series(%d,%d) y
		CROSS JOIN generate_series(1,3) k
		LEFT OUTER JOIN (
			SELECT 1 kind,year,COALESCE(cmt,0) cmt,COALESCE(pmt,0) pmt FROM housingFig
			UNION ALL 
			SELECT 2 kind,year,COALESCE(cmt,0) cmt,COALESCE(pmt,0) pmt FROM coproFig
			UNION ALL 
			SELECT 3 kind,year,COALESCE(cmt,0) cmt,COALESCE(pmt,0) pmt FROM renewProjectFig
		) q
			ON q.year=y AND q.kind=k
	ORDER BY 1,2;`, firstYear, lastYear)
	rows, err := db.Query(qry, inseeCode)
	if err != nil {
		return err
	}
	var r CityReportLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Year, &r.Kind, &r.Commitment, &r.Payment); err != nil {
			return err
		}
		c.Lines = append(c.Lines, r)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	if len(c.Lines) == 0 {
		c.Lines = []CityReportLine{}
	}
	return err
}
