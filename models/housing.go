package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/lib/pq"
)

// Housing model
type Housing struct {
	ID                   int64      `json:"ID"`
	Reference            string     `json:"Reference"`
	Address              NullString `json:"Address"`
	ZipCode              NullInt64  `json:"ZipCode"`
	CityName             NullString `json:"CityName"`
	PLAI                 int64      `json:"PLAI"`
	PLUS                 int64      `json:"PLUS"`
	PLS                  int64      `json:"PLS"`
	ANRU                 bool       `json:"ANRU"`
	HousingTypeID        NullInt64  `json:"HousingTypeID"`
	HousingTypeShortName NullString `json:"HousingTypeShortName"`
	HousingTypeLongName  NullString `json:"HousingTypeLongName"`
}

// Housings embeddes an array of Housing for json export
type Housings struct {
	Housings []Housing `json:"Housing"`
}

// HousingLine is used to decode a line of Housing batch
type HousingLine struct {
	Reference string     `json:"Reference"`
	Address   NullString `json:"Address"`
	ZipCode   NullInt64  `json:"ZipCode"`
	PLAI      int64      `json:"PLAI"`
	PLUS      int64      `json:"PLUS"`
	PLS       int64      `json:"PLS"`
	ANRU      bool       `json:"ANRU"`
}

// HousingBatch embeddes an array of HousingLine for json export
type HousingBatch struct {
	Lines []HousingLine `json:"Housing"`
}

// PaginatedHousings embeddes an array of housing for paginated get request
type PaginatedHousings struct {
	Housings   []Housing `json:"Housing"`
	Page       int64     `json:"Page"`
	ItemsCount int64     `json:"ItemsCount"`
}

// Validate checks if Housing's fields are correctly filled
func (h *Housing) Validate() error {
	if h.Reference == "" {
		return errors.New("Champ Reference incorrect")
	}
	return nil
}

// Get fetches a housing from database using the ID field
func (h *Housing) Get(db *sql.DB) error {
	return db.QueryRow(`SELECT h.id,h.reference,h.address,h.zip_code,c.name,
	h.plai,h.plus,h.pls,h.anru,h.housing_type_id,ht.short_name,ht.long_name
	FROM housing h
	LEFT JOIN city c ON h.zip_code=c.insee_code 
	LEFT JOIN housing_type ht ON h.housing_type_id=ht.id
	WHERE h.id=$1`, h.ID).Scan(&h.ID, &h.Reference, &h.Address, &h.ZipCode,
		&h.CityName, &h.PLAI, &h.PLUS, &h.PLS, &h.ANRU, &h.HousingTypeID,
		&h.HousingTypeShortName, &h.HousingTypeLongName)
}

// Create insert a new Housing into database
func (h *Housing) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO housing 
	(reference,address,zip_code,plai,plus,pls,anru,housing_type_id)
	 VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`, &h.Reference, &h.Address,
		&h.ZipCode, &h.PLAI, &h.PLUS, &h.PLS, &h.ANRU, &h.HousingTypeID).Scan(&h.ID)
	if err != nil {
		return err
	}
	err = db.QueryRow(`SELECT name FROM city WHERE insee_code=$1`, h.ZipCode).
		Scan(&h.CityName)
	if err != nil {
		return fmt.Errorf("city scan %v", err)
	}
	err = db.QueryRow(`SELECT short_name,long_name FROM housing_type
		WHERE id=$1`, h.HousingTypeID).
		Scan(&h.CityName)
	if err == sql.ErrNoRows {
		h.HousingTypeShortName.Valid = false
		h.HousingTypeLongName.Valid = false
		err = nil
	}
	return err
}

// Update modifies a housing in database
func (h *Housing) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE housing SET reference=$1,address=$2,zip_code=$3,
	plai=$4,plus=$5,pls=$6,anru=$7,housing_type_id=$8 WHERE id=$9`, h.Reference,
		h.Address, h.ZipCode, h.PLAI, h.PLUS, h.PLS, h.ANRU, h.HousingTypeID, h.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Logement introuvable")
	}
	err = db.QueryRow(`SELECT name FROM city WHERE insee_code=$1`, h.ZipCode).
		Scan(&h.CityName)
	if err == sql.ErrNoRows {
		h.CityName.Valid = false
		err = nil
	}
	err = db.QueryRow(`SELECT short_name,long_name FROM housing_type
		WHERE id=$1`, h.HousingTypeID).
		Scan(&h.HousingTypeShortName, &h.HousingTypeLongName)
	if err == sql.ErrNoRows {
		h.HousingTypeShortName.Valid = false
		h.HousingTypeLongName.Valid = false
		err = nil
	}
	return err
}

