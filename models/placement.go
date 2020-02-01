package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// Placement model
type Placement struct {
	ID              int64      `json:"ID"`
	IrisCode        string     `json:"IrisCode"`
	Count           NullInt64  `json:"Count"`
	ContractYear    NullInt64  `json:"ContractYear"`
	Comment         NullString `json:"Comment"`
	CreationDate    NullTime   `json:"CreationDate"`
	BeneficiaryName NullString `json:"BeneficiaryName"`
	BeneficiaryCode NullInt64  `json:"BeneficiaryCode"`
	CommitmentValue NullInt64  `json:"CommitmentValue"`
}

// Placements embeddes an array of Placement for json export and dedicated queries
type Placements struct {
	Lines []Placement `json:"Placement"`
}

// Update changes the comment of a placement
func (p *Placement) Update(db *sql.DB) error {
	_, err := db.Exec(`UPDATE placement SET comment=$1 WHERE id=$2`, p.Comment, p.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	if err = db.QueryRow(`SELECT p.iris_code,p.count,p.contract_year,
		MIN(c.creation_date),c.value,b.code,b.name
	FROM placement p
	LEFT OUTER JOIN commitment c ON p.iris_code=c.iris_code
	JOIN beneficiary b ON c.beneficiary_id=b.id
	WHERE p.id=$1
	GROUP BY 1,2,3,5,6,7`, p.ID).
		Scan(&p.IrisCode, &p.Count, &p.ContractYear, &p.CreationDate, &p.CommitmentValue,
			&p.BeneficiaryCode, &p.BeneficiaryName); err != nil {
		return fmt.Errorf("select %v", err)
	}
	return nil
}

// Get fetches all placements from database
func (p *Placements) Get(db *sql.DB) error {
	rows, err := db.Query(`SELECT p.id,p.iris_code,p.count,p.contract_year,
	p.comment,MIN(c.creation_date),c.value,b.code,b.name FROM placement p
		LEFT OUTER JOIN commitment c ON p.iris_code=c.iris_code
		JOIN beneficiary b ON c.beneficiary_id=b.id
		GROUP BY 1,2,3,4,5,7,8,9`)
	if err != nil {
		return err
	}
	var row Placement
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.IrisCode, &row.Count, &row.ContractYear,
			&row.Comment, &row.CreationDate, &row.CommitmentValue, &row.BeneficiaryCode,
			&row.BeneficiaryName); err != nil {
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
	rows, err := db.Query(`SELECT p.id,p.iris_code,p.count,p.contract_year,
	p.comment,c.creation_date FROM placement p
	JOIN commitment c ON p.iris_code=c.iris_code
	WHERE c.beneficiary_id=$1`, bID)
	if err != nil {
		return err
	}
	var row Placement
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.IrisCode, &row.Count,
			&row.ContractYear, &row.Comment, &row.CreationDate); err != nil {
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

// GetByBeneficiaryGroup fetches all placements linked to the beneficiaries that
// belong to a group
func (p *Placements) GetByBeneficiaryGroup(bID int64, db *sql.DB) error {
	rows, err := db.Query(`SELECT p.id,p.iris_code,c.value,b.code,b.name,p.count,p.contract_year,
	p.comment,MIN(c.creation_date) FROM placement p
	JOIN commitment c ON p.iris_code=c.iris_code
	JOIN beneficiary b ON c.beneficiary_id=b.id
	WHERE c.beneficiary_id IN 
		(SELECT beneficiary_id FROM beneficiary_belong WHERE group_id=$1)
	GROUP BY 1,2,3,4,5,6,7,8`, bID)
	if err != nil {
		return err
	}
	var row Placement
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.IrisCode, &row.CommitmentValue, &row.BeneficiaryCode, &row.BeneficiaryName, &row.Count,
			&row.ContractYear, &row.Comment, &row.CreationDate); err != nil {
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
		"contract_year"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range p.Lines {
		if _, err = stmt.Exec(r.IrisCode, r.Count, r.ContractYear); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	queries := []string{`INSERT INTO placement(iris_code,count,contract_year,
		commitment_id) 
	SELECT t.iris_code,t.count,t.contract_year,MIN(c.id) FROM temp_placement t
	LEFT OUTER JOIN commitment c ON t.iris_code=c.iris_code
	WHERE t.iris_code NOT IN (SELECT DISTINCT iris_code FROM placement)
  GROUP BY 1,2,3`,
		`UPDATE placement SET count=t.count,contract_year=t.contract_year,
		commitment_id=t.id FROM
		(SELECT t.iris_code,t.count,t.contract_year,MIN(c.id) id FROM temp_placement t
			LEFT OUTER JOIN commitment c ON t.iris_code=c.iris_code
			WHERE t.iris_code IN (SELECT DISTINCT iris_code FROM placement)
			GROUP BY 1,2,3) t
		WHERE placement.iris_code=t.iris_code`,
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
