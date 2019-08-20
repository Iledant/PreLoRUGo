package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// Prog model includes fields for a better readability in frontend and matching
// pre programmation fields
type Prog struct {
	ID             NullInt64  `json:"ID"`
	Year           int64      `json:"Year"`
	CommissionID   int64      `json:"CommissionID"`
	CommissionDate NullTime   `json:"CommissionDate"`
	CommissionName string     `json:"CommissionName"`
	Value          NullInt64  `json:"Value"`
	PreProgValue   NullInt64  `json:"PreProgValue"`
	Kind           int64      `json:"Kind"`
	KindID         NullInt64  `json:"KindID"`
	KindName       NullString `json:"KindName"`
	ProgComment    NullString `json:"ProgComment"`
	PreProgComment NullString `json:"PreProgComment"`
	ActionID       int64      `json:"ActionID"`
	ActionCode     int64      `json:"ActionCode"`
	ActionName     string     `json:"ActionName"`
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
 SELECT q.*,b.code,b.name,c.date,c.name FROM
 (SELECT p.id,$1 AS year,COALESCE(p.commission_id,pre.commission_id) AS commission_id,
	 p.value,pre.value AS preprog_value,NULL::int as kind_id,1 as kind,
	 NULL::varchar(150) as kind_name,p.comment as prog_comment,
	 pre.comment as pre_prog_comment,COALESCE(p.action_id,pre.action_id) as action_id
	 FROM p
 FULL OUTER JOIN (SELECT DISTINCT pp.ID,pp.commission_id,pp.value,pp.kind,
	 NULL::int as kind_id,NULL::varchar(150) as name,pp.comment,pp.action_id
 FROM pp WHERE pp.kind=1) pre
 ON p.commission_id=pre.commission_id AND p.kind=pre.kind AND p.action_id=pre.action_id
 WHERE p.kind ISNULL or p.kind=1
 UNION ALL
 SELECT pc.id,$1 AS year,COALESCE(pc.commission_id,pre.commission_id) AS commission_id,
	 pc.value,pre.value AS preprog_value,COALESCE(pc.kind_id,pre.kind_id) AS kind_id,
	 2 as kind,COALESCE(pc.name,pre.name) as kind_name,pc.comment as prog_comment,
	 pre.comment as pre_prog_comment,COALESCE(pc.action_id,pre.action_id) as action_id
	 FROM (SELECT p.id,p.year,p.commission_id,p.value,p.kind_id,p.kind,p.comment,
			c.name,p.action_id
		 FROM p JOIN copro c ON p.kind_id=c.id WHERE kind=2) pc
 FULL OUTER JOIN (SELECT DISTINCT pp.ID,pp.commission_id,pp.value,pp.kind,
	 pp.kind_id,c.name,pp.comment,pp.action_id
 FROM pp JOIN copro c ON pp.kind_id=c.id WHERE pp.kind=2) pre
 ON pc.commission_id=pre.commission_id AND pc.action_id=pre.action_id AND pc.kind_id=pre.kind_id
 WHERE pc.kind ISNULL or pc.kind=2
 UNION ALL
 SELECT pc.id,$1 AS year,COALESCE(pc.commission_id,pre.commission_id) AS commission_id,
	 pc.value,pre.value AS preprog_value,COALESCE(pc.kind_id,pre.kind_id) AS kind_id,
	 3 AS kind,COALESCE(pc.name,pre.name) as kind_name,
	 pc.comment as prog_comment,pre.comment as pre_prog_comment,
	 COALESCE(pc.action_id,pre.action_id) as action_id
	 FROM (SELECT p.id,p.year,p.commission_id,p.value,p.kind_id,p.kind,p.comment,
			c.name,p.action_id
		 FROM p JOIN renew_project c ON p.kind_id=c.id WHERE kind=3) pc
 FULL OUTER JOIN (SELECT DISTINCT pp.ID,pp.commission_id,pp.value,pp.kind,
	 pp.kind_id,c.name,pp.comment,pp.action_id
 FROM pp JOIN renew_project c ON pp.kind_id=c.id WHERE pp.kind=3) pre
 ON pc.commission_id=pre.commission_id AND pc.action_id=pre.action_id AND pc.kind_id=pre.kind_id
 WHERE pc.kind ISNULL or pc.kind=3) q
 JOIN budget_action b ON q.action_id=b.id
 JOIN commission c ON q.commission_id=c.id`, year)
	if err != nil {
		return err
	}
	var row Prog
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CommissionID, &row.Value,
			&row.PreProgValue, &row.KindID, &row.Kind, &row.KindName, &row.ProgComment,
			&row.PreProgComment, &row.ActionID, &row.ActionCode, &row.ActionName,
			&row.CommissionDate, &row.CommissionName); err != nil {
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
