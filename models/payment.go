package models

import (
	"database/sql"
	"errors"
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
}

// PaymentBatch embeddes an array of PaymentLine for json export
type PaymentBatch struct {
	Lines []PaymentLine `json:"Payment"`
}

// GetAll fetches all Payments from database
func (p *Payments) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,commitment_id,commitment_year,commitment_code,
	commitment_number,commitment_line,year,creation_date,modification_date,value FROM payment`)
	if err != nil {
		return err
	}
	var row Payment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CommitmentYear,
			&row.CommitmentCode, &row.CommitmentNumber, &row.CommitmentLine, &row.Year,
			&row.CreationDate, &row.ModificationDate, &row.Value); err != nil {
			return err
		}
		p.Payments = append(p.Payments, row)
	}
	err = rows.Err()
	return err
}

// Save insert a batch of PaymentLine into database
func (p *PaymentBatch) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_payment (commitment_year,commitment_code,
		commitment_number,commitment_line,year,creation_date,modification_date,value)
		VALUES ($1,$2,$3,$4,$5,make_date($6,$7,$8),make_date($9,$10,$11),$12)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range p.Lines {
		if r.CommitmentYear == 0 || r.CommitmentCode == "" || r.CommitmentNumber == 0 ||
			r.CommitmentLine == 0 || r.Year == 0 || r.CreationDate < 20090101 ||
			r.ModificationDate < 20090101 {
			tx.Rollback()
			return errors.New("Champs incorrects")
		}
		if _, err = stmt.Exec(r.CommitmentYear, r.CommitmentCode, r.CommitmentNumber,
			r.CommitmentLine, r.Year, r.CreationDate/10000, (r.CreationDate/100)%100,
			r.CreationDate%100, r.ModificationDate/10000, (r.ModificationDate/100)%100,
			r.ModificationDate%100, r.Value); err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`UPDATE payment SET value=t.value, modification_date=t.modification_date
	FROM temp_payment t WHERE t.commitment_year=payment.commitment_year AND 
		t.commitment_code=payment.commitment_code AND
		t.commitment_number=payment.commitment_number AND
		t.commitment_line=payment.commitment_line AND t.year=payment.year AND 
		t.creation_date=payment.creation_date`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO payment (commitment_id,commitment_year,commitment_code,
		commitment_number,commitment_line,year,creation_date,modification_date,value)
	SELECT c.id,t.commitment_year,t.commitment_code,t.commitment_number,
		t.commitment_line,t.year,t.creation_date,t.modification_date,t.value 
		FROM temp_payment t
		LEFT JOIN commitment c ON t.commitment_year=c.year AND t.commitment_code=c.code
			AND t.commitment_number=c.number AND t.commitment_line=c.line
		WHERE (t.commitment_year,t.commitment_code,t.commitment_number,t.commitment_line,
			t.year,t.creation_date,t.modification_date) 
		NOT IN (SELECT DISTINCT commitment_year,commitment_code,commitment_number,
			commitment_line,year,creation_date,modification_date from payment)`)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
