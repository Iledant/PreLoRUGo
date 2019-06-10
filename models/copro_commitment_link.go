package models

import (
	"database/sql"
	"fmt"
)

// CoproCommitmentLine is used to decode a line of a batch of links between
// commitments and copro operations using the IRIS Code
type CoproCommitmentLine struct {
	Reference string `json:"Reference"`
	IRISCode  string `json:"IRISCode"`
}

// CoproCommitmentBach embeddes an array of CoproCommitmentLine to
// link commitment to copro operations
type CoproCommitmentBach struct {
	Lines []CoproCommitmentLine `json:"CoproCommitmentBach"`
}

// Save takes a batch of housing commitment links and updates the database
func (h *CoproCommitmentBach) Save(db *sql.DB) error {
	for i, l := range h.Lines {
		if l.Reference == "" {
			return fmt.Errorf("Ligne %d, Reference vide", i+1)
		}
		if l.IRISCode == "" {
			return fmt.Errorf("Ligne %d, IRISCode vide", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction %v", err)
	}
	stmt, err := tx.Prepare(`UPDATE commitment SET copro_id=copro.id, 
	housing_id=NULL, renew_project_id=NULL FROM copro 
	WHERE copro.reference = $1 AND commitment.iris_code=$2`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("statement creation %v", err)
	}
	for i, l := range h.Lines {
		res, err := stmt.Exec(l.Reference, l.IRISCode)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("statement exec %d %v", i, err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("statement %d count %v", i, err)
		}
		if count == 0 {
			tx.Rollback()
			return fmt.Errorf("ligne %d Reference ou code IRIS introuvable", i)
		}
	}
	tx.Commit()
	return nil
}
