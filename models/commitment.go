package models

import (
	"database/sql"
	"errors"
	"time"
)

// Commitment model
type Commitment struct {
	ID               int64      `json:"ID"`
	Year             int64      `json:"Year"`
	Code             string     `json:"Code"`
	Number           int64      `json:"Number"`
	Line             int64      `json:"Line"`
	CreationDate     time.Time  `json:"CreationDate"`
	ModificationDate time.Time  `json:"ModificationDate"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	BeneficiaryID    int64      `json:"BeneficiaryID"`
	IrisCode         NullString `json:"IrisCode"`
}

// Commitments embeddes an array of Commitment for json export
type Commitments struct {
	Commitments []Commitment `json:"Commitment"`
}

// CommitmentLine is used to decode a line of Commitment batch
type CommitmentLine struct {
	Year             int64      `json:"Year"`
	Code             string     `json:"Code"`
	Number           int64      `json:"Number"`
	Line             int64      `json:"Line"`
	CreationDate     int64      `json:"Creation"`
	ModificationDate int64      `json:"Modification"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	BeneficiaryCode  int64      `json:"BeneficiaryCode"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	IrisCode         NullString `json:"IrisCode"`
}

// CommitmentBatch embeddes an array of CommitmentLine for json export
type CommitmentBatch struct {
	Lines []CommitmentLine `json:"Commitment"`
}

// GetAll fetches all Commitments from database
func (c *Commitments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,year,code,number,line,creation_date,
	modification_date,name,value,beneficiary_id,iris_code FROM commitment`)
	if err != nil {
		return err
	}
	var row Commitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.Name, &row.Value,
			&row.BeneficiaryID, &row.IrisCode); err != nil {
			return err
		}
		c.Commitments = append(c.Commitments, row)
	}
	err = rows.Err()
	return err
}

// Save insert a batch of CommitmentLine into database
func (c *CommitmentBatch) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_commitment (year,code,number,line,
		creation_date,modification_date,name,value,beneficiary_code,beneficiary_name,iris_code) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range c.Lines {
		if r.Year < 2009 || r.Number == 0 || r.Line == 0 || r.CreationDate < 20090101 ||
			r.ModificationDate < 20090101 || r.Name == "" || r.BeneficiaryCode == 0 ||
			r.BeneficiaryName == "" {
			tx.Rollback()
			return errors.New("Champs incorrects")
		}
		if _, err = stmt.Exec(r.Year, r.Code, r.Number, r.Line, r.CreationDate,
			r.ModificationDate, r.Name, r.Value, r.BeneficiaryCode, r.BeneficiaryName,
			r.IrisCode); err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO beneficiary (code,name) SELECT DISTINCT beneficiary_code,beneficiary_name 
		FROM temp_commitment WHERE beneficiary_code not in (SELECT code from beneficiary)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO commitment (year,code,number,line,creation_date,modification_date,
		name,value,beneficiary_id,iris_code)
  	(SELECT ic.year,ic.code,ic.number,ic.line,make_date(ic.creation_date/10000,(ic.creation_date/100)%100,ic.creation_date%100),
			make_date(ic.modification_date/10000,(ic.modification_date/100)%100,ic.modification_date%100),
			ic.name,ic.value,b.id,ic.iris_code
  	FROM temp_commitment ic
  	JOIN beneficiary b on ic.beneficiary_code=b.code
  	WHERE (ic.year,ic.code,ic.number,ic.line,make_date(ic.creation_date/10000,(ic.creation_date/100)%100,ic.creation_date%100),
    make_date(ic.modification_date/10000,(ic.modification_date/100)%100,ic.modification_date%100),ic.name, ic.value) 
    NOT IN (select year,code,number,line,creation_date,modification_date,name,value FROM commitment));`)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
