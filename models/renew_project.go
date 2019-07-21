package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// RenewProject model
type RenewProject struct {
	ID             int64      `json:"ID"`
	Reference      string     `json:"Reference"`
	Name           string     `json:"Name"`
	Budget         int64      `json:"Budget"`
	PRIN           bool       `json:"PRIN"`
	CityCode1      int64      `json:"CityCode1"`
	CityName1      string     `json:"CityName1"`
	BudgetCity1    NullInt64  `json:"BudgetCity1"`
	CityCode2      NullInt64  `json:"CityCode2"`
	CityName2      NullString `json:"CityName2"`
	BudgetCity2    NullInt64  `json:"BudgetCity2"`
	CityCode3      NullInt64  `json:"CityCode3"`
	CityName3      NullString `json:"CityName3"`
	BudgetCity3    NullInt64  `json:"BudgetCity3"`
	Population     NullInt64  `json:"Population"`
	CompositeIndex NullInt64  `json:"CompositeIndex"`
}

// RenewProjects embeddes an array of RenewProject for json export
type RenewProjects struct {
	RenewProjects []RenewProject `json:"RenewProject"`
}

// RenewProjectLine is used to decode one line of renew projects batch
type RenewProjectLine struct {
	Reference      string    `json:"Reference"`
	Name           string    `json:"Name"`
	Budget         int64     `json:"Budget"`
	PRIN           bool      `json:"PRIN"`
	CityCode1      int64     `json:"CityCode1"`
	BudgetCity1    NullInt64 `json:"BudgetCity1"`
	CityCode2      NullInt64 `json:"CityCode2"`
	BudgetCity2    NullInt64 `json:"BudgetCity2"`
	CityCode3      NullInt64 `json:"CityCode3"`
	Population     NullInt64 `json:"Population"`
	BudgetCity3    NullInt64 `json:"BudgetCity3"`
	CompositeIndex NullInt64 `json:"CompositeIndex"`
}

// RenewProjectBatch embeddes an array of RenewProjectLine
type RenewProjectBatch struct {
	Lines []RenewProjectLine `json:"RenewProject"`
}

// Validate checks if the fields of a renew project are correctly filled
func (r *RenewProject) Validate() error {
	if r.Reference == "" || r.Name == "" || r.Budget == 0 || r.CityCode1 == 0 {
		return errors.New("Champ reference, name ou budget incorrect")
	}
	return nil
}

// GetByID fetches all fields from a renew project whose ID is given
func (r *RenewProject) GetByID(db *sql.DB) error {
	return db.QueryRow(`SELECT r.reference,r.name, r.budget,r.prin,r.city_code1, c1.name,
	r.city_code2,c2.name,r.city_code3,c3.name,r.population,r.composite_index,
	r.budget_city_1,r.budget_city_2,r.budget_city_3
	FROM renew_project r
	JOIN city c1 ON r.city_code1=c1.insee_code
	LEFT JOIN city c2 ON r.city_code2=c2.insee_code
	LEFT JOIN city c3 ON r.city_code3=c3.insee_code
	WHERE r.id=$1`, r.ID).Scan(&r.Reference, &r.Name, &r.Budget, &r.PRIN,
		&r.CityCode1, &r.CityName1, &r.CityCode2, &r.CityName2, &r.CityCode3,
		&r.CityName3, &r.Population, &r.CompositeIndex, &r.BudgetCity1,
		&r.BudgetCity2, &r.BudgetCity3)
}

