package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// KindHousing is the PreProg and Prog database field content for housing
const KindHousing = "Housing"

// KindCopro is the PreProg and Prog database field content for copro
const KindCopro = "Copro"

// KindRenewProject is the PreProg and Prog database field content for renew project
const KindRenewProject = "RenewProject"

const preProgHousingQry = `SELECT DISTINCT pp.ID,pp.year,pp.commission_id,c.date,
c.name,pp.value,pp.kind,NULL::int,NULL::varchar(150),pp.comment,pp.action_id,b.code,b.name 
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
WHERE pp.kind='Housing' AND pp.year=$1`
const preProgCoproQry = `SELECT DISTINCT pp.ID,pp.year,pp.commission_id,c.date,
c.name, pp.value,pp.kind,pp.kind_id,copro.name,pp.comment,pp.action_id,b.code,
b.name 
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
LEFT JOIN copro ON copro.id=pp.kind_id
WHERE pp.kind='Copro' AND pp.year=$1`
const preProgRPQry = `SELECT DISTINCT pp.ID,pp.year,pp.commission_id,c.date,
c.name,pp.value,pp.kind,pp.kind_id,rp.name,pp.comment,pp.action_id,b.code,b.name 
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
LEFT JOIN renew_project rp ON rp.id=pp.kind_id
WHERE pp.kind='RenewProject' AND pp.year=$1`

// PreProg model includes fields for a better readability in frontend
type PreProg struct {
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

// PreProgs embeddes an array of PreProg for json export
type PreProgs struct {
	PreProgs []PreProg `json:"PreProg"`
}

// PreProgLine is used to decode one line of PreProg batch
type PreProgLine struct {
	CommissionID int64      `json:"CommissionID"`
	Year         int64      `json:"Year"`
	Value        int64      `json:"Value"`
	KindID       NullInt64  `json:"KindID"`
	Comment      NullString `json:"Comment"`
	ActionID     int64      `json:"ActionID"`
}

// PreProgBatch embeddes an array of PreProgLine to import a batch
type PreProgBatch struct {
	Lines []PreProgLine `json:"PreProg"`
}

// GetAll fetches all PreProg of a given year from the database
func (p *PreProgs) GetAll(year int64, db *sql.DB) error {
	qry := fmt.Sprintf("%s UNION ALL %s UNION ALL %s", preProgHousingQry,
		preProgCoproQry, preProgRPQry)
	rows, err := db.Query(qry, year)
	if err != nil {
		return err
	}
	var row PreProg
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Kind, &row.KindID, &row.KindName,
			&row.Comment, &row.ActionID, &row.ActionCode, &row.ActionName); err != nil {
			return err
		}
		p.PreProgs = append(p.PreProgs, row)
	}
	err = rows.Err()
	if len(p.PreProgs) == 0 {
		p.PreProgs = []PreProg{}
	}
	return err
}

// GetAllOfKind fetches all PreProg of a given year and kind from the database
func (p *PreProgs) GetAllOfKind(year int64, kind string, db *sql.DB) error {
	var qry string
	switch kind {
	case KindHousing:
		qry = preProgHousingQry
	case KindCopro:
		qry = preProgCoproQry
	case KindRenewProject:
		qry = preProgRPQry
	}
	rows, err := db.Query(qry, year)
	if err != nil {
		return err
	}
	var row PreProg
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Kind, &row.KindID, &row.KindName,
			&row.Comment, &row.ActionID, &row.ActionCode, &row.ActionName); err != nil {
			return err
		}
		p.PreProgs = append(p.PreProgs, row)
	}
	err = rows.Err()
	if len(p.PreProgs) == 0 {
		p.PreProgs = []PreProg{}
	}
	return err
}

// Save insert a batch of PreProgLine into the database. It checks if the
// batch includes only one year, otherwise throw an error. It replaces all
// the datas of the given year and kinds, deleting PreProgData of that year and
// kind in the database
func (p *PreProgBatch) Save(kind string, year int64, db *sql.DB) error {
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
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("début de transaction %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_pre_prog", "commission_id",
		"year", "value", "kind", "kind_id", "comment", "action_id"))
	if err != nil {
		return fmt.Errorf("insert statement %v", err)
	}
	defer stmt.Close()
	for _, l := range p.Lines {
		if _, err = stmt.Exec(l.CommissionID, year, l.Value, kind, l.KindID,
			l.Comment, l.ActionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement execution %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	if _, err = tx.Exec(`DELETE FROM pre_prog WHERE year=$1 AND kind=$2`,
		year, kind); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete query %v", err)
	}
	queries := []string{`INSERT INTO pre_prog (commission_id,year,value,kind,kind_id,
		comment,action_id) SELECT commission_id,year,value,kind,kind_id,
		comment,action_id FROM temp_pre_prog `,
		`DELETE from temp_pre_prog`,
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
