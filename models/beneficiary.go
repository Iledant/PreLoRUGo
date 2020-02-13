package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// Beneficiary model
type Beneficiary struct {
	ID   int64  `json:"ID"`
	Code int64  `json:"Code"`
	Name string `json:"Name"`
}

// Beneficiaries embeddes an array of Beneficiary for json export
type Beneficiaries struct {
	Beneficiaries []Beneficiary `json:"Beneficiary"`
}

// PaginatedBeneficiaries embeddes an array of Beneficiary for paginated
// display
type PaginatedBeneficiaries struct {
	Beneficiaries []Beneficiary `json:"Beneficiary"`
	Page          int64         `json:"Page"`
	ItemsCount    int64         `json:"ItemsCount"`
}

// GetAll fetches all Beneficiaries from database
func (b *Beneficiaries) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,code,name FROM beneficiary`)
	if err != nil {
		return err
	}
	var row Beneficiary
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Code, &row.Name); err != nil {
			return err
		}
		b.Beneficiaries = append(b.Beneficiaries, row)
	}
	err = rows.Err()
	if len(b.Beneficiaries) == 0 {
		b.Beneficiaries = []Beneficiary{}
	}
	return err
}

// Get fetches beneficiaries from database according to PaginatedQuery where
// only Page and Search fields are used
func (p *PaginatedBeneficiaries) Get(db *sql.DB, q *PaginatedQuery) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM beneficiary 
		WHERE name ILIKE $1 OR code::varchar ILIKE $1`, "%"+q.Search+"%").
		Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT id,code,name FROM beneficiary b
	WHERE name ILIKE $1 OR code::varchar ILIKE $1
	ORDER BY 2,1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $2`,
		"%"+q.Search+"%", offset)
	if err != nil {
		return err
	}
	var row Beneficiary
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Code, &row.Name); err != nil {
			return err
		}
		p.Beneficiaries = append(p.Beneficiaries, row)
	}
	err = rows.Err()
	if len(p.Beneficiaries) == 0 {
		p.Beneficiaries = []Beneficiary{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}

// Validate checks if the fields matches the database constraints
func (b *Beneficiary) Validate() error {
	if b.Name == "" {
		return errors.New("name vide")
	}
	return nil
}

// Create insert a new beneficiary into the database
func (b *Beneficiary) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO beneficiary(code,name) VALUES($1,$2)
	RETURNING ID`, b.Code, b.Name).Scan(&b.ID)
}

// Update modifies a beneficiary in the database
func (b *Beneficiary) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE beneficiary SET code=$1,name=$2 WHERE id=$3`,
		b.Code, b.Name, b.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsaffected %v", err)
	}
	if count != 1 {
		return fmt.Errorf("bénéficiaire introuvable")
	}
	return nil
}

// Delete removes a beneficiary from database
func (b *Beneficiary) Delete(db *sql.DB) error {
	res, err := db.Exec(`DELETE FROM beneficiary WHERE id=$1`, b.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsaffected %v", err)
	}
	if count != 1 {
		return fmt.Errorf("bénéficiaire introuvable")
	}
	return nil
}
