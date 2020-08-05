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
	housingCmt AS (SELECT d.code as dpt_code,cmt.year,SUM(cmt.value) as cmt
		FROM cumulated_commitment cmt
		JOIN housing h ON cmt.housing_id=h.id
		JOIN department d ON d.code=h.zip_code/1000
		GROUP BY 1,2),
  housingPmt AS (SELECT d.code AS dpt_code,p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN housing h ON c.housing_id=h.id
		JOIN department d ON d.code=h.zip_code/1000
		GROUP BY 1,2),
  housingFig AS (SELECT COALESCE(housingCmt.dpt_code,housingPmt.dpt_code) AS dpt_code,
			COALESCE(housingCmt.year,housingPmt.year) AS year,housingCmt.cmt,housingPmt.pmt
		FROM housingCmt
		FULL OUTER JOIN housingPmt ON housingCmt.dpt_code=housingPmt.dpt_code 
			AND housingCmt.year=housingPmt.year),
	coproCmt AS (SELECT d.code as dpt_code,cmt.year,SUM(cmt.value) as cmt 
		FROM cumulated_commitment cmt
		JOIN copro co ON cmt.copro_id=co.id
		JOIN department d ON d.code=co.zip_code/1000
		GROUP BY 1,2),
  coproPmt AS (SELECT d.code AS dpt_code,p.year,SUM(p.value) as pmt FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN copro co ON c.copro_id=co.id
		JOIN department d ON d.code=co.zip_code/1000
		GROUP BY 1,2),
  coproFig AS (SELECT COALESCE(coproCmt.dpt_code,coproPmt.dpt_code) AS dpt_code,
			COALESCE(coproCmt.year,coproPmt.year) AS year,coproCmt.cmt,coproPmt.pmt
			FROM coproCmt
		FULL OUTER JOIN coproPmt ON coproCmt.dpt_code=coproPmt.dpt_code 
			AND coproCmt.year=coproPmt.year),
	renewProjectCmt AS (SELECT d.code as dpt_code,cmt.year,SUM(cmt.value) as cmt 
		FROM cumulated_commitment cmt
		JOIN rp_cmt_city_join rp ON cmt.id=rp.commitment_id
		JOIN department d ON d.code=rp.city_code/1000
		GROUP BY 1,2),
	renewProjectPmt AS (SELECT d.code AS dpt_code,p.year,SUM(p.value) as pmt 
		FROM payment p
		JOIN cumulated_commitment c ON p.commitment_id=c.id
		JOIN rp_cmt_city_join rp ON c.id=rp.commitment_id
		JOIN department d ON d.code=rp.city_code/1000
		GROUP BY 1,2),
  renewProjectFig AS (SELECT COALESCE(rc.dpt_code,rp.dpt_code) dpt_code,
		COALESCE(rc.year,rp.year) AS year,rc.cmt,rp.pmt
		FROM renewProjectCmt rc
		FULL OUTER JOIN renewProjectPmt rp ON rp.dpt_code=rc.dpt_code
			AND rp.year=rc.year),
  totalFig AS (SELECT q.dpt_code,q.year,SUM(q.cmt) as cmt,SUM(q.pmt) as pmt FROM 
		(SELECT * FROM coproFig UNION ALL SELECT * FROM housingFig
			UNION ALL SELECT * FROM renewProjectFig) q
    GROUP BY 1,2)
	SELECT dpt.code,dpt.name,y.year,COALESCE(totalFig.cmt,0)::bigint,
			COALESCE(totalFig.pmt,0)::bigint
		FROM department dpt
		CROSS JOIN (SELECT generate_series(%d,%d) as year) y
		LEFT OUTER JOIN totalFig ON totalFig.dpt_code=dpt.code AND totalFig.year=y.year
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
