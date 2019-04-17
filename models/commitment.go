package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	HousingID        NullInt64  `json:"HousingID"`
	CoproID          NullInt64  `json:"CoproID"`
	RenewProjectID   NullInt64  `json:"RenewProjectID"`
}

// PaginatedCommitment is used for paginated query providing beneficiary name
type PaginatedCommitment struct {
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
	BeneficiaryName  string     `json:"BeneficiaryName"`
	IrisCode         NullString `json:"IrisCode"`
	HousingID        NullInt64  `json:"HousingID"`
	CoproID          NullInt64  `json:"CoproID"`
	RenewProjectID   NullInt64  `json:"RenewProjectID"`
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
	CreationDate     int64      `json:"CreationDate"`
	ModificationDate int64      `json:"ModificationDate"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	BeneficiaryCode  int64      `json:"BeneficiaryCode"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	IrisCode         NullString `json:"IrisCode"`
	Sector           string     `json:"Sector"`
	ActionCode       NullInt64  `json:"ActionCode"`
	ActionName       NullString `json:"ActionName"`
}

// CommitmentBatch embeddes an array of CommitmentLine for json export
type CommitmentBatch struct {
	Lines []CommitmentLine `json:"Commitment"`
}

// CommitmentQuery embeddes the request to fetch some commitments from database
// according to the given pattern
type CommitmentQuery struct {
	Page   int64  `json:"Page"`
	Year   int64  `json:"Year"`
	Search string `json:"Search"`
}

// PaginatedCommitments embeddes the query results of a CommitmentQuery
type PaginatedCommitments struct {
	Commitments []PaginatedCommitment `json:"Commitment"`
	Page        int64                 `json:"Page"`
	ItemsCount  int64                 `json:"ItemsCount"`
}

// Get fetches the results of a paginated commitment query
func (p *PaginatedCommitments) Get(db *sql.DB, c *CommitmentQuery) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM commitment c 
		JOIN beneficiary b on c.beneficiary_id=b.id WHERE year >= $1 AND
		(c.name ILIKE $2 OR c.code ILIKE $2  OR b.name ILIKE $2)`, c.Year, "%"+c.Search+"%").
		Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(c.Page, count)

	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,c.creation_date,
	c.modification_date,c.name,c.value,c.beneficiary_id, b.name, c.iris_code FROM commitment c
	JOIN beneficiary b ON c.beneficiary_id = b.id
	WHERE year >= $1 AND (c.name ILIKE $2 OR c.code ILIKE $2 OR b.name ILIKE $2)
	ORDER BY 2,6,7,3,4,5 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $3`,
		c.Year, "%"+c.Search+"%", offset)
	if err != nil {
		return err
	}
	var row PaginatedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.Name, &row.Value,
			&row.BeneficiaryID, &row.BeneficiaryName, &row.IrisCode); err != nil {
			return err
		}
		p.Commitments = append(p.Commitments, row)
	}
	err = rows.Err()
	if len(p.Commitments) == 0 {
		p.Commitments = []PaginatedCommitment{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
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
	if len(c.Commitments) == 0 {
		c.Commitments = []Commitment{}
	}
	return err
}

// Save insert a batch of CommitmentLine into database
func (c *CommitmentBatch) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_commitment (year,code,number,line,
		creation_date,modification_date,name,value,beneficiary_code,beneficiary_name,
		iris_code,sector,action_code,action_name)
	VALUES ($1,$2,$3,$4,make_date($5,$6,$7),make_date($8,$9,$10),$11,$12,$13,$14,
	$15,$16,$17,$18)`)
	if err != nil {
		return errors.New("Statement creation " + err.Error())
	}
	defer stmt.Close()
	for _, r := range c.Lines {
		if r.Year < 2009 || r.Number == 0 || r.Line == 0 || r.CreationDate < 20090101 ||
			r.ModificationDate < 20090101 || r.Name == "" || r.BeneficiaryCode == 0 ||
			r.BeneficiaryName == "" || r.Sector == "" {
			tx.Rollback()
			return fmt.Errorf("Champs incorrects dans %+v", r)
		}
		if _, err = stmt.Exec(r.Year, r.Code, r.Number, r.Line, r.CreationDate/10000,
			r.CreationDate/100%100, r.CreationDate%100, r.ModificationDate/10000,
			r.ModificationDate/100%100, r.ModificationDate%100, strings.TrimSpace(r.Name),
			r.Value, r.BeneficiaryCode, strings.TrimSpace(r.BeneficiaryName), r.IrisCode,
			strings.TrimSpace(r.Sector), r.ActionCode, r.ActionName.TrimSpace()); err != nil {
			tx.Rollback()
			return errors.New("Statement execution " + err.Error())
		}
	}
	_, err = tx.Exec(`INSERT INTO beneficiary (code,name) SELECT DISTINCT beneficiary_code,beneficiary_name 
		FROM temp_commitment WHERE beneficiary_code not in (SELECT code from beneficiary)`)
	if err != nil {
		tx.Rollback()
		return errors.New("Beneficiary insertion " + err.Error())
	}
	_, err = tx.Exec(`INSERT INTO budget_sector (name) SELECT DISTINCT sector
	FROM temp_commitment WHERE sector not in (SELECT name from budget_sector)`)
	if err != nil {
		tx.Rollback()
		return errors.New("Budget sector insertion " + err.Error())
	}
	_, err = tx.Exec(`INSERT INTO budget_action (code,name,sector_id) 
		SELECT DISTINCT ic.action_code,ic.action_name, s.id
		FROM temp_commitment ic
		LEFT JOIN budget_sector s ON ic.sector = s.name
		WHERE action_code not in (SELECT code from budget_action)`)
	if err != nil {
		tx.Rollback()
		return errors.New("Budget action insertion " + err.Error())
	}
	_, err = tx.Exec(`INSERT INTO commitment (year,code,number,line,creation_date,modification_date,
		name,value,beneficiary_id,iris_code,action_id)
  	(SELECT ic.year,ic.code,ic.number,ic.line,ic.creation_date,ic.modification_date,
			ic.name,ic.value,b.id,ic.iris_code,a.id
  	FROM temp_commitment ic
		JOIN beneficiary b on ic.beneficiary_code=b.code
		LEFT JOIN budget_action a on ic.action_code = a.code
  	WHERE (ic.year,ic.code,ic.number,ic.line,ic.creation_date,ic.modification_date,ic.name, ic.value) 
    NOT IN (select year,code,number,line,creation_date,modification_date,name,value FROM commitment));`)
	if err != nil {
		tx.Rollback()
		return errors.New("Commitment insertion " + err.Error())
	}
	tx.Commit()
	return nil
}
