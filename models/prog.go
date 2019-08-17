package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// Prog model includes fields for a better readability in frontend
type Prog struct {
	ID             int64      `json:"ID"`
	Year           int64      `json:"Year"`
	CommissionID   int64      `json:"CommissionID"`
	CommissionDate NullTime   `json:"CommissionDate"`
	CommissionName string     `json:"CommissionName"`
	Value          int64      `json:"Value"`
	Kind           string     `json:"Kind"`
	KindID         NullInt64  `json:"KindID"`
	KindName       NullString `json:"KindName"`
	Comment        NullString `json:"Comment"`
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
	Kind         string     `json:"Kind"`
	KindID       NullInt64  `json:"KindID"`
	Comment      NullString `json:"Comment"`
	ActionID     int64      `json:"ActionID"`
}

// ProgBatch embeddes an array of ProgLine to import a batch
type ProgBatch struct {
	Lines []ProgLine `json:"Prog"`
}

// GetAll fetches all Prog of a given year from the database
func (p *Progs) GetAll(year int64, db *sql.DB) error {
	rows, err := db.Query(`SELECT DISTINCT p.ID,p.year,p.commission_id,c.date,
	c.name,p.value,p.kind,NULL::int,NULL::varchar(150),p.comment,p.action_id,b.code,b.name 
	FROM prog p
	JOIN commission c ON c.id=p.commission_id
	JOIN budget_action b ON b.id=p.action_id
	WHERE p.kind='Housing' AND p.year=$1
	UNION ALL
	SELECT DISTINCT p.ID,p.year,p.commission_id,c.date,
	c.name, p.value,p.kind,p.kind_id,copro.name,p.comment,p.action_id,b.code,
	b.name 
	FROM prog p
	JOIN commission c ON c.id=p.commission_id
	JOIN budget_action b ON b.id=p.action_id
	LEFT JOIN copro ON copro.id=p.kind_id
	WHERE p.kind='Copro' AND p.year=$1
	UNION ALL
	SELECT DISTINCT p.ID,p.year,p.commission_id,c.date,
	c.name,p.value,p.kind,p.kind_id,rp.name,p.comment,p.action_id,b.code,b.name 
	FROM prog p
	JOIN commission c ON c.id=p.commission_id
	JOIN budget_action b ON b.id=p.action_id
	LEFT JOIN renew_project rp ON rp.id=p.kind_id
	WHERE p.kind='RenewProject' AND p.year=$1`, year)
	if err != nil {
		return err
	}
	var row Prog
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Kind, &row.KindID, &row.KindName,
			&row.Comment, &row.ActionID, &row.ActionCode, &row.ActionName); err != nil {
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
