package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// Prog model includes fields for a better readability in frontend and matching
// pre programmation fields
type Prog struct {
	CommissionID    int64      `json:"CommissionID"`
	CommissionDate  NullTime   `json:"CommissionDate"`
	CommissionName  string     `json:"CommissionName"`
	ActionID        int64      `json:"ActionID"`
	ActionCode      int64      `json:"ActionCode"`
	ActionName      string     `json:"ActionName"`
	Kind            int64      `json:"Kind"`
	KindID          NullInt64  `json:"KindID"`
	KindName        NullString `json:"KindName"`
	ForecastValue   NullInt64  `json:"ForecastValue"`
	ForecastComment NullString `json:"ForecastComment"`
	PreProgValue    NullInt64  `json:"PreProgValue"`
	PreProgComment  NullString `json:"PreProgComment"`
	ID              NullInt64  `json:"ID"`
	Value           NullInt64  `json:"Value"`
	Comment         NullString `json:"Comment"`
}

// Progs embeddes an array of Prog for json export
type Progs struct {
	Progs []Prog `json:"Prog"`
}

// ProgLine is used to decode one line of Prog batch
type ProgLine struct {
	CommissionID int64      `json:"CommissionID"`
	Value        int64      `json:"Value"`
	Kind         int64      `json:"Kind"`
	KindID       NullInt64  `json:"KindID"`
	Comment      NullString `json:"Comment"`
	ActionID     int64      `json:"ActionID"`
}

// ProgBatch embeddes an array of ProgLine to import a batch
type ProgBatch struct {
	Lines []ProgLine `json:"Prog"`
}

// ProgYears embeddes an array of int64 for json export fetching the available
// years with programmation data in the database
type ProgYears struct {
	Years []int64 `json:"ProgYear"`
}

// GetAll fetches all Prog of a given year from the database
func (p *Progs) GetAll(year int64, db *sql.DB) error {
	rows, err := db.Query(`WITH p AS (SELECT * FROM prog WHERE year=$1),
	pp AS (SELECT * FROM pre_prog WHERE year=$1)
 SELECT q.commission_id,c.date,c.name,q.action_id,b.code,b.name,q.kind,q.kind_id,
	 q.kind_name,q.forecast_value,q.forecast_comment,q.pre_prog_value,
	 q.pre_prog_comment,q.id,q.value,q.prog_comment
 FROM
 (
   SELECT p.id,COALESCE(p.commission_id,pre.commission_id) AS commission_id,
	 	COALESCE(p.action_id,pre.action_id) as action_id,
   	p.value,pre.pre_prog_value,NULL::int as kind_id,1 as kind,
	 	NULL::varchar(150) as kind_name,p.comment as prog_comment,
	 	pre.pre_prog_comment,pre.forecast_value,pre.forecast_comment
	 FROM p
		FULL OUTER JOIN
		(SELECT DISTINCT COALESCE(pp.commission_id,hf.commission_id) AS commission_id,
			pp.value AS pre_prog_value,pp.kind,pp.comment AS pre_prog_comment,
			COALESCE(pp.action_id,hf.action_id) AS action_id,hf.value AS forecast_value,
			hf.comment AS forecast_comment
 		FROM pp 
		FULL OUTER JOIN
		(SELECT DISTINCT hf.commission_id, hf.action_id, hf.value,hf.comment 
 		FROM housing_forecast hf 
 		JOIN commission co ON hf.commission_id=co.id
 		WHERE extract(year FROM co.date) = $1) hf
 		ON pp.commission_id=hf.commission_id AND pp.action_id=hf.action_id
 		WHERE pp.kind=1) pre
 	ON p.commission_id=pre.commission_id AND p.kind=pre.kind AND p.action_id=pre.action_id
 	WHERE p.kind ISNULL or p.kind=1

	UNION ALL

	SELECT pc.id,COALESCE(pc.commission_id,pre.commission_id) AS commission_id,
		COALESCE(pc.action_id,pre.action_id) as action_id,
		pc.value,pre.pre_prog_value,COALESCE(pc.kind_id,pre.kind_id) AS kind_id,
		2 as kind,COALESCE(pc.name,pre.name) as kind_name,pc.comment as prog_comment,
		pre.pre_prog_comment,pre.forecast_value, pre.forecast_comment
	FROM
		(SELECT p.id,p.year,p.commission_id,p.value,p.kind_id,p.kind,p.comment,
			c.name,p.action_id
		 FROM p
		 JOIN copro c ON p.kind_id=c.id
		 WHERE kind=2
		) pc
	 FULL OUTER JOIN
	 	(SELECT DISTINCT COALESCE(pp.commission_id,cf.commission_id) AS commission_id,
			pp.value AS pre_prog_value,2::int as kind,COALESCE(pp.kind_id,cf.kind_id) AS kind_id,
			COALESCE(c.name,cf.name) as name,pp.comment as pre_prog_comment,
			 COALESCE(pp.action_id,cf.action_id) AS action_id,cf.value as forecast_value,
			 cf.comment AS forecast_comment
 		FROM pp
 		JOIN copro c ON pp.kind_id=c.id 
		FULL OUTER JOIN
			(SELECT DISTINCT cf.commission_id, cf.action_id, cf.value,cf.comment,
 				cf.copro_id as kind_id,c.name
 			FROM copro_forecast cf 
 			JOIN commission co ON cf.commission_id=co.id
 			JOIN copro c ON cf.copro_id=c.id
 			WHERE extract(year FROM co.date) = $1) cf
		ON pp.commission_id=cf.commission_id AND pp.action_id=cf.action_id AND pp.kind_id=cf.kind_id
 		WHERE pp.kind=2
 		) pre
 		ON pc.commission_id=pre.commission_id AND pc.action_id=pre.action_id AND pc.kind_id=pre.kind_id
 		WHERE pc.kind ISNULL or pc.kind=2
 
	UNION ALL
	 
	SELECT pc.id,COALESCE(pc.commission_id,pre.commission_id) AS commission_id,
		COALESCE(pc.action_id,pre.action_id) as action_id,pc.value,
		pre.value AS pre_prog_value,COALESCE(pc.kind_id,pre.kind_id) AS kind_id,
		3 AS kind,COALESCE(pc.name,pre.name) as kind_name,pc.comment as prog_comment,
		pre.pre_prog_comment,pre.forecast_value,pre.forecast_comment
	FROM
		(SELECT p.id,p.year,p.commission_id,p.value,p.kind_id,p.kind,p.comment,
			c.name,p.action_id
		 FROM p
		 JOIN renew_project c
		 ON p.kind_id=c.id
		 WHERE kind=3) pc
		 FULL OUTER JOIN
		 	(SELECT DISTINCT pp.ID,
				COALESCE(pp.commission_id,rf.commission_id) AS commission_id,pp.value,
				COALESCE(pp.kind_id,rf.kind_id) AS kind_id,COALESCE(c.name,rf.name) AS name,
				pp.comment AS pre_prog_comment,COALESCE(pp.action_id,rf.action_id) AS action_id,
				rf.value AS forecast_value,rf.comment AS forecast_comment
			 FROM pp
			 JOIN renew_project c ON pp.kind_id=c.id
			 FULL OUTER JOIN
			 	(SELECT DISTINCT rf.commission_id, rf.action_id, rf.value,rf.comment,
 					rf.renew_project_id as kind_id,rp.name
 				FROM renew_project_forecast rf 
 				JOIN commission co ON rf.commission_id=co.id
 				JOIN renew_project rp ON rf.renew_project_id=rp.id
 				WHERE extract(year FROM co.date) = $1) rf
 			ON pp.commission_id=rf.commission_id AND pp.action_id=rf.action_id AND pp.kind_id=rf.kind_id
			WHERE pp.kind=3
		) pre
 		ON pc.commission_id=pre.commission_id AND pc.action_id=pre.action_id AND pc.kind_id=pre.kind_id
  	WHERE pc.kind ISNULL or pc.kind=3
	) q
	JOIN budget_action b ON q.action_id=b.id
	JOIN commission c ON q.commission_id=c.id
	ORDER BY date,kind,code,kind_id`, year)
	if err != nil {
		return err
	}
	var row Prog
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.CommissionID, &row.CommissionDate, &row.CommissionName,
			&row.ActionID, &row.ActionCode, &row.ActionName, &row.Kind, &row.KindID,
			&row.KindName, &row.ForecastValue, &row.ForecastComment, &row.PreProgValue,
			&row.PreProgComment, &row.ID, &row.Value, &row.Comment); err != nil {
			return err
		}
		p.Progs = append(p.Progs, row)
	}
	err = rows.Err()
	if len(p.Progs) == 0 {
		p.Progs = []Prog{}
	}
	return err
}

