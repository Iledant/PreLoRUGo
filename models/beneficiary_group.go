package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// BeneficiaryGroup model
type BeneficiaryGroup struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

// BeneficiaryGroups embeddes an array of BeneficiaryGroup for json export and
// dedicated queries
type BeneficiaryGroups struct {
	Lines []BeneficiaryGroup `json:"BeneficiaryGroup"`
}

// Valid checks if fields are correctly filled
func (b *BeneficiaryGroup) Valid() error {
	if b.Name == "" {
		return errors.New("name vide")
	}
	return nil
}

// Get fetches all beneficiary groups from database
func (b *BeneficiaryGroups) Get(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,name FROM beneficiary_group`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var bb BeneficiaryGroup
	for rows.Next() {
		if err = rows.Scan(&bb.ID, &bb.Name); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		b.Lines = append(b.Lines, bb)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(b.Lines) == 0 {
		b.Lines = []BeneficiaryGroup{}
	}
	return nil
}

// Create insert a new beneficiary group into database
func (b *BeneficiaryGroup) Create(db *sql.DB) error {
	if err := db.QueryRow(`INSERT INTO beneficiary_group (name) VALUES ($1)
	RETURNING ID`, b.Name).Scan(&b.ID); err != nil {
		return fmt.Errorf("insert %v", err)
	}
	return nil
}

// Delete remove a beneficiary from database using database cascading to also
// remove beneficiary_belong
func (b *BeneficiaryGroup) Delete(db *sql.DB) error {
	res, err := db.Exec(`DELETE FROM beneficiary_group WHERE id=$1`, b.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", res)
	}
	if count != 1 {
		return fmt.Errorf("groupe non trouvé")
	}
	return nil
}

// Update change beneficiary group name
func (b *BeneficiaryGroup) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE beneficiary_group SET name=$1 WHERE id=$2`,
		b.Name, b.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", res)
	}
	if count != 1 {
		return fmt.Errorf("groupe non trouvé")
	}
	return nil
}

// Set replace the beneficiaries ID of a beneficiary group
func (b *BeneficiaryGroup) Set(IDs []int64, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}
	if _, err := tx.Exec(`DELETE FROM beneficiary_belong WHERE group_id=$1`,
		b.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("beneficiary_belong", "beneficiary_id",
		"group_id"))
	if err != nil {
		return fmt.Errorf("prepare %v", err)
	}
	defer stmt.Close()
	for _, ID := range IDs {
		if _, err := stmt.Exec(ID, b.ID); err != nil {
			tx.Rollback()
			return fmt.Errorf("stmt %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	tx.Commit()
	return nil
}

// GroupGet fetches the beneficiaries linked to a beneficiary group
func (b *Beneficiaries) GroupGet(ID int64, db *sql.DB) error {
	rows, err := db.Query(`SELECT b.id,b.code,b.name FROM beneficiary_belong bb
	JOIN beneficiary b ON b.id=bb.beneficiary_id WHERE bb.group_id=$1`, ID)
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