// GetAll fetches all Housings from database
func (h *Housings) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT h.id,h.reference,h.address,h.zip_code,c.name,
	h.plai,h.plus,h.pls,h.anru,h.housing_type_id,ht.short_name,ht.long_name 
	FROM housing h
	LEFT JOIN city c ON h.zip_code=c.insee_code
	LEFT JOIN housing_type ht ON h.housing_type_id=ht.id`)
	if err != nil {
		return err
	}
	var row Housing
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Reference, &row.Address, &row.ZipCode,
			&row.CityName, &row.PLAI, &row.PLUS, &row.PLS, &row.ANRU,
			&row.HousingTypeID, &row.HousingTypeShortName, &row.HousingTypeLongName); err != nil {
			return err
		}
		h.Housings = append(h.Housings, row)
	}
	err = rows.Err()
	if len(h.Housings) == 0 {
		h.Housings = []Housing{}
	}
	return err
}

// Delete removes housing whose ID is given from database
func (h *Housing) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM housing WHERE id = $1", h.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Logement introuvable")
	}
	return nil
}

// Save insert a batch of HousingLine into database
func (h *HousingBatch) Save(db *sql.DB) (err error) {
	for i, r := range h.Lines {
		if r.Reference == "" {
			return fmt.Errorf("ligne %d, champ Reference incorrect", i)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_housing", "reference", "address",
		"zip_code", "plai", "plus", "pls", "anru"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range h.Lines {
		if _, err = stmt.Exec(r.Reference, r.Address, r.ZipCode, r.PLAI, r.PLUS,
			r.PLS, r.ANRU); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	queries := []string{`UPDATE housing SET address=t.address,zip_code=t.zip_code,
	plai=t.plai,plus=t.plus,pls=t.pls,anru=t.anru FROM temp_housing t 
	WHERE t.reference = housing.reference`,
		`INSERT INTO housing
	(reference,address,zip_code,plai,plus,pls,anru)
	SELECT reference,address,zip_code,plai,plus,pls,anru from temp_housing 
		WHERE reference NOT IN (SELECT reference from housing)`,
		`DELETE from temp_housing`,
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

// Get fetches a bath of paginated housings form database that fetch a search
// pattern
func (p *PaginatedHousings) Get(db *sql.DB, q *PaginatedQuery) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM housing
		WHERE reference ILIKE $1 OR address ILIKE $1 OR zip_code::varchar ILIKE $1`,
		"%"+q.Search+"%").Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT h.id,h.reference,h.address,h.zip_code,c.name,
	h.plai,h.plus,h.pls,h.anru FROM housing h
	LEFT JOIN city c ON h.zip_code=c.insee_code
	WHERE reference ILIKE $1 OR address ILIKE $1 OR zip_code::varchar ILIKE $1
	ORDER BY 1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $2`,
		"%"+q.Search+"%", offset)
	if err != nil {
		return err
	}
	var row Housing
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Reference, &row.Address, &row.ZipCode,
			&row.CityName, &row.PLAI, &row.PLUS, &row.PLS, &row.ANRU); err != nil {
			return err
		}
		p.Housings = append(p.Housings, row)
	}
	err = rows.Err()
	if len(p.Housings) == 0 {
		p.Housings = []Housing{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}
