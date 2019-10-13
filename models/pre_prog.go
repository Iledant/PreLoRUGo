package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// KindHousing is the PreProg and Prog database field content for housing
const KindHousing = 1

// KindCopro is the PreProg and Prog database field content for copro
const KindCopro = 2

// KindRenewProject is the PreProg and Prog database field content for renew project
const KindRenewProject = 3

const preProgHousingQry = `SELECT DISTINCT pp.ID,pp.year,pp.commission_id,c.date,
c.name,pp.value,pp.kind,NULL::int,NULL::varchar(150),pp.project,pp.comment,pp.action_id,b.code,b.name 
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
WHERE pp.kind=1 AND pp.year=$1`
const preProgCoproQry = `SELECT DISTINCT pp.ID,pp.year,pp.commission_id,c.date,
c.name, pp.value,pp.kind,pp.kind_id,copro.name,pp.project,pp.comment,pp.action_id,b.code,
b.name 
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
LEFT JOIN copro ON copro.id=pp.kind_id
WHERE pp.kind=2 AND pp.year=$1`
const preProgRPQry = `SELECT DISTINCT pp.ID,pp.year,pp.commission_id,c.date,
c.name,pp.value,pp.kind,pp.kind_id,rp.name,pp.project,pp.comment,pp.action_id,b.code,b.name 
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
LEFT JOIN renew_project rp ON rp.id=pp.kind_id
WHERE pp.kind=3 AND pp.year=$1`

const fcPreProgHousingQry = `SELECT COALESCE(pp.commission_id,hf.commission_id),
COALESCE(pp.date,hf.date), COALESCE(pp.name,hf.name),COALESCE(pp.action_id,hf.action_id),
COALESCE(pp.code,hf.code),COALESCE(pp.action_name,hf.action_name),
NULL::int,NULL::varchar(150),hf.id,hf.value,NULL::varchar(150),hf.comment,pp.id,pp.value,
NULL::varchar(150),pp.comment
FROM
(SELECT DISTINCT pp.ID,pp.commission_id,c.date,
c.name,pp.value,pp.comment,pp.action_id,b.code,b.name as action_name
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
WHERE pp.kind=1 AND pp.year=$1) pp
FULL OUTER JOIN
(SELECT hf.id,hf.commission_id,c.date,c.name,hf.value,hf.comment,hf.action_id,
b.code,b.name  as action_name
FROM housing_forecast hf
JOIN commission c ON c.id=hf.commission_id
JOIN budget_action b ON b.id=hf.action_id
WHERE EXTRACT(year FROM c.date)=$1) hf
ON pp.commission_id=hf.commission_id AND pp.action_id=hf.action_id`

const fcPreProgCoproQry = `SELECT DISTINCT COALESCE(pp.commission_id,cf.commission_id),
COALESCE(pp.date,cf.date),COALESCE(pp.name,cf.name),COALESCE(pp.action_id,cf.action_id),
COALESCE(pp.code,cf.code),COALESCE(pp.action_name,cf.action_name),COALESCE(pp.kind_id,cf.copro_id),
COALESCE(pp.copro_name,cf.copro_name),cf.id,cf.value,cf.project,cf.comment,pp.id,
pp.value,pp.project,pp.comment
FROM
(SELECT DISTINCT pp.ID,pp.commission_id,c.date,
c.name,pp.value,pp.project,pp.comment,pp.action_id,pp.kind_id,co.name as copro_name,b.code,b.name as action_name
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
JOIN copro co ON pp.kind_id=co.id
WHERE pp.kind=2 AND pp.year=$1) pp
FULL OUTER JOIN
(SELECT cf.id,cf.commission_id,c.date,c.name,cf.value,cf.project,cf.comment,cf.action_id,
cf.copro_id,co.name as copro_name,b.code,b.name  as action_name
FROM copro_forecast cf
JOIN commission c ON c.id=cf.commission_id
JOIN budget_action b ON b.id=cf.action_id
JOIN copro co ON cf.copro_id=co.id
WHERE EXTRACT(year FROM c.date)=$1) cf
ON pp.commission_id=cf.commission_id AND pp.action_id=cf.action_id AND pp.kind_id=cf.copro_id`

const fcPreProgRPQry = `SELECT DISTINCT COALESCE(pp.commission_id,rf.commission_id),
COALESCE(pp.date,rf.date),COALESCE(pp.name,rf.name),COALESCE(pp.action_id,rf.action_id),
COALESCE(pp.code,rf.code),COALESCE(pp.action_name,rf.action_name),
COALESCE(pp.kind_id,rf.renew_project_id),COALESCE(pp.rp_name,rf.rp_name),rf.id,
rf.value,rf.project,rf.comment,pp.id,pp.value,pp.project,pp.comment
FROM
(SELECT DISTINCT pp.ID,pp.commission_id,c.date,
c.name,pp.value,pp.project,pp.comment,pp.action_id,pp.kind_id,city.name || ' - ' || rp.name as rp_name,b.code,b.name as action_name
FROM pre_prog pp
JOIN commission c ON c.id=pp.commission_id
JOIN budget_action b ON b.id=pp.action_id
JOIN renew_project rp ON pp.kind_id=rp.id
JOIN city ON rp.city_code1=city.insee_code
WHERE pp.kind=3 AND pp.year=$1) pp
FULL OUTER JOIN
(SELECT rf.id,rf.commission_id,c.date,c.name,rf.value,rf.project,rf.comment,rf.action_id,
rf.renew_project_id,rp.name as rp_name,b.code,b.name  as action_name
FROM renew_project_forecast rf
JOIN commission c ON c.id=rf.commission_id
JOIN budget_action b ON b.id=rf.action_id
JOIN renew_project rp ON rf.renew_project_id=rp.id
WHERE EXTRACT(year FROM c.date)=$1) rf
ON pp.commission_id=rf.commission_id AND pp.action_id=rf.action_id AND pp.kind_id=rf.renew_project_id`

