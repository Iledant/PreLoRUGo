package models

import (
	"database/sql"
	"fmt"
)

// CoproReport model
type CoproReport struct {
	InseeCode  int64     `json:"InseeCode"`
	CityName   string    `json:"CityName"`
	CoproName  string    `json:"CoproName"`
	Budget     NullInt64 `json:"Budget"`
	Commitment NullInt64 `json:"Commitment"`
	Prog       NullInt64 `json:"Prog"`
	Y1         NullInt64 `json:"Y1"`
	Y2         NullInt64 `json:"Y2"`
	Y3         NullInt64 `json:"Y3"`
	Y4         NullInt64 `json:"Y4"`
	Y5         NullInt64 `json:"Y5"`
}

// CoproReports embeddes an array of CoproReport for json export and database fetching
type CoproReports struct {
	Lines []CoproReport `json:"CoproReport"`
}

// GetAll fetches all lines of the copro report from database
func (c *CoproReports) GetAll(db *sql.DB) error {
	qry := `
  WITH max_cmt_dat AS (SELECT max(creation_date) AS d FROM commitment)
  SELECT ci.insee_code,ci.name,co.name,co.budget,cmt.value,prg.value,
    q.y1,q.y2,q.y3,q.y4,q.y5 FROM
  (SELECT * FROM
    crosstab('SELECT cf.id,EXTRACT(year FROM c.date),cf.value from copro_forecast cf
      JOIN commission c ON cf.commission_id=c.id ORDER BY 1',
      'select * FROM generate_series(extract(year FROM CURRENT_DATE)::bigint+1,
      extract(year FROM CURRENT_DATE)::bigint+5)')
    AS (id int, y1 BIGINT, y2 BIGINT, y3 BIGINT, y4 BIGINT, y5 BIGINT)) q
  JOIN copro_forecast cf ON q.id=cf.id
  JOIN copro co ON cf.copro_id=co.id
  JOIN (SELECT sum(value)::bigint AS value,copro_id FROM commitment 
    WHERE copro_id NOTNULL GROUP BY 2) cmt ON cmt.copro_id=co.id
  LEFT OUTER JOIN (SELECT sum(value)::bigint AS value,kind_id AS copro_id FROM prog 
    JOIN commission co ON prog.commission_id=co.id
    WHERE kind=2 AND co.date>(SELECT d FROM max_cmt_dat) GROUP BY 2) prg 
    ON prg.copro_id=co.id
  JOIN city ci ON co.zip_code=ci.insee_code;`
	rows, err := db.Query(qry)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var l CoproReport
	for rows.Next() {
		if err = rows.Scan(&l.InseeCode, &l.CityName, &l.CoproName, &l.Budget, &l.Commitment,
			&l.Prog, &l.Y1, &l.Y2, &l.Y3, &l.Y4, &l.Y5); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		c.Lines = append(c.Lines, l)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(c.Lines) == 0 {
		c.Lines = []CoproReport{}
	}
	return nil
}
