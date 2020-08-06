package models

import (
	"database/sql"
	"fmt"
)

// RPMultiAnnualReport model
type RPMultiAnnualReport struct {
	InseeCode        int64     `json:"InseeCode"`
	CityName         string    `json:"CityName"`
	RenewProjectName string    `json:"RenewProjectName"`
	Budget           NullInt64 `json:"Budget"`
	Commitment       NullInt64 `json:"Commitment"`
	Prog             NullInt64 `json:"Prog"`
	Y1               NullInt64 `json:"Y1"`
	Y2               NullInt64 `json:"Y2"`
	Y3               NullInt64 `json:"Y3"`
	Y4               NullInt64 `json:"Y4"`
	Y5               NullInt64 `json:"Y5"`
}

// RPMultiAnnualReports embeddes an array of RPMultiAnnualReport for json export
// and database fetching
type RPMultiAnnualReports struct {
	Lines []RPMultiAnnualReport `json:"RPMultiAnnualReport"`
}

// GetAll fetches all lines of the copro report from database
func (c *RPMultiAnnualReports) GetAll(db *sql.DB) error {
	qry := `
  WITH max_cmt_dat AS (SELECT MAX(creation_date) AS d FROM commitment)
  SELECT ci.insee_code,ci.name,rp.name,rp.budget,cmt.value,prg.value,
    q.y1,q.y2,q.y3,q.y4,q.y5 FROM
  (SELECT * FROM
		crosstab('SELECT rf.id,EXTRACT(year FROM c.date),rf.value
			FROM renew_project_forecast rf
      JOIN commission c ON rf.commission_id=c.id ORDER BY 1',
      'SELECT * FROM generate_series(EXTRACT(year FROM CURRENT_DATE)::bigint+1,
      EXTRACT(year FROM CURRENT_DATE)::bigint+5)')
    AS (id int, y1 BIGINT, y2 BIGINT, y3 BIGINT, y4 BIGINT, y5 BIGINT)) q
  JOIN renew_project_forecast rf ON q.id=rf.id
  JOIN renew_project rp ON rf.renew_project_id=rp.id
  JOIN (SELECT SUM(value)::bigint AS value,renew_project_id FROM commitment 
    WHERE renew_project_id NOTNULL GROUP BY 2) cmt ON cmt.renew_project_id=rp.id
	LEFT OUTER JOIN (SELECT SUM(value)::bigint AS value,kind_id AS renew_project_id
		FROM prog 
    JOIN commission rp ON prog.commission_id=rp.id
    WHERE kind=2 AND rp.date>(SELECT d FROM max_cmt_dat) GROUP BY 2) prg 
    ON prg.renew_project_id=rp.id
  JOIN city ci ON rp.city_code1=ci.insee_code;`
	rows, err := db.Query(qry)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var l RPMultiAnnualReport
	for rows.Next() {
		if err = rows.Scan(&l.InseeCode, &l.CityName, &l.RenewProjectName, &l.Budget,
			&l.Commitment, &l.Prog, &l.Y1, &l.Y2, &l.Y3, &l.Y4, &l.Y5); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		c.Lines = append(c.Lines, l)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(c.Lines) == 0 {
		c.Lines = []RPMultiAnnualReport{}
	}
	return nil
}