const preProgQry = preProgHousingQry + " UNION ALL " + preProgCoproQry +
	" UNION ALL " + preProgRPQry

// PreProg model includes fields for a better readability in frontend
type PreProg struct {
	ID             int64      `json:"ID"`
	Year           int64      `json:"Year"`
	CommissionID   int64      `json:"CommissionID"`
	CommissionDate NullTime   `json:"CommissionDate"`
	CommissionName string     `json:"CommissionName"`
	Value          int64      `json:"Value"`
	Kind           int64      `json:"Kind"`
	KindID         NullInt64  `json:"KindID"`
	KindName       NullString `json:"KindName"`
	KindProject    NullString `json:"KindProject"`
	Comment        NullString `json:"Comment"`
	ActionID       int64      `json:"ActionID"`
	ActionCode     int64      `json:"ActionCode"`
	ActionName     string     `json:"ActionName"`
}

// PreProgs embeddes an array of PreProg for json export
type PreProgs struct {
	PreProgs []PreProg `json:"PreProg"`
}

// FcPreProg model includes fields for a better readability in frontend
type FcPreProg struct {
	CommissionID    int64      `json:"CommissionID"`
	CommissionDate  NullTime   `json:"CommissionDate"`
	CommissionName  string     `json:"CommissionName"`
	ActionID        int64      `json:"ActionID"`
	ActionCode      int64      `json:"ActionCode"`
	ActionName      string     `json:"ActionName"`
	KindID          NullInt64  `json:"KindID"`
	KindName        NullString `json:"KindName"`
	ForecastID      NullInt64  `json:"ForecastID"`
	ForecastValue   NullInt64  `json:"ForecastValue"`
	ForecastComment NullString `json:"ForecastComment"`
	ForecastProject NullString `json:"ForecastProject"`
	PreProgID       NullInt64  `json:"PreProgID"`
	PreProgValue    NullInt64  `json:"PreProgValue"`
	PreProgComment  NullString `json:"PreProgComment"`
	PreProgProject  NullString `json:"PreProgProject"`
}

// FcPreProgs embeddes an array of FcPreProg for json export
type FcPreProgs struct {
	FcPreProgs []FcPreProg `json:"FcPreProg"`
}

// PreProgLine is used to decode one line of PreProg batch
type PreProgLine struct {
	CommissionID int64      `json:"CommissionID"`
	Year         int64      `json:"Year"`
	Value        int64      `json:"Value"`
	KindID       NullInt64  `json:"KindID"`
	Comment      NullString `json:"Comment"`
	Project      NullString `json:"Project"`
	ActionID     int64      `json:"ActionID"`
}

// PreProgBatch embeddes an array of PreProgLine to import a batch
type PreProgBatch struct {
	Lines []PreProgLine `json:"PreProg"`
}

// GetAll fetches all PreProg of a given year from the database
func (p *PreProgs) GetAll(year int64, db *sql.DB) error {
	rows, err := db.Query(preProgQry, year)
	if err != nil {
		return err
	}
	var row PreProg
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CommissionID, &row.CommissionDate,
			&row.CommissionName, &row.Value, &row.Kind, &row.KindID, &row.KindName,
			&row.KindProject, &row.Comment, &row.ActionID, &row.ActionCode, &row.ActionName); err != nil {
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
func (p *PreProgs) GetAllOfKind(year int64, kind int64, db *sql.DB) error {
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
			&row.KindProject, &row.Comment, &row.ActionID, &row.ActionCode, &row.ActionName); err != nil {
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

// GetAllOfKind fetches all FcPreProg of a given year and kind from the database
func (p *FcPreProgs) GetAllOfKind(year int64, kind int64, db *sql.DB) error {
	var qry string
	switch kind {
	case KindHousing:
		qry = fcPreProgHousingQry
	case KindCopro:
		qry = fcPreProgCoproQry
	case KindRenewProject:
		qry = fcPreProgRPQry
	}
	rows, err := db.Query(qry, year)
	if err != nil {
		return err
	}
	var row FcPreProg
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.CommissionID, &row.CommissionDate, &row.CommissionName,
			&row.ActionID, &row.ActionCode, &row.ActionName, &row.KindID, &row.KindName,
			&row.ForecastID, &row.ForecastValue, &row.ForecastProject, &row.ForecastComment,
			&row.PreProgID, &row.PreProgValue, &row.PreProgProject, &row.PreProgComment); err != nil {
			return err
		}
		p.FcPreProgs = append(p.FcPreProgs, row)
	}
	err = rows.Err()
	if len(p.FcPreProgs) == 0 {
		p.FcPreProgs = []FcPreProg{}
	}
	return err
}

// Save insert a batch of PreProgLine into the database. It checks if the
// batch includes only one year, otherwise throw an error. It replaces all
// the datas of the given year and kinds, deleting PreProgData of that year and
// kind in the database
func (p *PreProgBatch) Save(kind int64, year int64, db *sql.DB) error {
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
		"year", "value", "kind", "kind_id", "project", "comment", "action_id"))
	if err != nil {
		return fmt.Errorf("insert statement %v", err)
	}
	defer stmt.Close()
	for _, l := range p.Lines {
		if _, err = stmt.Exec(l.CommissionID, year, l.Value, kind, l.KindID,
			l.Project, l.Comment, l.ActionID); err != nil {
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
		project,comment,action_id) SELECT commission_id,year,value,kind,kind_id,
		project,comment,action_id FROM temp_pre_prog `,
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
