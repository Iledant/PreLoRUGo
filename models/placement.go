package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// Placement model
type Placement struct {
	ID           int64      `json:"ID"`
	IrisCode     string     `json:"IrisCode"`
	Count        NullInt64  `json:"Count"`
	ContractYear NullInt64  `json:"ContractYear"`
	Comment      NullString `json:"Comment"`
}

// Placements embeddes an array of Placement for json export and dedicated queries
type Placements struct {
	Lines []Placement `json:"Placement"`
}

// Update changes the comment of a placement
func (p *Placement) Update(db *sql.DB) error {
	err := db.QueryRow(`UPDATE placement SET comment=$1 WHERE id=$2
		RETURNING iris_code,count,contract_year`,
		p.Comment, p.ID).Scan(&p.IrisCode, &p.Count, &p.ContractYear)
	if err != nil {
		return err
	}
	return nil
}

// Get fetches all placements from database
func (p *Placements) Get(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,iris_code,count,contract_year,comment
		FROM placement`)
	if err != nil {
		return err
	}
	var row Placement
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.IrisCode, &row.Count,
			&row.ContractYear, &row.Comment); err != nil {
			return err
		}
		p.Lines = append(p.Lines, row)
	}
	err = rows.Err()
	if len(p.Lines) == 0 {
		p.Lines = []Placement{}
	}
	return err
}

// GetByBeneficiary fetches all placements linked to a beneficiary
func (p *Placements) GetByBeneficiary(bID int64, db *sql.DB) error {
	rows, err := db.Query(`SELECT p.id,p.iris_code,p.count,p.contract_year,p.comment
	FROM placement p
	JOIN commitment c ON p.iris_code=c.iris_code
	WHERE c.beneficiary_id=$1`, bID)
	if err != nil {
		return err
	}
	var row Placement
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.IrisCode, &row.Count,
			&row.ContractYear, &row.Comment); err != nil {
			return err
		}
		p.Lines = append(p.Lines, row)
	}
	err = rows.Err()
	if len(p.Lines) == 0 {
		p.Lines = []Placement{}
	}
	return err
}

// Save update the database with a set of Placement
func (p *Placements) Save(db *sql.DB) error {
	for i, r := range p.Lines {
		if r.IrisCode == "" {
			return fmt.Errorf("Ligne %d, iris_code vide", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_placement", "iris_code", "count",
		"contract_year", "comment"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range p.Lines {
		if _, err = stmt.Exec(r.IrisCode, r.Count, r.ContractYear, r.Comment); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	queries := []string{`INSERT INTO placement(iris_code,count,contract_year,comment,
		commitment_id) 
	SELECT t.iris_code,t.count,t.contract_year,t.comment,c.id FROM temp_placement t
	LEFT OUTER JOIN commitment c ON t.iris_code=c.iris_code
	ON CONFLICT (iris_code) DO UPDATE SET count=excluded.count,
		contract_year=excluded.contract_year, comment=excluded.comment, 
		commitment_id=excluded.commitment_id`,
		`DELETE FROM temp_placement`,
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