// Create insert a renew project into database returning it's ID
func (r *RenewProject) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO renew_project (reference,name,budget,prin,
		city_code1,city_code2,city_code3,population,composite_index,budget_city_1,
		budget_city_2,budget_city_3) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`, r.Reference,
		r.Name, r.Budget, r.PRIN, r.CityCode1, r.CityCode2, r.CityCode3,
		r.Population, r.CompositeIndex, r.BudgetCity1, r.BudgetCity2, r.BudgetCity3).
		Scan(&r.ID)
	if err != nil {
		return fmt.Errorf("insert query %v", err)
	}
	err = db.QueryRow(`SELECT c1.name,c2.name,c3.name FROM renew_project r
	JOIN city c1 ON r.city_code1=c1.insee_code
	LEFT JOIN city c2 ON r.city_code2=c2.insee_code
	LEFT JOIN city c3 ON r.city_code3=c3.insee_code
	WHERE r.id=$1`, r.ID).Scan(&r.CityName1, &r.CityName2, &r.CityName3)
	return err
}

// Update modifies a renew program into database
func (r *RenewProject) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE renew_project SET reference=$1, name=$2, budget=$3,
	prin=$4,city_code1=$5,city_code2=$6,city_code3=$7,population=$8,
	composite_index=$9,budget_city_1=$10,budget_city_2=$11,budget_city_3=$12
	 WHERE id = $13`, r.Reference, r.Name, r.Budget, r.PRIN, r.CityCode1,
		r.CityCode2, r.CityCode3, r.Population, r.CompositeIndex, r.BudgetCity1,
		r.BudgetCity2, r.BudgetCity3, r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Projet de renouvellement introuvable")
	}
	err = db.QueryRow(`SELECT c1.name,c2.name,c3.name FROM renew_project r
	JOIN city c1 ON r.city_code1=c1.insee_code
	LEFT JOIN city c2 ON r.city_code2=c2.insee_code
	LEFT JOIN city c3 ON r.city_code3=c3.insee_code
	WHERE r.id=$1`, r.ID).Scan(&r.CityName1, &r.CityName2, &r.CityName3)
	return err
}

// GetAll fetches all renew projects from database
func (r *RenewProjects) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT r.id,r.reference,r.name,r.budget,r.prin,
	r.city_code1,c1.name,r.city_code2,c2.name,r.city_code3,c3.name,
	r.population,r.composite_index,r.budget_city_1,r.budget_city_2,r.budget_city_3
	FROM renew_project r
	JOIN city c1 ON c1.insee_code= r.city_code1
	LEFT JOIN city c2 ON c2.insee_code= r.city_code2
	LEFT JOIN city c3 ON c3.insee_code= r.city_code3`)
	if err != nil {
		return err
	}
	var row RenewProject
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Reference, &row.Name, &row.Budget,
			&row.PRIN, &row.CityCode1, &row.CityName1, &row.CityCode2, &row.CityName2,
			&row.CityCode3, &row.CityName3, &row.Population, &row.CompositeIndex,
			&row.BudgetCity1, &row.BudgetCity2, &row.BudgetCity3); err != nil {
			return err
		}
		r.RenewProjects = append(r.RenewProjects, row)
	}
	err = rows.Err()
	if len(r.RenewProjects) == 0 {
		r.RenewProjects = []RenewProject{}
	}
	return err
}

// Delete removes a renew project whose ID is given from database
func (r *RenewProject) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM renew_project WHERE id = $1", r.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Projet de renouvellement introuvable")
	}
	return nil
}

// Save validate the array of project and update or save all renew projects
// against the database
func (r *RenewProjectBatch) Save(db *sql.DB) error {
	for i, l := range r.Lines {
		if l.Name == "" || l.Reference == "" || l.Budget == 0 || l.CityCode1 == 0 {
			return fmt.Errorf("ligne %d : champs incorrects", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_renew_project", "reference", "name",
		"budget", "prin", "city_code1", "city_code2", "city_code3", "population",
		"composite_index", "budget_city_1", "budget_city_2", "budget_city_3"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, l := range r.Lines {
		if _, err = stmt.Exec(l.Reference, l.Name, l.Budget, l.PRIN, l.CityCode1,
			l.CityCode2, l.CityCode3, l.Population, l.CompositeIndex, l.BudgetCity1,
			l.BudgetCity2, l.BudgetCity3); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v : %s", r, err.Error())
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{`UPDATE renew_project SET name=t.name,budget=t.budget,
	prin=t.prin,city_code1=t.city_code1,city_code2=t.city_code2,
	city_code3=t.city_code3,population=t.population,composite_index=t.composite_index,
	budget_city_1=t.budget_city_1,budget_city_2=t.budget_city_2,
	budget_city_3=t.budget_city_3
	FROM temp_renew_project t WHERE t.reference = renew_project.reference`,
		`INSERT INTO renew_project (reference,name,budget,prin,city_code1,city_code2,
			city_code3,population,composite_index,budget_city_1,budget_city_2,
			budget_city_3)
	SELECT reference,name,budget, prin,city_code1,city_code2,city_code3,population,
		composite_Index,budget_city_1,budget_city_2,budget_city_3
		FROM temp_renew_project 
		WHERE reference NOT IN (SELECT reference from renew_project)`,
		`DELETE from temp_renew_project`,
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
