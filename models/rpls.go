package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// RPLS is the model to store
type RPLS struct {
	ID        int64   `json:"ID"`
	InseeCode int64   `json:"InseeCode"`
	CityName  string  `json:"CityName"`
	Year      int64   `json:"Year"`
	Ratio     float64 `json:"Ratio"`
}

// RPLSArray embeddes an array of RPLS for json export
type RPLSArray struct {
	Lines []RPLS `json:"RPLS"`
}

// RPLSLine is used t decode one line of RPLS batch
type RPLSLine struct {
	InseeCode int64   `json:"InseeCode"`
	Year      int64   `json:"Year"`
	Ratio     float64 `json:"Ratio"`
}

// RPLSBatch embeddes a batch of RPSLine for import
type RPLSBatch struct {
	Lines []RPLSLine `json:"RPLS"`
}

// RPLSYears embeddes the distinct years in the rpls table
type RPLSYears struct {
	Lines []int64 `json:"RPLSYear"`
}

// Validate checks of fields can be inserted into database
func (r *RPLS) Validate() error {
	if r.InseeCode == 0 {
		return fmt.Errorf("InseeCode nul")
	}
	if r.Year == 0 {
		return fmt.Errorf("Year nul")
	}
	return nil
}

// Create insert a new RPLS into database
func (r *RPLS) Create(db *sql.DB) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if err := db.QueryRow(`INSERT INTO rpls (insee_code,year,ratio) 
	VALUES ($1,$2,$3) RETURNING id`, r.InseeCode, r.Year, r.Ratio).Scan(&r.ID); err != nil {
		return fmt.Errorf("insert %v", err)
	}
	return db.QueryRow(`SELECT name FROM city WHERE insee_code=$1`, r.InseeCode).
		Scan(&r.CityName)
}

// Update modifies a RPLS in database
func (r *RPLS) Update(db *sql.DB) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.ID == 0 {
		return fmt.Errorf("ID nul")
	}
	res, err := db.Exec(`UPDATE rpls SET insee_code=$1, year=$2, ratio=$3
	WHERE id=$4`, r.InseeCode, r.Year, r.Ratio, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return fmt.Errorf("RPLS introuvable")
	}
	return db.QueryRow(`SELECT name FROM city WHERE insee_code=$1`, r.InseeCode).
		Scan(&r.CityName)
}

// Delete removes a RPLS from database
func (r *RPLS) Delete(db *sql.DB) error {
	res, err := db.Exec(`DELETE FROM rpls WHERE id=$1`, r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count != 1 {
		return fmt.Errorf("RPLS introuvable")
	}
	return nil
}

// Save insert a batch of rpls into database
func (r *RPLSBatch) Save(db *sql.DB) error {
	for i, l := range r.Lines {
		if l.InseeCode == 0 {
			return fmt.Errorf("ligne %d InseeCode null", i+1)
		}
		if l.Year == 0 {
			return fmt.Errorf("ligne %d Year null", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_rpls", "insee_code", "year", "ratio"))
	if err != nil {
		return fmt.Errorf("copy in %v", err)
	}
	defer stmt.Close()
	for i, l := range r.Lines {
		if _, err = stmt.Exec(l.InseeCode, l.Year, l.Ratio); err != nil {
			tx.Rollback()
			return fmt.Errorf("ligne %d %v", i+1, err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	queries := []string{`UPDATE rpls SET ratio=t.ratio FROM temp_rpls t 
	WHERE t.insee_code=rpls.insee_code AND t.year=rpls.year`,
		`INSERT INTO rpls (insee_code,year,ratio)
	SELECT insee_code,year,ratio from temp_rpls 
		WHERE (insee_code,year) NOT IN (SELECT insee_code,year from rpls)`,
		`DELETE from temp_rpls`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %v", i, err)
		}
	}
	return tx.Commit()
}

// GetAll fetches all RPLS from database and add CityName
func (r *RPLSArray) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT r.id,r.insee_code,c.name,r.year,r.ratio FROM rpls r
		JOIN city c ON r.insee_code=c.insee_code`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var l RPLS
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&l.ID, &l.InseeCode, &l.CityName, &l.Year,
			&l.Ratio); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.Lines = append(r.Lines, l)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.Lines) == 0 {
		r.Lines = []RPLS{}
	}
	return nil
}

// GetAll fetches all years from the rpls table
func (r *RPLSYears) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT DISTINCT year FROM rpls`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var y int64
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&y); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.Lines = append(r.Lines, y)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.Lines) == 0 {
		r.Lines = []int64{}
	}
	return nil
}
