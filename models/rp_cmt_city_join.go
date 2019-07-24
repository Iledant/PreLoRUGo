package models

import (
	"database/sql"
	"errors"
)

// RPCmtCityJoin is the model used to save the link between a renew project
// commitment and the city it's attached to
type RPCmtCityJoin struct {
	ID           int64 `json:"ID"`
	CommitmentID int64 `json:"CommitmentID"`
	CityCode     int64 `json:"CityCode"`
}

// RPCmtCityJoins embeddes an array of RPCmtCityJoin for json export
type RPCmtCityJoins struct {
	RPCmtCityJoins []RPCmtCityJoin `json:"RPCmtCityJoin"`
}

// Validate checks if RPCmtCityJoin's fields are correctly filled
func (r *RPCmtCityJoin) Validate() error {
	if r.CommitmentID == 0 || r.CityCode == 0 {
		return errors.New("Champ CommitmentID ou CityCode incorrect")
	}
	return nil
}

// Create insert a new RPCmtCityJoin into database
func (r *RPCmtCityJoin) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO rp_cmt_city_join (commitment_id,city_code)
 VALUES($1,$2) RETURNING id`, &r.CommitmentID, &r.CityCode).Scan(&r.ID)
	return err
}

// Get fetches a RPCmtCityJoin from database using ID field
func (r *RPCmtCityJoin) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT commitment_id,city_code 
	FROM rp_cmt_city_join WHERE ID=$1`, r.ID).Scan(&r.CommitmentID, &r.CityCode)
	return err
}

// Update modifies a RPCmtCityJoin in database
func (r *RPCmtCityJoin) Update(db *sql.DB) (err error) {
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

// GetAll fetches all RPCmtCityJoin from database
func (r *RPCmtCityJoins) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,commitment_id,city_code FROM rp_cmt_city_join`)
	if err != nil {
		return err
	}
	var row RPCmtCityJoin
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CityCode); err != nil {
			return err
		}
		r.RPCmtCityJoins = append(r.RPCmtCityJoins, row)
	}
	err = rows.Err()
	if len(r.RPCmtCityJoins) == 0 {
		r.RPCmtCityJoins = []RPCmtCityJoin{}
	}
	return err
}

// Delete removes RPCmtCityJoin whose ID is given from database
func (r *RPCmtCityJoin) Delete(db *sql.DB) (err error) {
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

// GetLinked fetches all commitment city joins attached to a renew project
func (r *RPCmtCityJoins) GetLinked(db *sql.DB, ID int64) (err error) {
	rows, err := db.Query(`SELECT r.id,r.commitment_id,r.city_code FROM rp_cmt_city_join r
	JOIN commitment c ON r.commitment_id=c.id
	WHERE c.renew_project_id=$1`, ID)
	if err != nil {
		return err
	}
	var row RPCmtCityJoin
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CommitmentID, &row.CityCode); err != nil {
			return err
		}
		r.RPCmtCityJoins = append(r.RPCmtCityJoins, row)
	}
	err = rows.Err()
	if len(r.RPCmtCityJoins) == 0 {
		r.RPCmtCityJoins = []RPCmtCityJoin{}
	}
	return err
}
