package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
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
}

// ExportedPayments embeddes an array of ExportedPayment for json export
type ExportedPayments struct {
	ExportedPayments []ExportedPayment `json:"ExportedPayment"`
}

// TwoYearsPayments is used to fetch payments of current and previous year
type TwoYearsPayments struct {
	CurrentYear  []MonthCumulatedValue `json:"CurrentYear"`
	PreviousYear []MonthCumulatedValue `json:"PreviousYear"`
}

// GetAll fetches all Payments from database
func (p *Payments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,commitment_id,commitment_year,commitment_code,
	commitment_number,commitment_line,year,creation_date,modification_date,
	number, value FROM payment`)
	if err != nil {
		return err
	}
	var row Payment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CommitmentYear,
			&row.CommitmentCode, &row.CommitmentNumber, &row.CommitmentLine, &row.Year,
			&row.CreationDate, &row.ModificationDate, &row.Number, &row.Value); err != nil {
			return err
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if len(p.Payments) == 0 {
		p.Payments = []Payment{}
	}
	return err
}

// Save insert a batch of PaymentLine into database
func (p *PaymentBatch) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_payment (commitment_year,commitment_code,
		commitment_number,commitment_line,year,creation_date,modification_date,number, value)
		VALUES ($1,$2,$3,$4,$5,make_date($6,$7,$8),make_date($9,$10,$11),$12,$13)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range p.Lines {
		if r.CommitmentYear == 0 || r.CommitmentCode == "" || r.CommitmentNumber == 0 ||
			r.CommitmentLine == 0 || r.Year == 0 || r.CreationDate < 20090101 ||
			r.ModificationDate < 20090101 || r.Number == 0 {
			tx.Rollback()
			return fmt.Errorf("Champ incorrect dans %+v", r)
		}
		if _, err = stmt.Exec(r.CommitmentYear, r.CommitmentCode, r.CommitmentNumber,
			r.CommitmentLine, r.Year, r.CreationDate/10000, (r.CreationDate/100)%100,
			r.CreationDate%100, r.ModificationDate/10000, (r.ModificationDate/100)%100,
			r.ModificationDate%100, r.Number, r.Value); err != nil {
			tx.Rollback()
			return err
		}
	}
	queries := []string{`UPDATE payment SET value=t.value, 
		modification_date=t.modification_date
	FROM temp_payment t WHERE t.commitment_year=payment.commitment_year AND 
		t.commitment_code=payment.commitment_code AND
		t.commitment_number=payment.commitment_number AND
		t.commitment_line=payment.commitment_line AND t.year=payment.year AND 
		t.creation_date=payment.creation_date AND
		t.number=payment.number`,
		`INSERT INTO payment (commitment_id,commitment_year,commitment_code,
			commitment_number,commitment_line,year,creation_date,modification_date,
			number, value)
		SELECT c.id,t.commitment_year,t.commitment_code,t.commitment_number,
			t.commitment_line,t.year,t.creation_date,t.modification_date,t.number,t.value 
			FROM temp_payment t
			LEFT JOIN cumulated_commitment c 
				ON t.commitment_year=c.year AND t.commitment_code=c.code
				AND t.commitment_number=c.number
			WHERE (t.commitment_year,t.commitment_code,t.commitment_number,
				t.commitment_line,t.year,t.creation_date,t.modification_date) 
			NOT IN (SELECT DISTINCT commitment_year,commitment_code,commitment_number,
				commitment_line,year,creation_date,modification_date from payment)`,
		`DELETE FROM temp_payment`,
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

// Get fetches all paginated payments from database that match the paginated query
func (p *PaginatedPayments) Get(db *sql.DB, q *PaginatedQuery) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM payment p 
		LEFT JOIN cumulated_commitment c on p.commitment_id=c.id
		JOIN budget_action a ON a.id = c.action_id
		JOIN budget_sector s ON s.id=a.sector_id 
		JOIN beneficiary b ON c.beneficiary_id = b.id
		WHERE p.year >= $1 AND
			(c.name ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2)`, q.Year, "%"+q.Search+"%").
		Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT p.id,p.year,p.creation_date,p.value, p.number,
	c.creation_date, c.name, c.value, b.name, s.name, a.name FROM payment p
	LEFT JOIN cumulated_commitment c ON p.commitment_id = c.id
	JOIN beneficiary b ON c.beneficiary_id = b.id
	JOIN budget_action a ON a.id = c.action_id
	JOIN budget_sector s ON s.id=a.sector_id 
	WHERE p.year >= $1 AND (c.name ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2)
	ORDER BY 2,5,3 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $3`,
		q.Year, "%"+q.Search+"%", offset)
	if err != nil {
		return err
	}
	var row PaginatedPayment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CreationDate, &row.Value,
			&row.Number, &row.CommitmentDate, &row.CommitmentName, &row.CommitmentValue,
			&row.Beneficiary, &row.Sector, &row.ActionName); err != nil {
			return err
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	if len(p.Payments) == 0 {
		p.Payments = []PaginatedPayment{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}

// Get fetches all exported payments from database that match the export query
func (p *ExportedPayments) Get(db *sql.DB, q *ExportQuery) error {
	rows, err := db.Query(`SELECT p.id,p.year,p.creation_date,p.modification_date,
	p.number, p.value * 0.01, p.commitment_year, p.commitment_code, p.commitment_number,
	c.creation_date, c.value * 0.01,c.name, b.name, s.name, a.name
	FROM payment p
	LEFT JOIN cumulated_commitment c ON p.commitment_id = c.id
	JOIN beneficiary b ON c.beneficiary_id = b.id
	JOIN budget_action a ON a.id = c.action_id
	JOIN budget_sector s ON s.id=a.sector_id 
	WHERE p.year >= $1 AND (c.name ILIKE $2 OR b.name ILIKE $2 OR a.name ILIKE $2)
	ORDER BY 1 `, q.Year, "%"+q.Search+"%")
	if err != nil {
		return err
	}
	var row ExportedPayment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Year, &row.CreationDate, &row.ModificationDate,
			&row.Number, &row.Value, &row.CommitmentYear, &row.CommitmentCode,
			&row.CommitmentNumber, &row.CommitmentCreationDate, &row.CommitmentValue,
			&row.CommitmentName, &row.BeneficiaryName, &row.Sector, &row.ActionName); err != nil {
			return err
		}
		p.ExportedPayments = append(p.ExportedPayments, row)
	}
	err = rows.Err()
	if len(p.ExportedPayments) == 0 {
		p.ExportedPayments = []ExportedPayment{}
	}
	return err
}

// Get fetches all payments per year for the current and the previous years
func (t *TwoYearsPayments) Get(db *sql.DB) error {
	var row MonthCumulatedValue
	rows, err := db.Query(`SELECT q.m,SUM(q.v) OVER (ORDER BY m) FROM
	(SELECT EXTRACT(MONTH FROM modification_date) as m,sum(value*0.01) as v
	FROM payment WHERE year=extract(year FROM CURRENT_DATE) GROUP BY 1) q`)
	if err != nil {
		return err
	}
	defer rows.Close()
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
	rows, err = db.Query(`SELECT q.m,SUM(q.v) OVER (ORDER BY m) FROM
	(SELECT EXTRACT(MONTH FROM modification_date) as m,sum(value*0.01) as v
	FROM payment WHERE year=extract(year FROM CURRENT_DATE)-1 GROUP BY 1) q`)
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
