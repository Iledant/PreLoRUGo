package models

import (
	"database/sql"
	"fmt"
)

// DptReportLine is used to decode one line of the department query to fetch
// commitments and payments linked to all kinds of projects
type DptReportLine struct {
	Code       int64  `json:"Code"`
	Name       string `json:"Name"`
	Year       int64  `json:"Year"`
	Commitment int64  `json:"Commitment"`
	Payment    int64  `json:"Payment"`
}

// DptReport embeddes an array par DptReportLine for json export
type DptReport struct {
	Lines []DptReportLine `json:"DptReport"`
}

// GetAll fetches the commitment and payment per department from database
func (d *DptReport) GetAll(db *sql.DB, firstYear, lastYear int64) (err error) {
	qry := fmt.Sprintf(`WITH
	hcmt AS (SELECT d.code as dpt_code,cmt.year,SUM(cmt.value) as cmt
		FROM cumulated_commitment cmt
		JOIN housing h ON cmt.housing_id=h.id
		JOIN department d ON d.code=h.zip_code/1000
		GROUP BY 1,2),
  hpmt AS (SELECT d.code AS dpt_code,p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN housing h ON c.housing_id=h.id
		JOIN department d ON d.code=h.zip_code/1000
		GROUP BY 1,2),
  hfig AS (SELECT hcmt.dpt_code,hcmt.year,hcmt.cmt,hpmt.pmt FROM hcmt
  	FULL OUTER JOIN hpmt ON hcmt.dpt_code=hpmt.dpt_code AND hcmt.year=hpmt.year),
	ccmt AS (SELECT d.code as dpt_code,cmt.year,SUM(cmt.value) as cmt 
		FROM cumulated_commitment cmt
		JOIN copro co ON cmt.copro_id=co.id
		JOIN department d ON d.code=co.zip_code/1000
		GROUP BY 1,2),
  cpmt AS (SELECT d.code AS dpt_code,p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN copro co ON c.copro_id=co.id
		JOIN department d ON d.code=co.zip_code/1000
		GROUP BY 1,2),
  cfig AS (SELECT ccmt.dpt_code,ccmt.year,ccmt.cmt,cpmt.pmt FROM ccmt
		FULL OUTER JOIN cpmt ON ccmt.dpt_code=cpmt.dpt_code AND ccmt.year=cpmt.year),
	rcmt AS (SELECT d.code as dpt_code,cmt.year,SUM(cmt.value) as cmt 
		FROM cumulated_commitment cmt
		JOIN rp_cmt_city_join rp ON cmt.id=rp.commitment_id
		JOIN department d ON d.code=rp.city_code/1000
		GROUP BY 1,2),
  rpmt AS (SELECT d.code AS dpt_code,p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN rp_cmt_city_join rp ON c.id=rp.commitment_id
		JOIN department d ON d.code=rp.city_code/1000
		GROUP BY 1,2),
  rfig AS (SELECT rcmt.dpt_code,rcmt.year,rcmt.cmt,rpmt.pmt FROM rcmt
    FULL OUTER JOIN rpmt ON rpmt.dpt_code=rcmt.dpt_code AND rpmt.year=rcmt.year),
  tfig AS (SELECT q.dpt_code,q.year,SUM(q.cmt) as cmt,SUM(q.pmt) as pmt FROM 
    (SELECT * FROM cfig UNION ALL SELECT * FROM hfig UNION ALL SELECT * FROM rfig)q
    GROUP BY 1,2)
	SELECT dpt.code,dpt.name,y.year,COALESCE(tfig.cmt,0)::bigint,
			COALESCE(tfig.pmt,0)::bigint
		FROM department dpt
		CROSS JOIN (SELECT generate_series(%d,%d) as year) y
		LEFT OUTER JOIN tfig ON tfig.dpt_code=dpt.code AND tfig.year=y.year
		ORDER BY 1,2;`, firstYear, lastYear)
	rows, err := db.Query(qry)
	if err != nil {
		return err
	}
	var r DptReportLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Code, &r.Name, &r.Year, &r.Commitment, &r.Payment); err != nil {
			return err
		}
		d.Lines = append(d.Lines, r)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	if len(d.Lines) == 0 {
		d.Lines = []DptReportLine{}
	}
	return err
}
