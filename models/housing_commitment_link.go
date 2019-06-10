package models

import (
	"database/sql"
	"fmt"
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
	stmt, err := tx.Prepare(`UPDATE commitment SET housing_id=housing.id, 
	copro_id=NULL, renew_project_id=NULL FROM housing 
	WHERE housing.reference = $1 AND commitment.iris_code=$2`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("statement creation %v", err)
	}
	for i, l := range h.Lines {
		_, err := stmt.Exec(l.Reference, l.IRISCode)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("statement exec %d %v", i, err)
		}
	}
	tx.Commit()
	return nil
}
