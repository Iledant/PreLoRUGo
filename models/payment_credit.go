package models

import (
	"database/sql"
	"fmt"
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
		return err
	}
	defer rows.Close()
	var row PaymentCredit
	for rows.Next() {
		if err = rows.Scan(&row.Year, &row.Chapter, &row.Function, &row.Primitive,
			&row.Reported, &row.Added, &row.Modified, &row.Movement); err != nil {
			return err
		}
		p.Lines = append(p.Lines, row)
	}
	if err = rows.Err(); err != nil {
		return err
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
		return err
	}
	for _, l := range p.Lines {
		if _, err = tx.Exec(`INSERT INTO temp_payment_credit (chapter,
			function,primitive,reported,added,modified,
			movement) VALUES($1,$2,$3,$4,$5,$6,$7)`, l.Chapter, l.Function,
			l.Primitive, l.Reported, l.Added, l.Modified,
			l.Movement); err != nil {
			tx.Rollback()
			return fmt.Errorf("temp insert %v", err)
		}
	}
	if _, err = tx.Exec(`UPDATE payment_credit SET year=$1,primitive=t.primitive,
		reported=t.reported,added=t.added,modified=t.modified,movement=t.movement
		FROM (SELECT * FROM temp_payment_credit WHERE (chapter,function) IN
					(SELECT chapter,function FROM payment_credit)) t`, year); err != nil {
		tx.Rollback()
		return fmt.Errorf("update %v", err)
	}
	if _, err = tx.Exec(`INSERT INTO payment_credit (year,chapter,function,
		primitive,reported,added,modified,movement)
		(SELECT $1, chapter, function,primitive,reported,added,modified,movement
			FROM temp_payment_credit
			WHERE (chapter,function) NOT IN
			(SELECT chapter,function FROM payment_credit))`, year); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert %v", err)
	}
	if _, err = tx.Exec(`DELETE FROM temp_payment_credit`); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete %v", err)
	}
	tx.Commit()
	return nil
}
