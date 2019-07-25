package models

import (
	"database/sql"
	"fmt"
)

// CityReportLine is used to decode a line of the query that fetches commitments
// and payments per policy and per year in a city
type CityReportLine struct {
	Policy     string `Json:"Policy"`
	Year       int64  `Json:"Year"`
	Commitment int64  `Json:"Commitment"`
	Payment    int64  `Json:"Payment"`
}

// CityReport embeddes an array of CityReportLine for json export
type CityReport struct {
	Lines []CityReportLine `json:"CityReport"`
}

// GetAll fetches commitments and payments per policy and year in a city
func (c *CityReport) GetAll(db *sql.DB, inseeCode, firstYear, lastYear int64) (err error) {
	qry := fmt.Sprintf(`WITH
	hcmt AS (SELECT cmt.year,SUM(cmt.value) as cmt
		FROM cumulated_commitment cmt
		JOIN housing h ON cmt.housing_id=h.id
		WHERE h.zip_code=$1
		GROUP BY 1),
  hpmt AS (SELECT p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN housing h ON c.housing_id=h.id
		WHERE h.zip_code=$1
		GROUP BY 1),
  hfig AS (SELECT hcmt.year,hcmt.cmt,hpmt.pmt FROM hcmt
  	FULL OUTER JOIN hpmt ON hcmt.year=hpmt.year),
	ccmt AS (SELECT cmt.year,SUM(cmt.value) as cmt 
		FROM cumulated_commitment cmt
		JOIN copro co ON cmt.copro_id=co.id
		WHERE co.zip_code=$1
		GROUP BY 1),
  cpmt AS (SELECT p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN copro co ON c.copro_id=co.id
		WHERE co.zip_code=$1
		GROUP BY 1),
  cfig AS (SELECT ccmt.year,ccmt.cmt,cpmt.pmt FROM ccmt
		FULL OUTER JOIN cpmt ON ccmt.year=cpmt.year),
	rcmt AS (SELECT cmt.year,SUM(cmt.value) as cmt 
		FROM cumulated_commitment cmt
		JOIN rp_cmt_city_join rp ON cmt.id=rp.commitment_id
		WHERE rp.city_code=$1
		GROUP BY 1),
  rpmt AS (SELECT p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN rp_cmt_city_join rp ON c.id=rp.commitment_id
		WHERE rp.city_code=$1
		GROUP BY 1),
  rfig AS (SELECT rcmt.year,rcmt.cmt,rpmt.pmt FROM rcmt
    FULL OUTER JOIN rpmt ON rpmt.year=rcmt.year)
	SELECT y.year,p.policy,COALESCE(q.cmt,0),COALESCE(q.pmt,0)
   FROM (SELECT generate_series(%d,%d) as year) y
   CROSS JOIN (SELECT * FROM (VALUES ('Housing'),('Copro'),('RenewProject')) AS t(policy)) p
	LEFT OUTER JOIN (
		SELECT 'Housing' as policy,year,COALESCE(cmt,0) AS cmt,COALESCE(pmt,0) AS pmt FROM hfig UNION ALL 
			SELECT 'Copro' as policy,year,COALESCE(cmt,0) AS cmt,COALESCE(pmt,0) AS pmt FROM cfig UNION ALL 
			SELECT 'RenewProject' as policy,year,COALESCE(cmt,0) AS cmt,COALESCE(pmt,0) AS pmt FROM rfig
	) q
		ON q.year=y.year AND q.policy=p.policy
		ORDER BY 1,2;`, firstYear, lastYear)
	rows, err := db.Query(qry, inseeCode)
	if err != nil {
		return err
	}
	var r CityReportLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Year, &r.Policy, &r.Commitment, &r.Payment); err != nil {
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
