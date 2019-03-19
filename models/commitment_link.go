package models

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/lib/pq"
)

// CommitmentLink is used to embed link or unlink data
type CommitmentLink struct {
	DestID int    `json:"DestID"`
	IDs    []int  `json:"IDs"`
	Link   bool   `json:"Link"`
	Type   string `json:"Type"`
}

// Validate checks commitments fields are correctly filled
func (c *CommitmentLink) Validate() error {
	if c.Type != "Housing" && c.Type != "Copro" && c.Type != "RenewProject" {
		return errors.New("Type incorrect")
	}
	if c.DestID == 0 && c.Link {
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
		commitmentField = "housing_id="
	case "Copro":
		commitmentField = "copro_id="
	case "RenewProject":
		commitmentField = "renew_project_id="
	}
	if c.Link {
		commitmentField += strconv.Itoa(c.DestID)
	} else {
		commitmentField += "null"
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("UPDATE commitment SET "+commitmentField+
		" WHERE id = ANY($1)", pq.Array(c.IDs))
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if int(count) != len(c.IDs) {
		tx.Rollback()
		return errors.New("Impossible de modifier tous les ID d'engagement")
	}
	tx.Commit()
	return nil
}
