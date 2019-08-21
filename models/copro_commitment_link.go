package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// CoproCommitmentLine is used to decode a line of a batch of links between
// commitments and copro operations using the IRIS Code
type CoproCommitmentLine struct {
	Reference string `json:"Reference"`
	IRISCode  string `json:"IRISCode"`
}

// CoproCommitmentBatch embeddes an array of CoproCommitmentLine to
// link commitment to copro operations
type CoproCommitmentBatch struct {
	Lines []CoproCommitmentLine `json:"CoproCommitmentBatch"`
}

// Save takes a batch of housing commitment links and updates the database
func (h *CoproCommitmentBatch) Save(db *sql.DB) error {
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
	if _, err = tx.Exec(`DELETE FROM copro_commitment`); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete query %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("copro_commitment", "reference", "iris_code"))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("statement creation %v", err)
	}
	for _, l := range h.Lines {
		if _, err = stmt.Exec(l.Reference, l.IRISCode); err != nil {
			tx.Rollback()
			return fmt.Errorf("statement execution %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	if err = stmt.Close(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement close %v", err)
	}
	queries := []string{`UPDATE commitment SET copro_id=q.copro_id, 
	housing_id=NULL, renew_project_id=NULL FROM
		(SELECT c.id AS commitment_id, co.id AS copro_id FROM copro co 
			JOIN copro_commitment cc ON co.reference = cc.reference
			JOIN commitment c ON cc.iris_code=c.iris_code) q 
	WHERE commitment.id=q.commitment_id`,
		`DELETE FROM copro_commitment`}
	for i, q := range queries {
		if _, err = tx.Exec(q); err != nil {
			tx.Rollback()
			return fmt.Errorf("query %d %v", i, err)
		}
	}
	tx.Commit()
	return nil
}
