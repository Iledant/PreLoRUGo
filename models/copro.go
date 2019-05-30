package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Copro is the model for a labelled condominium
type Copro struct {
	ID        int64     `json:"ID"`
	Reference string    `json:"Reference"`
	Name      string    `json:"Name"`
	Address   string    `json:"Address"`
	ZipCode   int       `json:"ZipCode"`
	CityName  string    `json:"CityName"`
	LabelDate NullTime  `json:"LabelDate"`
	Budget    NullInt64 `json:"Budget"`
}

// Copros embeddes an array of Copro for JSON export
type Copros struct {
	Copros []Copro `json:"Copro"`
}

// CoproLine is used to decode a line of copro batch
type CoproLine struct {
	Reference string    `json:"Reference"`
	Name      string    `json:"Name"`
	Address   string    `json:"Address"`
	ZipCode   int       `json:"ZipCode"`
	LabelDate NullInt64 `json:"LabelDate"`
	Budget    NullInt64 `json:"Budget"`
}

// CoproBatch embeddes an array of CoproLine for batch import
type CoproBatch struct {
	Lines []CoproLine `json:"Copro"`
}

// Validate checks copro's fields and return an error if they don't
// fit with database constraints
func (c *Copro) Validate() error {
	if c.Reference == "" || c.Name == "" || c.Address == "" || c.ZipCode == 0 {
		return errors.New("Champ reference, name, address ou zipcode vide")
	}
	return nil
}

// Create a new copro entry according to fields and returning id
func (c *Copro) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO copro (reference,name,address,
		zip_code,label_date,budget) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		c.Reference, c.Name, c.Address, c.ZipCode, c.LabelDate, c.Budget).Scan(&c.ID)
	if err != nil {
		return err
	}
	err = db.QueryRow(`SELECT city.name FROM copro 
	JOIN city ON copro.zip_code=city.insee_code 
	WHERE copro.id=$1`, c.ID).Scan(&c.CityName)
	return err
}

// Update modifies the copro fields whose ID is given
func (c *Copro) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE copro SET reference=$1, name=$2,
	address=$3, zip_code=$4, label_date=$5, budget=$6 WHERE id = $7`,
		c.Reference, c.Name, c.Address, c.ZipCode, c.LabelDate, c.Budget, c.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Copro introuvable")
	}
	err = db.QueryRow(`SELECT city.name FROM copro 
	JOIN city ON copro.zip_code=city.insee_code 
	WHERE id=$1`, c.ID).Scan(&c.CityName)
	return err
}

// Delete removes the copro whose ID is given
func (c *Copro) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM copro WHERE id = $1", c.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Copro introuvable")
	}
	return nil
}

// Get fetches the copro whose ID is given
func (c *Copro) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT co.id,co.reference,co.name,co.address,co.zip_code,
	ci.name,co.label_date,co.budget FROM copro co
	JOIN city ci ON co.zip_code=ci.insee_code
	WHERE id = $1`, c.ID).
		Scan(&c.ID, &c.Reference, &c.Name, &c.Address, &c.ZipCode, &c.CityName,
			&c.LabelDate, &c.Budget)
	return err
}

// GetAll fetches all Copros from database
func (c *Copros) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT co.id,co.reference,co.name,co.address,co.zip_code,
	ci.name,co.label_date,co.budget FROM copro co
	JOIN city ci ON co.zip_code=ci.insee_code`)
	if err != nil {
		return err
	}
	var r Copro
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ID, &r.Reference, &r.Name, &r.Address, &r.ZipCode,
			&r.CityName, &r.LabelDate, &r.Budget); err != nil {
			return err
		}
		c.Copros = append(c.Copros, r)
	}
	err = rows.Err()
	if len(c.Copros) == 0 {
		c.Copros = []Copro{}
	}
	return err
}

// nullExcel2NullTime convert a null int64 corresponding to an Excel integer
// date to a NullTime struct
func nullExcel2NullTime(i NullInt64) NullTime {
	var n NullTime
	if !i.Valid {
		return n
	}
	n.Valid = true
	n.Time = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC).
		Add(time.Duration(i.Int64*24) * time.Hour)
	return n
}

// Save insert a batch of CoproLine into database
func (c *CoproBatch) Save(db *sql.DB) (err error) {
	for i, r := range c.Lines {
		if r.Reference == "" || r.Name == "" || r.Address == "" || r.ZipCode == 0 {
			return fmt.Errorf("ligne %d : champs incorrects", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_copro 
	(reference,name,address,zip_code,label_date,budget) 
	VALUES ($1,$2,$3,$4,$5,$6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range c.Lines {
		if _, err = stmt.Exec(r.Reference, r.Name, r.Address, r.ZipCode,
			nullExcel2NullTime(r.LabelDate), r.Budget); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	queries := []string{`UPDATE copro SET name=t.name, address=t.address, zip_code=t.zip_code,
	label_date=t.label_date,budget=t.budget FROM temp_copro t WHERE t.reference = copro.reference`,
		`INSERT INTO copro (reference, name,address,zip_code,label_date,budget)
	SELECT reference,name,address,zip_code,label_date,budget from temp_copro 
		WHERE reference NOT IN (SELECT reference from copro)`,
		`DELETE from temp_copro`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	tx.Commit()
	return nil
}
