package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
)

// Payment model
type Payment struct {
	ID               int64     `json:"ID"`
	CommitmentID     NullInt64 `json:"CommitmentID"`
	CommitmentYear   int64     `json:"CommitmentYear"`
	CommitmentCode   string    `json:"CommitmentCode"`
	CommitmentNumber int64     `json:"CommitmentNumber"`
	CommitmentLine   int64     `json:"CommitmentLine"`
	Year             int64     `json:"Year"`
	CreationDate     time.Time `json:"CreationDate"`
	ModificationDate time.Time `json:"ModificationDate"`
	Number           int64     `json:"Number"`
	Value            int64     `json:"Value"`
	ReceiptDate      NullTime  `json:"ReceiptDate"`
}

// Payments embeddes an array of Payment for json export
type Payments struct {
	Payments []Payment `json:"Payment"`
}

// PaymentLine is used to decode a line of Payment batch
type PaymentLine struct {
	CommitmentYear   int64  `json:"CommitmentYear"`
	CommitmentCode   string `json:"CommitmentCode"`
	CommitmentNumber int64  `json:"CommitmentNumber"`
	CommitmentLine   int64  `json:"CommitmentLine"`
	Year             int64  `json:"Year"`
	CreationDate     int    `json:"CreationDate"`
	ModificationDate int    `json:"ModificationDate"`
	Value            int64  `json:"Value"`
	Number           int64  `json:"Number"`
	ReceiptDate      int    `json:"ReceiptDate"`
}

// PaymentBatch embeddes an array of PaymentLine for json export
type PaymentBatch struct {
	Lines []PaymentLine `json:"Payment"`
}

// PaginatedPayment is used for paginated request to fetch some payments that
// match a search pattern using PaginatedQuery
type PaginatedPayment struct {
	ID              int64      `json:"ID"`
	Year            int64      `json:"Year"`
	CreationDate    time.Time  `json:"CreationDate"`
	Value           int64      `json:"Value"`
	Number          int64      `json:"Number"`
	CommitmentDate  NullTime   `json:"CommitmentDate"`
	CommitmentName  NullString `json:"CommitmentName"`
	CommitmentValue NullInt64  `json:"CommitmentValue"`
	Beneficiary     NullString `json:"Beneficiary"`
	Sector          NullString `json:"Sector"`
	ActionName      NullString `json:"ActionName"`
	ReceiptDate     NullTime   `json:"ReceiptDate"`
}

// PaginatedPayments embeddes an array of PaginatedPayment for json export with
// paginated informations
type PaginatedPayments struct {
	Payments   []PaginatedPayment `json:"Payments"`
	Page       int64              `json:"Page"`
	ItemsCount int64              `json:"ItemsCount"`
}

// ExportedPayment is used for the excel export query to fetch payment with all
// fields according to a certain search pattern
type ExportedPayment struct {
	ID                     int64       `json:"ID"`
	Year                   int64       `json:"Year"`
	CreationDate           time.Time   `json:"CreationDate"`
	ModificationDate       time.Time   `json:"ModificationDate"`
	Number                 int64       `json:"Number"`
	Value                  float64     `json:"Value"`
	CommitmentYear         int64       `json:"CommitmentYear"`
	CommitmentCode         string      `json:"CommitmentCode"`
	CommitmentNumber       int64       `json:"CommitmentNumber"`
	CommitmentCreationDate NullTime    `json:"CommitmentCreationDate"`
	CommitmentValue        NullFloat64 `json:"CommitmentValue"`
	CommitmentName         NullString  `json:"CommitmentName"`
	BeneficiaryName        NullString  `json:"BeneficiaryName"`
	Sector                 NullString  `json:"Sector"`
	ActionName             NullString  `json:"ActionName"`
	ReceiptDate            NullTime    `json:"ReceiptDate"`
}

// ExportedPayments embeddes an array of ExportedPayment for json export
type ExportedPayments struct {
	ExportedPayments []ExportedPayment `json:"ExportedPayment"`
}

// SectorPayment is used to fetch cumulated payments per sector
type SectorPayment struct {
	Month  int64   `json:"Month"`
	Sector string  `json:"Sector"`
	Value  float64 `json:"Value"`
}

// TwoYearsPayments is used to fetch payments of current and previous year
type TwoYearsPayments struct {
	CurrentYear  []SectorPayment `json:"CurrentYear"`
	PreviousYear []SectorPayment `json:"PreviousYear"`
}

// GetAll fetches all Payments FROM database
func (p *Payments) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,commitment_id,commitment_year,commitment_code,
	commitment_number,commitment_line,year,creation_date,modification_date,
	number, value, receipt_date FROM payment`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row Payment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CommitmentYear,
			&row.CommitmentCode, &row.CommitmentNumber, &row.CommitmentLine, &row.Year,
			&row.CreationDate, &row.ModificationDate, &row.Number, &row.Value,
			&row.ReceiptDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Payments) == 0 {
		p.Payments = []Payment{}
	}
	return nil
}

const commonLinkedQuery = `SELECT p.id,p.commitment_id,p.commitment_year,
	p.commitment_code,p.commitment_number,p.commitment_line,p.year,
	p.creation_date,p.modification_date,p.number,p.value,p.receipt_date