// Save insert a batch of ProgLine into the database. It checks if the
// batch includes only one year , otherwise throw an error. It replaces all the
//  datas of the given year, deleting the programming data of that year
// in the database
func (p *ProgBatch) Save(year int64, db *sql.DB) error {
	for i, l := range p.Lines {
		if l.CommissionID == 0 {
			return fmt.Errorf("ligne %d, CommissionID nul", i+1)
		}
		if l.Value == 0 {
			return fmt.Errorf("ligne %d, Value nul", i+1)
		}
		if l.ActionID == 0 {
			return fmt.Errorf("ligne %d, ActionID nul", i+1)
		}
		if l.Kind != KindCopro && l.Kind != KindRenewProject && l.Kind != KindHousing {
			return fmt.Errorf("linge %d, Kind de mauvais type", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("début de transaction %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_prog", "commission_id",
		"year", "value", "kind", "kind_id", "comment", "action_id"))
	if err != nil {
		return fmt.Errorf("insert statement %v", err)
	}
	defer stmt.Close()
	for _, l := range p.Lines {
		if _, err = stmt.Exec(l.CommissionID, year, l.Value, l.Kind, l.KindID,
			l.Comment, l.ActionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement execution %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	if _, err = tx.Exec(`DELETE FROM prog WHERE year=$1`, year); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete query %v", err)
	}
	queries := []string{`INSERT INTO prog (commission_id,year,value,kind,kind_id,
		comment,action_id) SELECT commission_id,year,value,kind,kind_id,
		comment,action_id FROM temp_prog `,
		`DELETE from temp_prog`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	tx.Commit()
	return nil
}

// GetAll fetches all programmation years in the database
func (p *ProgYears) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT DISTINCT year from prog`)
	if err != nil {
		return err
	}
	var y int64
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&y); err != nil {
			return err
		}
		p.Years = append(p.Years, y)
	}
	err = rows.Err()
	if len(p.Years) == 0 {
		p.Years = []int64{}
	}
	return err
}
