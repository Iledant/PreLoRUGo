package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
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
	CaducityDate     NullTime   `json:"CaducityDate"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	SoldOut          bool       `json:"SoldOut"`
	BeneficiaryID    int64      `json:"BeneficiaryID"`
	ActionID         int64      `json:"ActionID"`
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
	CaducityDate     NullTime   `json:"CaducityDate"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	SoldOut          bool       `json:"SoldOut"`
	BeneficiaryID    int64      `json:"BeneficiaryID"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	ActionName       string     `json:"ActionName"`
	Sector           string     `json:"Sector"`
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
	CaducityDate     int64      `json:"CaducityDate"`
	Name             string     `json:"Name"`
	Value            int64      `json:"Value"`
	SoldOut          string     `json:"SoldOut"`
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

// PaginatedQuery embeddes the request to fetch some commitments from database
// according to the given pattern
type PaginatedQuery struct {
	Page   int64  `json:"Page"`
	Year   int64  `json:"Year"`
	Search string `json:"Search"`
}

// PaginatedCommitments embeddes the query results of a PaginatedQuery
type PaginatedCommitments struct {
	Commitments []PaginatedCommitment `json:"Commitment"`
	Page        int64                 `json:"Page"`
	ItemsCount  int64                 `json:"ItemsCount"`
}

// ExportQuery embeddes the request to fetch some commitments from database
// according to the given pattern for export purpose
type ExportQuery struct {
	Year   int64  `json:"Year"`
	Search string `json:"Search"`
}

// ExportedCommitment is dedicated to commitments exports with explicit fields
type ExportedCommitment struct {
	ID               int64      `json:"ID"`
	Year             int64      `json:"Year"`
	Code             string     `json:"Code"`
	Number           int64      `json:"Number"`
	Line             int64      `json:"Line"`
	CreationDate     time.Time  `json:"CreationDate"`
	ModificationDate time.Time  `json:"ModificationDate"`
	CaducityDate     NullTime   `json:"CaducityDate"`
	Name             string     `json:"Name"`
	Value            float64    `json:"Value"`
	SoldOut          bool       `json:"SoldOut"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	Sector           string     `json:"Sector"`
	ActionName       string     `json:"ActionName"`
	IrisCode         NullString `json:"IrisCode"`
	HousingName      NullString `json:"HousingName"`
	CoproName        NullString `json:"CoproName"`
	RenewProjectName NullString `json:"RenewProjectName"`
}

// ExportedCommitments embeddes an array of ExportedCommitment for json export
type ExportedCommitments struct {
	ExportedCommitments []ExportedCommitment `json:"ExportedCommitment"`
}

// MonthCumulatedValue is used for window query to fetch the value of a given
// month
type MonthCumulatedValue struct {
	Month int     `json:"Month"`
	Value float64 `json:"Value"`
}

// TwoYearsCommitments is used to fetch commitments of current and previous year
type TwoYearsCommitments struct {
	CurrentYear  []MonthCumulatedValue `json:"CurrentYear"`
	PreviousYear []MonthCumulatedValue `json:"PreviousYear"`
}

// RPLinkedCommitment is used for the renew project linked data project and add
// beneficiary name to the commitment fields
type RPLinkedCommitment struct {
	ID               int64      `json:"ID"`
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
	BeneficiaryID    int64      `json:"BeneficiaryID"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	ActionID         int64      `json:"ActionID"`
	IrisCode         NullString `json:"IrisCode"`
	HousingID        NullInt64  `json:"HousingID"`
	CoproID          NullInt64  `json:"CoproID"`
	RenewProjectID   NullInt64  `json:"RenewProjectID"`
}

// RPLinkedCommitments embeddes an array of RPLinkedCommitment for json export
type RPLinkedCommitments struct {
	Commitments []RPLinkedCommitment `json:"Commitment"`
}

// CoproLinkedCommitment is used for the renew project linked data project and
// add beneficiary name to the commitment fields
type CoproLinkedCommitment struct {
	ID               int64      `json:"ID"`
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
	BeneficiaryID    int64      `json:"BeneficiaryID"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	ActionID         int64      `json:"ActionID"`
	IrisCode         NullString `json:"IrisCode"`
	HousingID        NullInt64  `json:"HousingID"`
	CoproID          NullInt64  `json:"CoproID"`
	RenewProjectID   NullInt64  `json:"RenewProjectID"`
}

// CoproLinkedCommitments embeddes an array of RPLinkedCommitment for json export
type CoproLinkedCommitments struct {
	Commitments []CoproLinkedCommitment `json:"Commitment"`
}

