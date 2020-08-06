package models

import (
	"database/sql"
	"fmt"
	"time"
)

// PaymentCredit model
type PaymentCredit struct {
	Year      int64 `json:"Year"`
	Chapter   int64 `json:"Chapter"`
	Function  int64 `json:"Function"`
	Primitive int64 `json:"Primitive"`
	Reported  int64 `json:"Reported"`
	Added     int64 `json:"Added"`
	Modified  int64 `json:"Modified"`
	Movement  int64 `json:"Movement"`
}

// PaymentCredits embeddes an array of PaymentCredit for json export
type PaymentCredits struct {
	Lines []PaymentCredit `json:"PaymentCredit"`
}

// PaymentCreditLine is used to decode one line of PaymentCreditBatch
type PaymentCreditLine struct {
	Chapter   int64 `json:"Chapter"`
	Function  int64 `json:"Function"`
	Primitive int64 `json:"Primitive"`
	Reported  int64 `json:"Reported"`
	Added     int64 `json:"Added"`
	Modified  int64 `json:"Modified"`
	Movement  int64 `json:"Movement"`
}

// PaymentCreditSum is used to calculate the total of payment credit of a current year
type PaymentCreditSum struct {
	Sum NullFloat64 `json:"PaymentCreditSum"`
}

// PaymentCreditBatch embeddes an array of PaumentCreditLine for batch import
type PaymentCreditBatch struct {
	Lines []PaymentCreditLine `json:"PaymentCredit"`
}

// GetAll fetches all PaymentCredits of a year from database
func (p *PaymentCredits) GetAll(year int, db *sql.DB) error {
	rows, err := db.Query(`SELECT year,chapter,function,primitive,reported,added,
		modified,movement
	 FROM payment_credit WHERE year=$1`, year)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	defer rows.Close()
	var row PaymentCredit
	for rows.Next() {
		if err = rows.Scan(&row.Year, &row.Chapter, &row.Function, &row.Primitive,
			&row.Reported, &row.Added, &row.Modified, &row.Movement); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Lines = append(p.Lines, row)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Lines) == 0 {
		p.Lines = []PaymentCredit{}
	}
	return nil
}

// Save import a batch of payment credits into database
func (p *PaymentCreditBatch) Save(year int64, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}
	if _, err = tx.Exec(`DELETE FROM payment_credit WHERE year=$1`, year); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete %v", err)
	}
	for i, l := range p.Lines {
		if _, err = tx.Exec(`INSERT INTO payment_credit (year,chapter,function,
			primitive,reported,added,modified,movement) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
			year, l.Chapter, l.Function, l.Primitive, l.Reported, l.Added, l.Modified,
			l.Movement); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert %d %v", i, err)
		}
	}
	return tx.Commit()
}

// Get fetches the payment credit sum of the current year
func (p *PaymentCreditSum) Get(db *sql.DB) error {
	q := `SELECT sum(primitive+reported+added+modified+movement)*0.01 from payment_credit
	where year=$1 and chapter='905' and function<>52`
	year := time.Now().Year()
	if err := db.QueryRow(q, year).Scan(&p.Sum); err != nil {
		return fmt.Errorf("select %v", err)
	}
	return nil
}
