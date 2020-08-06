package models

import (
	"database/sql"
	"fmt"
	"time"
)

// SoldCommitment model
type SoldCommitment struct {
	Year             int64      `json:"Year"`
	Code             string     `json:"Code"`
	Number           int64      `json:"Number"`
	Line             int64      `json:"Line"`
	CreationDate     time.Time  `json:"CreationDate"`
	ModificationDate time.Time  `json:"ModificationDate"`
	CaducityDate     NullTime   `json:"CaducityDate"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	SoldOut          bool       `json:"SoldOut"`
	Beneficiary      string     `json:"Beneficiary"`
	Sector           string     `json:"Sector"`
	ActionName       string     `json:"ActionName"`
	ActionCode       string     `json:"ActionCode"`
	IrisCode         NullString `json:"IrisCode"`
}

// SoldCommitments embeddes an array of SoldCommitment for json export and
// dedicated queries
type SoldCommitments struct {
	Lines []SoldCommitment `json:"SoldCommitment"`
}

// GetOld fetches all commitments older than 9 years wether they are sold out or
// not in order to deal with the annual exercice
func (s *SoldCommitments) GetOld(db *sql.DB) error {
	rows, err := db.Query(`SELECT c.year,c.code,c.number,c.line,c.creation_date,
	c.modification_date,c.caducity_date,c.name,c.value,c.sold_out,b.name,s.name,
	a.name,a.code,c.iris_code
	FROM commitment c
	JOIN beneficiary b on c.beneficiary_id=b.id
	JOIN budget_action a ON c.action_id=a.id
	JOIN budget_sector s ON a.sector_id=s.id
	WHERE EXTRACT(year FROM creation_date) <= EXTRACT(year FROM current_date)-9
	ORDER BY creation_date, value DESC;`)
	if err != nil {
		return fmt.Errorf("select %d", err)
	}
	var l SoldCommitment
	for rows.Next() {
		if err = rows.Scan(&l.Year, &l.Code, &l.Number, &l.Line, &l.CreationDate,
			&l.ModificationDate, &l.CaducityDate, &l.Name, &l.Value, &l.SoldOut,
			&l.Beneficiary, &l.Sector, &l.ActionName, &l.ActionCode, &l.IrisCode); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		s.Lines = append(s.Lines, l)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("scan err %v", err)
	}
	if len(s.Lines) == 0 {
		s.Lines = []SoldCommitment{}
	}
	return nil
}

// GetUnpaid fetches all commitments older than 3 years with no payments and that
// are not sold out
func (s *SoldCommitments) GetUnpaid(db *sql.DB) error {
	rows, err := db.Query(`SELECT q.year,q.code,q.number,q.line,q.creation_date,
	q.modification_date,q.caducity_date,q.name,q.value,q.sold_out,b.name,s.name,
	a.name,a.code,q.iris_code
	FROM 
	 (SELECT c.year,c.code,c.number,c.line,c.creation_date,c.modification_date,
    c.caducity_date,c.name,c.value,c.sold_out,c.iris_code,c.beneficiary_id,
    c.action_id,p.value AS payment FROM commitment c 
		LEFT OUTER JOIN payment p ON p.commitment_id=c.id
		WHERE EXTRACT(year FROM c.creation_date)<=EXTRACT(year FROM current_date)-3) q
	JOIN beneficiary b on q.beneficiary_id=b.id
	JOIN budget_action a ON q.action_id=a.id
	JOIN budget_sector s ON a.sector_id=s.id
	WHERE q.payment IS NULL AND q.sold_out=FALSE
  ORDER BY q.creation_date,q.value DESC;`)
	if err != nil {
		return fmt.Errorf("select %d", err)
	}
	var l SoldCommitment
	for rows.Next() {
		if err = rows.Scan(&l.Year, &l.Code, &l.Number, &l.Line, &l.CreationDate,
			&l.ModificationDate, &l.CaducityDate, &l.Name, &l.Value, &l.SoldOut,
			&l.Beneficiary, &l.Sector, &l.ActionName, &l.ActionCode, &l.IrisCode); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		s.Lines = append(s.Lines, l)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("scan err %v", err)
	}
	if len(s.Lines) == 0 {
		s.Lines = []SoldCommitment{}
	}
	return nil
}
