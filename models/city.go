package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/lib/pq"
)

// City model
type City struct {
	InseeCode     int64      `json:"InseeCode"`
	Name          string     `json:"Name"`
	CommunityID   NullInt64  `json:"CommunityID"`
	CommunityName NullString `json:"CommunityName"`
	QPV           bool       `json:"QPV"`
}

// Cities embeddes an array of City for json export
type Cities struct {
	Cities []City `json:"City"`
}

// CityLine is used to decode a line of City batch
type CityLine struct {
	InseeCode     int64      `json:"InseeCode"`
	Name          string     `json:"Name"`
	CommunityCode NullString `json:"CommunityCode"`
	QPV           bool       `json:"QPV"`
}

// CityBatch embeddes an array of CityLine for json export
type CityBatch struct {
	Lines []CityLine `json:"City"`
}

// PaginatedCity is used to fetch a city with community full name
type PaginatedCity struct {
	InseeCode     int64      `json:"InseeCode"`
	Name          string     `json:"Name"`
	QPV           bool       `json:"QPV"`
	CommunityID   NullInt64  `json:"CommunityID"`
	CommunityName NullString `json:"CommunityName"`
}

// PaginatedCities embeddes an array of City for paginated display
type PaginatedCities struct {
	Cities     []PaginatedCity `json:"City"`
	Page       int64           `json:"Page"`
	ItemsCount int64           `json:"ItemsCount"`
}

// Validate checks if City's fields are correctly filled
func (c *City) Validate() error {
	if c.InseeCode == 0 || c.Name == "" {
		return errors.New("Champ incorrect")
	}
	return nil
}

// Create insert a new City into database
func (c *City) Create(db *sql.DB) (err error) {
	res, err := db.Exec(`INSERT INTO city (insee_code,name,community_id,qpv)
 VALUES($1,$2,$3,$4)`, &c.InseeCode, &c.Name, &c.CommunityID, &c.QPV)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Insertion de la ville non réussie")
	}
	if c.CommunityID.Valid {
		err = db.QueryRow(`SELECT name from community WHERE id=$1`, c.CommunityID).
			Scan(&c.CommunityName)
	}
	return err
}

// Get fetches a City from database using ID field
func (c *City) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT c.name, c.community_id, c.qpv, o.name FROM city c
		LEFT JOIN community o ON c.community_id = o.id	
	WHERE insee_code=$1`,
		c.InseeCode).Scan(&c.Name, &c.CommunityID, &c.QPV, &c.CommunityName)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a city in database
func (c *City) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE city SET name=$1,community_id=$2,qpv=$3 
		WHERE insee_code=$4`, c.Name, c.CommunityID, c.QPV, c.InseeCode)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Ville introuvable")
	}
	if c.CommunityID.Valid {
		err = db.QueryRow(`SELECT name FROM community WHERE id=$1`, c.CommunityID).
			Scan(&c.CommunityName)
	}
	return err
}

// GetAll fetches all Cities from database
func (c *Cities) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT c.insee_code,c.name,c.community_id,c.qpv,o.name
		FROM city c
		LEFT JOIN community o ON c.community_id=o.id`)
	if err != nil {
		return err
	}
	var row City
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.InseeCode, &row.Name, &row.CommunityID, &row.QPV,
			&row.CommunityName); err != nil {
			return err
		}
		c.Cities = append(c.Cities, row)
	}
	err = rows.Err()
	if len(c.Cities) == 0 {
		c.Cities = []City{}
	}
	return err
}

// Delete removes city whose ID is given from database
func (c *City) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM city WHERE insee_code = $1", c.InseeCode)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Ville introuvable")
	}
	return nil
}

// Save insert a batch of CityLine into database
func (c *CityBatch) Save(db *sql.DB) (err error) {
	for i, r := range c.Lines {
		if r.InseeCode == 0 || r.Name == "" {
			return fmt.Errorf("Ligne %d, champ incorrect", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_city", "insee_code", "name",
		"community_code", "qpv"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range c.Lines {
		if _, err = stmt.Exec(r.InseeCode, r.Name, r.CommunityCode, r.QPV); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	queries := []string{`UPDATE city SET name=q.name,community_id=q.id,qpv=q.qpv 
	FROM (SELECT t.*, c.id FROM temp_city t 
					LEFT JOIN community c ON t.community_code = c.code) q 
	WHERE q.insee_code = city.insee_code`,
		`INSERT INTO city (insee_code,name,community_id,qpv)
	SELECT t.insee_code,t.name,c.id,t.qpv from temp_city t 
		LEFT JOIN community c ON t.community_code = c.code
	WHERE insee_code NOT IN (SELECT DISTINCT insee_code from city)`,
		`DELETE FROM temp_city`,
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	return tx.Commit()
}

// Get fetches all cities that matches the search pattern
func (p *PaginatedCities) Get(db *sql.DB, q *PaginatedQuery) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM city c
	LEFT JOIN community o on o.id = c.community_id
		WHERE c.name ILIKE $1 OR c.insee_code::varchar ILIKE $1`, "%"+q.Search+"%").
		Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT c.insee_code,c.name, o.id, o.name,c.qpv FROM city c
	LEFT JOIN community o on o.id = c.community_id
	WHERE c.name ILIKE $1 OR c.insee_code::varchar ILIKE $1
	ORDER BY 1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $2`,
		"%"+q.Search+"%", offset)
	if err != nil {
		return err
	}
	var row PaginatedCity
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.InseeCode, &row.Name, &row.CommunityID,
			&row.CommunityName, &row.QPV); err != nil {
			return err
		}
		p.Cities = append(p.Cities, row)
	}
	err = rows.Err()
	if len(p.Cities) == 0 {
		p.Cities = []PaginatedCity{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}
