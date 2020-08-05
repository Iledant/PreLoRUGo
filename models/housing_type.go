package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// HousingType is used for normalizing housing type
type HousingType struct {
	ID        int64      `json:"ID"`
	ShortName string     `json:"ShortName"`
	LongName  NullString `json:"LongName"`
}

// HousingTypes embeddes an array of housing types for json export
type HousingTypes struct {
	HousingTypes []HousingType `json:"HousingType"`
}

// Create insert a new housing type into database
func (r *HousingType) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO housing_type(short_name,long_name)
	 VALUES($1,$2) RETURNING id`, r.ShortName, r.LongName).Scan(&r.ID)
}

// Valid checks if fields complies with database constraints
func (r *HousingType) Valid() error {
	if r.ShortName == "" {
		return fmt.Errorf("Nom court vide")
	}
	return nil
}

// Get fetches a HousingType from database using ID field
func (r *HousingType) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT short_name,long_name FROM housing_type WHERE ID=$1`,
		r.ID).Scan(&r.ShortName, &r.LongName)
}

// Update modifies a HousingType in database
func (r *HousingType) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE housing_type SET short_name=$1,long_name=$2
	WHERE id=$3`, r.ShortName, r.LongName, r.ID)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	if count != 1 {
		return errors.New("Type introuvable")
	}
	return nil
}

// GetAll fetches all housing transfers from database
func (r *HousingTypes) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT id,short_name,long_name FROM housing_type`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row HousingType
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.ShortName, &row.LongName); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.HousingTypes = append(r.HousingTypes, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.HousingTypes) == 0 {
		r.HousingTypes = []HousingType{}
	}
	return nil
}

// Delete removes housing transfer whose ID is given from database
func (r *HousingType) Delete(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM housing_type WHERE id = $1", r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return errors.New("Type introuvable")
	}
	return nil
}

// IRISHousingType model used to decode one line of batch
type IRISHousingType struct {
	IRISCode             string `json:"IRISCode"`
	HousingTypeShortName string `json:"HousingTypeShortName"`
}

// IRISHousingTypes embeddes an array of IRISHousingType for dedicated query
type IRISHousingTypes struct {
	Lines []IRISHousingType `json:"IRISHousingType"`
}

// Save import a batch of IRISHousingTypes, update the HousingType database and
// update all Housings with the housing types
func (i *IRISHousingTypes) Save(db *sql.DB) error {
	for j, ii := range i.Lines {
		if ii.IRISCode == "" {
			return fmt.Errorf("line %d IrisCode vide", j+1)
		}
		if ii.HousingTypeShortName == "" {
			return fmt.Errorf("ligne %d HousingTypeShortName vide", j+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_iris_housing_type", "iris_code",
		"housing_type_short_name"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, ii := range i.Lines {
		if _, err = stmt.Exec(ii.IRISCode, ii.HousingTypeShortName); err != nil {
			tx.Rollback()
			return fmt.Errorf("stmt exec %v", err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{
		`INSERT INTO housing_type(short_name,long_name)
			SELECT DISTINCT t.housing_type_short_name,NULL FROM temp_iris_housing_type t
			WHERE t.housing_type_short_name NOT IN 
				(SELECT short_name FROM housing_type)`,
		`UPDATE housing SET housing_type_id=q.id
			FROM (SELECT ht.id,hs.housing_ref
			FROM temp_iris_housing_type t
			JOIN housing_type ht ON t.housing_type_short_name=ht.short_name
			JOIN housing_summary hs ON t.iris_code=hs.iris_code) q
			WHERE housing.reference=q.housing_ref`,
	}
	for j, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d: %v", j+1, err)
		}
	}
	return tx.Commit()
}
