package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// CommitmentLink is used to embed link or unlink data
type CommitmentLink struct {
	DestID int    `json:"DestID"`
	IDs    []int  `json:"IDs"`
	Type   string `json:"Type"`
}

// CommitmentUnlink is used to embed an array of commitment's IDs
type CommitmentUnlink struct {
	IDs []int `json:"IDs"`
}

// Validate checks commitments fields are correctly filled
func (c *CommitmentLink) Validate() error {
	if c.Type != "Housing" && c.Type != "Copro" && c.Type != "RenewProject" {
		return errors.New("Type incorrect")
	}
	if c.DestID == 0 {
		return errors.New("ID d'engagement incorrect")
	}
	return nil
}

// Set updates database link between commitments, copros, housings and renew_projects
// according to CommitmentLink datas
func (c *CommitmentLink) Set(db *sql.DB) error {
	var commitmentField string
	switch c.Type {
	case "Housing":
		commitmentField = "housing_id"
	case "Copro":
		commitmentField = "copro_id"
	case "RenewProject":
		commitmentField = "renew_project_id"
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	qry := fmt.Sprintf("UPDATE commitment SET %s=%d WHERE id=ANY($1)",
		commitmentField, c.DestID)
	res, err := tx.Exec(qry, pq.Array(c.IDs))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("rows affected %v", err)
	}
	if int(count) != len(c.IDs) {
		tx.Rollback()
		return errors.New("Impossible de lier tous les engagements")
	}
	return tx.Commit()
}

// Set updates the commitment table to set all copro, renew projects and housings
// link Ids to null
func (c *CommitmentUnlink) Set(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(`UPDATE commitment SET renew_project_id=NULL,
		housing_id=NULL,copro_id=NULL WHERE id=ANY($1)`, pq.Array(c.IDs))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("rows affected %v", err)
	}
	if int(count) != len(c.IDs) {
		tx.Rollback()
		return errors.New("Impossible de supprimer tous les liens")
	}
	return tx.Commit()
}