// HousingLinkedCommitment is used for the renew project linked data project and
// add beneficiary name to the commitment fields
type HousingLinkedCommitment struct {
	ID               int64      `json:"ID"`
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
	BeneficiaryID    int64      `json:"BeneficiaryID"`
	BeneficiaryName  string     `json:"BeneficiaryName"`
	ActionID         int64      `json:"ActionID"`
	IrisCode         NullString `json:"IrisCode"`
	HousingID        NullInt64  `json:"HousingID"`
	CoproID          NullInt64  `json:"CoproID"`
	RenewProjectID   NullInt64  `json:"RenewProjectID"`
}

// HousingLinkedCommitments embeddes an array of RPLinkedCommitment for json export
type HousingLinkedCommitments struct {
	Commitments []HousingLinkedCommitment `json:"Commitment"`
}

// Get fetches the results of a paginated commitment query
func (p *PaginatedCommitments) Get(db *sql.DB, c *PaginatedQuery) error {
	var count int64
	commonQryPart := `FROM commitment c 
	JOIN beneficiary b on c.beneficiary_id=b.id
	JOIN budget_action a ON a.id = c.action_id
	JOIN budget_sector s ON s.id=a.sector_id 
	WHERE year >= $1 AND
		(c.name ILIKE $2 OR c.code ILIKE $2 OR c.number::varchar ILIKE $2 
			OR b.name ILIKE $2 OR a.name ILIKE $2 OR iris_code ILIKE $2)`
	if err := db.QueryRow("SELECT count(1) "+commonQryPart, c.Year, "%"+c.Search+"%").
		Scan(&count); err != nil {
		return fmt.Errorf("count query failed %v", err)
	}
	offset, newPage := GetPaginateParams(c.Page, count)
	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,
	c.creation_date,c.modification_date,c.caducity_date,c.name,c.value,c.sold_out,
	c.beneficiary_id,b.name,c.iris_code,a.name,s.name,c.housing_id,
	c.renew_project_id,c.copro_id `+commonQryPart+
		`ORDER BY 1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $3`,
		c.Year, "%"+c.Search+"%", offset)
	if err != nil {
		return err
	}
	var row PaginatedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryID, &row.BeneficiaryName,
			&row.IrisCode, &row.ActionName, &row.Sector, &row.HousingID,
			&row.RenewProjectID, &row.CoproID); err != nil {
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

// GetUnlinked fetches the commitments whose housing_id, copro_id and
// renew_project_id are null and that matches the query using paginated format
func (p *PaginatedCommitments) GetUnlinked(db *sql.DB, c *PaginatedQuery) error {
	var count int64
	commonQryPart := `FROM commitment c
	JOIN beneficiary b on c.beneficiary_id=b.id
	JOIN budget_action a ON a.id=c.action_id
	JOIN budget_sector s ON s.id=a.sector_id 
	WHERE year>=$1 AND housing_id IS NULL AND renew_project_id IS NULL AND
		copro_id IS NULL AND (c.name ILIKE $2 OR c.code ILIKE $2 OR
			c.number::varchar ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2 OR 
			iris_code ILIKE $2) `
	if err := db.QueryRow(`SELECT count(1) `+commonQryPart, c.Year, "%"+c.Search+"%").
		Scan(&count); err != nil {
		return fmt.Errorf("count query failed %v", err)
	}
	offset, newPage := GetPaginateParams(c.Page, count)
	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,
	c.creation_date,c.modification_date,c.caducity_date,c.name,c.value,c.sold_out,
	c.beneficiary_id,b.name,c.iris_code,a.name,s.name `+commonQryPart+
		`ORDER BY 1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $3`,
		c.Year, "%"+c.Search+"%", offset)
	if err != nil {
		return err
	}
	var row PaginatedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryID, &row.BeneficiaryName,
			&row.IrisCode, &row.ActionName, &row.Sector); err != nil {
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

// Get fetches the results of exported commitments
func (e *ExportedCommitments) Get(db *sql.DB, q *ExportQuery) error {
	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,
	c.creation_date,c.modification_date,c.caducity_date,c.name,c.value*0.01,
	c.sold_out, b.name, c.iris_code,a.name,s.name,copro.name,housing.address,
	renew_project.name
	FROM commitment c
	JOIN beneficiary b ON c.beneficiary_id = b.id
	JOIN budget_action a ON a.id = c.action_id
	JOIN budget_sector s ON s.id=a.sector_id
	LEFT JOIN copro ON copro.id = c.copro_id
	LEFT JOIN housing ON housing.id = c.housing_id
	LEFT JOIN renew_project ON renew_project.id = c.renew_project_id
	WHERE year >= $1 AND (c.name ILIKE $2  OR c.number::varchar ILIKE $2 OR 
		c.code ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2)
	ORDER BY 2,6,7,3,4,5`, q.Year, "%"+q.Search+"%")
	if err != nil {
		return err
	}
	var row ExportedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryName, &row.IrisCode,
			&row.ActionName, &row.Sector, &row.CoproName, &row.HousingName,
			&row.RenewProjectName); err != nil {
			return err
		}
		e.ExportedCommitments = append(e.ExportedCommitments, row)
	}
	err = rows.Err()
	if len(e.ExportedCommitments) == 0 {
		e.ExportedCommitments = []ExportedCommitment{}
	}
	return err
}

// GetAll fetches all Commitments from database
func (c *Commitments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,year,code,number,line,creation_date,
	modification_date,caducity_date,name,value,sold_out, beneficiary_id,iris_code, 
	action_id FROM commitment`)
	if err != nil {
		return err
	}
	var row Commitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryID, &row.IrisCode,
			&row.ActionID); err != nil {
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

// Get fetches all Commitments from database linked to a renew
// project whose ID is given
func (c *RPLinkedCommitments) Get(ID int64, db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,
	c.creation_date,c.modification_date,c.caducity_date,c.name,c.value,c.sold_out,
	c.beneficiary_id,b.name,c.iris_code,c.action_id
	FROM commitment c
	JOIN beneficiary b on c.beneficiary_id=b.id 
	WHERE renew_project_id=$1`, ID)
	if err != nil {
		return err
	}
	var row RPLinkedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryID, &row.BeneficiaryName,
			&row.IrisCode, &row.ActionID); err != nil {
			return err
		}
		c.Commitments = append(c.Commitments, row)
	}
	err = rows.Err()
	if len(c.Commitments) == 0 {
		c.Commitments = []RPLinkedCommitment{}
	}
	return err
}

// Get fetches all Commitments from database linked to a copro whose ID is given
func (c *CoproLinkedCommitments) Get(ID int64, db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,
	c.creation_date,c.modification_date,c.caducity_date,c.name,c.value,c.sold_out,
	c.beneficiary_id,b.name,c.iris_code,c.action_id
	FROM commitment c
	JOIN beneficiary b on c.beneficiary_id=b.id 
	WHERE copro_id=$1`, ID)
	if err != nil {
		return err
	}
	var row CoproLinkedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryID, &row.BeneficiaryName,
			&row.IrisCode, &row.ActionID); err != nil {
			return err
		}
		c.Commitments = append(c.Commitments, row)
	}
	err = rows.Err()
	if len(c.Commitments) == 0 {
		c.Commitments = []CoproLinkedCommitment{}
	}
	return err
}

