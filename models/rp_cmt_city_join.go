package models

import (
	"database/sql"
	"errors"
)

// RPCmtCiyJoin is the model used to save the link between a renew project
// commitment and the city it's attached to
type RPCmtCiyJoin struct {
	ID           int64 `json:"ID"`
	CommitmentID int64 `json:"CommitmentID"`
	CityCode     int64 `json:"CityCode"`
}

// RPCmtCiyJoins embeddes an array of RPCmtCiyJoin for json export
type RPCmtCiyJoins struct {
	RPCmtCiyJoins []RPCmtCiyJoin `json:"RPCmtCiyJoin"`
}

// Validate checks if RPCmtCiyJoin's fields are correctly filled
func (r *RPCmtCiyJoin) Validate() error {
	if r.CommitmentID == 0 || r.CityCode == 0 {
		return errors.New("Champ CommitmentID ou CityCode incorrect")
	}
	return nil
}

// Create insert a new RPCmtCiyJoin into database
func (r *RPCmtCiyJoin) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO rp_cmt_city_join (commitment_id,city_code)
 VALUES($1,$2) RETURNING id`, &r.CommitmentID, &r.CityCode).Scan(&r.ID)
	return err
}

// Get fetches a RPCmtCiyJoin from database using ID field
func (r *RPCmtCiyJoin) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT commitment_id,city_code 
	FROM rp_cmt_city_join WHERE ID=$1`, r.ID).Scan(&r.CommitmentID, &r.CityCode)
	return err
}

// Update modifies a RPCmtCiyJoin in database
func (r *RPCmtCiyJoin) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE rp_cmt_city_join SET commitment_id=$1,city_code=$2 
	WHERE id=$3`, r.CommitmentID, r.CityCode, r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Lien introuvable")
	}
	return err
}

// GetAll fetches all Communities from database
func (r *RPCmtCiyJoins) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,commitment_id,city_code FROM rp_cmt_city_join`)
	if err != nil {
		return err
	}
	var row RPCmtCiyJoin
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CityCode); err != nil {
			return err
		}
		r.RPCmtCiyJoins = append(r.RPCmtCiyJoins, row)
	}
	err = rows.Err()
	if len(r.RPCmtCiyJoins) == 0 {
		r.RPCmtCiyJoins = []RPCmtCiyJoin{}
	}
	return err
}

// Delete removes community whose ID is given from database
func (r *RPCmtCiyJoin) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM rp_cmt_city_join WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Lien introuvable")
	}
	return nil
}