FROM payment p
JOIN commitment c ON p.commitment_id = c.id WHERE `

// GetLinkedToRenewProject fetches all Payments FROM database
func (p *Payments) GetLinkedToRenewProject(ID int64, db *sql.DB) error {
	rows, err := db.Query(commonLinkedQuery+"c.renew_project_id=$1", ID)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row Payment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CommitmentYear,
			&row.CommitmentCode, &row.CommitmentNumber, &row.CommitmentLine, &row.Year,
			&row.CreationDate, &row.ModificationDate, &row.Number, &row.Value,
			&row.ReceiptDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Payments) == 0 {
		p.Payments = []Payment{}
	}
	return nil
}

// GetLinkedToCopro fetches all Payments FROM database
func (p *Payments) GetLinkedToCopro(ID int64, db *sql.DB) error {
	rows, err := db.Query(commonLinkedQuery+"c.copro_id=$1", ID)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row Payment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CommitmentYear,
			&row.CommitmentCode, &row.CommitmentNumber, &row.CommitmentLine, &row.Year,
			&row.CreationDate, &row.ModificationDate, &row.Number, &row.Value,
			&row.ReceiptDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Payments) == 0 {
		p.Payments = []Payment{}
	}
	return nil
}

// GetLinkedToHousing fetches all Payments FROM database
func (p *Payments) GetLinkedToHousing(ID int64, db *sql.DB) error {
	rows, err := db.Query(commonLinkedQuery+"c.housing_id=$1", ID)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row Payment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CommitmentYear,
			&row.CommitmentCode, &row.CommitmentNumber, &row.CommitmentLine, &row.Year,
			&row.CreationDate, &row.ModificationDate, &row.Number, &row.Value,
			&row.ReceiptDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Payments) == 0 {
		p.Payments = []Payment{}
	}
	return nil
}

// Save insert a batch of PaymentLine into database
func (p *PaymentBatch) Save(db *sql.DB) error {
	for _, r := range p.Lines {
		if r.CommitmentYear == 0 || r.CommitmentCode == "" || r.CommitmentNumber == 0 ||
			r.CommitmentLine == 0 || r.Year == 0 || r.CreationDate < 20090101 ||
			r.ModificationDate < 20090101 || r.Number == 0 {
			return fmt.Errorf("Champ incorrect dans %+v", r)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_payment", "commitment_year",
		"commitment_code", "commitment_number", "commitment_line", "year",
		"creation_date", "modification_date", "number", "value", "receipt_date"))
	if err != nil {
		return fmt.Errorf("copy in %v", err)
	}
	defer stmt.Close()
	var (
		cd, md time.Time
		rd     NullTime
	)
	for _, r := range p.Lines {
		cd = time.Date(int(r.CreationDate/10000), time.Month(r.CreationDate/100%100),
			int(r.CreationDate%100), 0, 0, 0, 0, time.UTC)
		md = time.Date(int(r.ModificationDate/10000),
			time.Month(r.ModificationDate/100%100), int(r.ModificationDate%100), 0, 0,
			0, 0, time.UTC)
		if r.ReceiptDate == 0 {
			rd.Valid = false
		} else {
			rd.Valid = true
			rd.Time = time.Date(int(r.ReceiptDate/10000),
				time.Month(r.ReceiptDate/100%100), int(r.ReceiptDate%100), 0, 0, 0, 0,
				time.UTC)
		}
		if _, err = stmt.Exec(r.CommitmentYear, r.CommitmentCode, r.CommitmentNumber,
			r.CommitmentLine, r.Year, cd, md, r.Number, r.Value, rd); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement exec %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{`UPDATE payment SET value=t.value, 
		modification_date=t.modification_date, receipt_date=t.receipt_date
	FROM temp_payment t WHERE t.commitment_year=payment.commitment_year AND 
		t.commitment_code=payment.commitment_code AND
		t.commitment_number=payment.commitment_number AND
		t.commitment_line=payment.commitment_line AND t.year=payment.year AND 
		t.creation_date=payment.creation_date AND
		t.number=payment.number`,
		`INSERT INTO payment (commitment_id,commitment_year,commitment_code,
			commitment_number,commitment_line,year,creation_date,modification_date,
			number, value, receipt_date)
		SELECT c.id,t.commitment_year,t.commitment_code,t.commitment_number,
			t.commitment_line,t.year,t.creation_date,t.modification_date,t.number,
			t.value,t.receipt_date
			FROM temp_payment t
			LEFT JOIN cumulated_commitment c 
				ON t.commitment_year=c.year AND t.commitment_code=c.code
				AND t.commitment_number=c.number
			WHERE (t.commitment_year,t.commitment_code,t.commitment_number,
				t.commitment_line,t.year,t.creation_date,t.modification_date) 
			NOT IN (SELECT DISTINCT commitment_year,commitment_code,commitment_number,
				commitment_line,year,creation_date,modification_date FROM payment)`,
		`DELETE FROM temp_payment`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	return tx.Commit()
}

// Get fetches all paginated payments FROM database that match the paginated query
func (p *PaginatedPayments) Get(db *sql.DB, q *PaginatedQuery) error {
	var count int64
	const commonPmtQry = ` FROM payment p 
	LEFT JOIN cumulated_commitment c on p.commitment_id=c.id
	JOIN budget_action a ON a.id = c.action_id
	JOIN budget_sector s ON s.id=a.sector_id 
	JOIN beneficiary b ON c.beneficiary_id = b.id
	WHERE p.year >= $1 AND
		(c.name ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2)`

	if err := db.QueryRow(`SELECT count(1)`+commonPmtQry, q.Year, "%"+q.Search+"%").
		Scan(&count); err != nil {
		return fmt.Errorf("select count %v", err)
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT p.id,p.year,p.creation_date,p.value,p.number,
	c.creation_date,c.name,c.value,b.name,s.name,a.name,p.receipt_date`+
		commonPmtQry+` ORDER BY 2,5,3 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $3`,
		q.Year, "%"+q.Search+"%", offset)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row PaginatedPayment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CreationDate, &row.Value,
			&row.Number, &row.CommitmentDate, &row.CommitmentName, &row.CommitmentValue,
			&row.Beneficiary, &row.Sector, &row.ActionName, &row.ReceiptDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Payments) == 0 {
		p.Payments = []PaginatedPayment{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return nil
}

// Get fetches all exported payments FROM database that match the export query
func (p *ExportedPayments) Get(db *sql.DB, q *ExportQuery) error {
	rows, err := db.Query(`SELECT p.id,p.year,p.creation_date,p.modification_date,
	p.number,p.value*0.01,p.commitment_year,p.commitment_code,p.commitment_number,
	c.creation_date,c.value*0.01,c.name,b.name,s.name,a.name,p.receipt_date
	FROM payment p
	LEFT JOIN cumulated_commitment c ON p.commitment_id = c.id
	JOIN beneficiary b ON c.beneficiary_id = b.id
	JOIN budget_action a ON a.id = c.action_id
	JOIN budget_sector s ON s.id=a.sector_id 
	WHERE p.year >= $1 AND (c.name ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2)
	ORDER BY 1 `, q.Year, "%"+q.Search+"%")
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row ExportedPayment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CreationDate, &row.ModificationDate,
			&row.Number, &row.Value, &row.CommitmentYear, &row.CommitmentCode,
			&row.CommitmentNumber, &row.CommitmentCreationDate, &row.CommitmentValue,
			&row.CommitmentName, &row.BeneficiaryName, &row.Sector, &row.ActionName,
			&row.ReceiptDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.ExportedPayments = append(p.ExportedPayments, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.ExportedPayments) == 0 {
		p.ExportedPayments = []ExportedPayment{}
	}
	return nil
}

// Get fetches all payments per year for the current and the previous years
func (t *TwoYearsPayments) Get(db *sql.DB) error {
	query := `WITH pmt_month as (
		SELECT max(extract(month FROM creation_date))::int max_month
		FROM payment WHERE year=$1)
	SELECT pmt.m,name,sum(pmt.v) OVER (PARTITION BY name ORDER BY m) FROM
	(SELECT q.m as m,sum_pmt.name,COALESCE(sum_pmt.v,0)*0.00000001 v FROM
	(SELECT generate_series(1,max_month) m FROM pmt_month) q
	LEFT OUTER JOIN
	(SELECT EXTRACT(month FROM p.creation_date)::int m,s.name,SUM(p.value)::bigint v
	FROM payment p
	JOIN commitment c on p.commitment_id=c.id
	JOIN budget_action ba ON c.action_id=ba.id
	JOIN budget_sector s ON ba.sector_id=s.id
	WHERE p.year=$1
	GROUP BY 1,2) sum_pmt
	ON sum_pmt.m=q.m) pmt;`
	actualYear := time.Now().Year()
	var row SectorPayment
	rows, err := db.Query(query, actualYear)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.Month, &row.Sector, &row.Value); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		t.CurrentYear = append(t.CurrentYear, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(t.CurrentYear) == 0 {
		t.CurrentYear = []SectorPayment{}
	}
	rows, err = db.Query(query, actualYear-1)
	if err != nil {
		return fmt.Errorf("select last year %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.Month, &row.Sector, &row.Value); err != nil {
			return fmt.Errorf("scan last year %v", err)
		}
		t.PreviousYear = append(t.PreviousYear, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(t.PreviousYear) == 0 {
		t.PreviousYear = []SectorPayment{}
	}
	return nil
}
