package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// HousingCommitmentLine is used to decode a line of a batch of links between
// commitments and housing operations using the IRIS Code
type HousingCommitmentLine struct {
	Reference string `json:"Reference"`
	IRISCode  string `json:"IRISCode"`
}

// HousingCommitmentBach embeddes an array of HousingCommitmentLine to
// link commitment to housing operations
type HousingCommitmentBach struct {
	Lines []HousingCommitmentLine `json:"HousingCommitmentBach"`
}

// Save takes a batch of housing commitment links and updates the database
func (h *HousingCommitmentBach) Save(db *sql.DB) error {
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
	if _, err = tx.Exec(`DELETE FROM housing_commitment`); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete query %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("housing_commitment", "reference", "iris_code"))
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
	queries := []string{`UPDATE commitment SET housing_id=q.housing_id, 
	copro_id=NULL, renew_project_id=NULL FROM
		(SELECT c.id AS commitment_id,h.id AS housing_id FROM housing h 
			JOIN housing_commitment hc ON h.reference=hc.reference
			JOIN commitment c ON hc.iris_code=c.iris_code) q 
	WHERE commitment.id=q.commitment_id`,
		`DELETE FROM housing_commitment`}
	for i, q := range queries {
		if _, err = tx.Exec(q); err != nil {
			tx.Rollback()
			return fmt.Errorf("query %d %v", i, err)
		}
	}
	return tx.Commit()
}