// Get fetches all Commitments from database linked to a hosing whose ID is given
func (c *HousingLinkedCommitments) Get(ID int64, db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT c.id,c.year,c.code,c.number,c.line,
	c.creation_date,c.modification_date,c.caducity_date,c.name,c.value,c.sold_out,
	c.beneficiary_id,b.name,c.iris_code,c.action_id
	FROM commitment c
	JOIN beneficiary b on c.beneficiary_id=b.id 
	WHERE housing_id=$1`, ID)
	if err != nil {
		return err
	}
	var row HousingLinkedCommitment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.Code, &row.Number, &row.Line,
			&row.CreationDate, &row.ModificationDate, &row.CaducityDate, &row.Name,
			&row.Value, &row.SoldOut, &row.BeneficiaryID, &row.BeneficiaryName,
			&row.IrisCode, &row.ActionID); err != nil {
			return err
		}
		c.Commitments = append(c.Commitments, row)
	}
	err = rows.Err()
	if len(c.Commitments) == 0 {
		c.Commitments = []HousingLinkedCommitment{}
	}
	return err
}

// Save insert a batch of CommitmentLine into database
func (c *CommitmentBatch) Save(db *sql.DB) (err error) {
	for i, r := range c.Lines {
		if r.Year < 2009 || r.Number == 0 || r.Line == 0 ||
			r.CreationDate < 20090101 || r.ModificationDate < 20090101 ||
			r.Name == "" || r.BeneficiaryCode == 0 || r.BeneficiaryName == "" ||
			r.Sector == "" {
			return fmt.Errorf("Ligne %d : champs incorrects dans %+v", i+1, r)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_commitment", "year", "code", "number",
		"line", "creation_date", "modification_date", "caducity_date", "name", "value",
		"sold_out", "beneficiary_code", "beneficiary_name", "iris_code", "sector",
		"action_code", "action_name"))
	if err != nil {
		return fmt.Errorf("Statement creation %v", err)
	}
	defer stmt.Close()
	var cd, md, ed time.Time
	for i, r := range c.Lines {
		cd = time.Date(int(r.CreationDate/10000), time.Month(r.CreationDate/100%100),
			int(r.CreationDate%100), 0, 0, 0, 0, time.UTC)
		md = time.Date(int(r.ModificationDate/10000),
			time.Month(r.ModificationDate/100%100), int(r.ModificationDate%100), 0, 0,
			0, 0, time.UTC)
		ed = time.Date(int(r.CaducityDate/10000), time.Month(r.CaducityDate/100%100),
			int(r.CreationDate%100), 0, 0, 0, 0, time.UTC)
		if _, err = stmt.Exec(r.Year, r.Code, r.Number, r.Line, cd, md, ed,
			strings.TrimSpace(r.Name), r.Value, r.SoldOut == "O", r.BeneficiaryCode,
			strings.TrimSpace(r.BeneficiaryName), r.IrisCode,
			strings.TrimSpace(r.Sector), r.ActionCode,
			r.ActionName.TrimSpace()); err != nil {
			tx.Rollback()
			return fmt.Errorf("Ligne %d statement execution %v", i+1, err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	queries := []string{`INSERT INTO beneficiary (code,name)
		SELECT DISTINCT beneficiary_code,beneficiary_name
		FROM temp_commitment
		WHERE beneficiary_code not in (SELECT code from beneficiary)`,
		`INSERT INTO budget_sector (name) SELECT DISTINCT sector
			FROM temp_commitment WHERE sector not in (SELECT name from budget_sector)`,
		`INSERT INTO budget_action (code,name,sector_id)
			SELECT DISTINCT ic.action_code,ic.action_name, s.id
			FROM temp_commitment ic
			LEFT JOIN budget_sector s ON ic.sector = s.name
			WHERE action_code not in (SELECT code from budget_action)`,
		`INSERT INTO commitment (year,code,number,line,creation_date,
			modification_date,caducity_date,name,value,sold_out,beneficiary_id,iris_code,action_id)
			(SELECT ic.year,ic.code,ic.number,ic.line,ic.creation_date,
				ic.modification_date,ic.caducity_date,ic.name,ic.value,ic.sold_out,b.id,
				ic.iris_code,a.id
			 FROM temp_commitment ic
			 JOIN beneficiary b on ic.beneficiary_code=b.code
			 LEFT JOIN budget_action a on ic.action_code = a.code
			 WHERE (ic.year,ic.code,ic.number,ic.line,ic.creation_date,
				ic.modification_date,ic.name, ic.value) NOT IN
					(SELECT year,code,number,line,creation_date,modification_date,
						name,value FROM commitment))`,
		`DELETE FROM temp_commitment`}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %v", i, err)
		}
	}
	return tx.Commit()
}

// Get fetches all commitments per year for the current and the previous years
func (t *TwoYearsCommitments) Get(db *sql.DB) error {
	query := `WITH cmt_month as 
	(SELECT MAX(EXTRACT(month FROM creation_date))::int max_month
  	FROM commitment WHERE year=$1)
	SELECT cmt.m,SUM(0.01 * cmt.v) OVER (ORDER BY m) FROM
		(SELECT q.m as m,COALESCE(sum_cmt.v,0) v FROM
			(SELECT GENERATE_SERIES(1,max_month) m FROM cmt_month) q
			LEFT OUTER JOIN
			(SELECT EXTRACT(month FROM creation_date)::int m,SUM(value)::bigint v
			FROM commitment WHERE year=$1 GROUP BY 1) sum_cmt
		ON sum_cmt.m=q.m) cmt;`
	actualYear := time.Now().Year()
	rows, err := db.Query(query, actualYear)
	if err != nil {
		return err
	}
	defer rows.Close()
	var row MonthCumulatedValue
	for rows.Next() {
		if err = rows.Scan(&row.Month, &row.Value); err != nil {
			return err
		}
		t.CurrentYear = append(t.CurrentYear, row)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	if len(t.CurrentYear) == 0 {
		t.CurrentYear = []MonthCumulatedValue{}
	}
	rows, err = db.Query(query, actualYear-1)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.Month, &row.Value); err != nil {
			return err
		}
		t.PreviousYear = append(t.PreviousYear, row)
	}
	err = rows.Err()
	if len(t.PreviousYear) == 0 {
		t.PreviousYear = []MonthCumulatedValue{}
	}
	return err
}
