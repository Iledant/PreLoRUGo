package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// Community model
type Community struct {
	ID           int64     `json:"ID"`
	Code         string    `json:"Code"`
	Name         string    `json:"Name"`
	DepartmentID NullInt64 `json:"DepartmentID"`
}

// Communities embeddes an array of Community for json export
type Communities struct {
	Communities []Community `json:"Community"`
}

// CommunityLine is used to decode a line of Community batch
type CommunityLine struct {
	Code           string `json:"Code"`
	Name           string `json:"Name"`
	DepartmentCode int    `json:"DepartmentCode"`
}

// CommunityBatch embeddes an array of CommunityLine for json export
type CommunityBatch struct {
	Lines []CommunityLine `json:"Community"`
}

// Validate checks if Community's fields are correctly filled
func (c *Community) Validate() error {
	if c.Code == "" || c.Name == "" {
		return errors.New("Champ code ou name incorrect")
	}
	return nil
}

// Create insert a new Community into database
func (c *Community) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO community (code,name,department_id)
 VALUES($1,$2,$3) RETURNING id`, &c.Code, &c.Name, &c.DepartmentID).Scan(&c.ID)
	return err
}

// Get fetches a Community from database using ID field
func (c *Community) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT code, name,department_id FROM community WHERE ID=$1`,
		c.ID).Scan(&c.Code, &c.Name, &c.DepartmentID)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a community in database
func (c *Community) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE community SET code=$1,name=$2,department_id=$3 
	WHERE id=$4`, c.Code, c.Name, c.DepartmentID, c.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Interco introuvable")
	}
	return err
}

// GetAll fetches all Communities from database
func (c *Communities) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id,code,name,department_id FROM community`)
	if err != nil {
		return err
	}
	var row Community
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Code, &row.Name, &row.DepartmentID); err != nil {
			return err
		}
		c.Communities = append(c.Communities, row)
	}
	err = rows.Err()
	if len(c.Communities) == 0 {
		c.Communities = []Community{}
	}
	return err
}

// Delete removes community whose ID is given from database
func (c *Community) Delete(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE city SET community_id=NULL where community_id = $1",
		c.ID); err != nil {
		tx.Rollback()
		return err
	}
	res, err := tx.Exec("DELETE FROM community WHERE id = $1", c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if count != 1 {
		tx.Rollback()
		return errors.New("Interco introuvable")
	}
	tx.Commit()
	return nil
}

// Save insert a batch of CommunityLine into database
func (c *CommunityBatch) Save(db *sql.DB) (err error) {
	for i, r := range c.Lines {
		if r.Code == "" || r.Name == "" || r.DepartmentCode == 0 {
			return fmt.Errorf("ligne %d, champ incorrect", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_community", "code", "name", "department_code"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range c.Lines {
		if _, err = stmt.Exec(r.Code, r.Name, r.DepartmentCode); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{`UPDATE community SET name=t.name, department_id=d.id 
	FROM temp_community t, department d 
	WHERE t.code = community.code AND d.code=t.department_code`,
		`INSERT INTO community (code,name,department_id)
	SELECT t.code,t.name,d.id FROM temp_community t 
	LEFT OUTER JOIN department d ON d.code=t.department_code
		WHERE t.code NOT IN (SELECT DISTINCT code from community)`,
		`DELETE FROM temp_community`,
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
