package models

import (
	"database/sql"
	"errors"
)

// Community model
type Community struct {
	ID   int64  `json:"ID"`
	Code string `json:"Code"`
	Name string `json:"Name"`
}

// Communities embeddes an array of Community for json export
type Communities struct {
	Communities []Community `json:"Community"`
}

// CommunityLine is used to decode a line of Community batch
type CommunityLine struct {
	Code string `json:"Code"`
	Name string `json:"Name"`
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
	err = db.QueryRow(`INSERT INTO community (code,name)
 VALUES($1,$2) RETURNING id`, &c.Code, &c.Name).Scan(&c.ID)
	return err
}

// Get fetches a Community from database using ID field
func (c *Community) Get(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT code, name FROM community WHERE ID=$1`, c.ID).
		Scan(&c.Code, &c.Name)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies a community in database
func (c *Community) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE community SET code=$1,name=$2 WHERE id=$3`,
		c.Code, c.Name, c.ID)
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
	rows, err := db.Query(`SELECT id,code,name FROM community`)
	if err != nil {
		return err
	}
	var row Community
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Code, &row.Name); err != nil {
			return err
		}
		c.Communities = append(c.Communities, row)
	}
	err = rows.Err()
	return err
}

// Delete removes community whose ID is given from database
func (c *Community) Delete(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// TODO : remove comment to update city table
	// if _, err := tx.Exec("UPDATE city SET community_id=NULL where community_id = $1",
	// 	c.ID); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
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
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO temp_community (code,name) VALUES ($1,$2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, r := range c.Lines {
		if r.Code == "" || r.Name == "" {
			tx.Rollback()
			return errors.New("Champs incorrects")
		}
		if _, err = stmt.Exec(r.Code, r.Name); err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`UPDATE community SET name=t.name FROM temp_community t 
	WHERE t.code = community.code`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO community (code,name)
	SELECT code,name from temp_community 
	  WHERE code NOT IN (SELECT DISTINCT code from community)`)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
