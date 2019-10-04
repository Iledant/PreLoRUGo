package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// HousingSummaryLine is used to decode one line of the housing summary excel file
// that is imported to created the housing lines
type HousingSummaryLine struct {
	InseeCode     int64  `json:"InseeCode"`
	Address       string `json:"Address"`
	PLS           int64  `json:"PLS"`
	PLAI          int64  `json:"PLAI"`
	PLUS          int64  `json:"PLUS"`
	IRISCode      string `json:"IRISCode"`
	ReferenceCode string `json:"ReferenceCode"`
	ANRU          bool   `json:"Anru"`
}

// HousingSummary embeddes an array of HousingSummaryLine for json upload
// and creating housing lines
type HousingSummary struct {
	Lines []HousingSummaryLine `json:"HousingSummary"`
}

// validate checks if fields are correctly filled
func (h *HousingSummary) validate() error {
	for i, l := range h.Lines {
		if l.InseeCode == 0 {
			return fmt.Errorf("line %d InseeCode nul", i+1)
		}
		if l.Address == "" {
			return fmt.Errorf("line %d Address nul", i+1)
		}
		if l.PLS == 0 && l.PLAI == 0 && l.PLUS == 0 {
			return fmt.Errorf("line %d PLS PLAI and PLUS nul", i+1)
		}
		if l.IRISCode == "" {
			return fmt.Errorf("line %d IrisCode nul", i+1)
		}
		if l.ReferenceCode == "" {
			return fmt.Errorf("line %d ReferenceCode nul", i+1)
		}
	}
	return nil
}

// fetchMaxDptRef calculates the max number of the housing reference in the
// database of the current year
func fetchMaxDptRef(decade int, tx *sql.Tx) (map[int]int, error) {
	maxDptRef := make(map[int]int)
	rows, err := tx.Query(`select substring(reference FROM 4 for 2)::int,
		max(substring(reference FROM 8 FOR 3)::int) 
	FROM housing WHERE substring(reference FROM 6 FOR 2)=$1::text
	GROUP BY 1 ORDER BY 1`, decade)
	if err != nil {
		return nil, fmt.Errorf("maxDptRef query %v", err)
	}
	defer rows.Close()
	var dpt, val int
	for d := range []int{75, 77, 78, 91, 92, 93, 94, 95} {
		maxDptRef[d] = 0
	}
	for rows.Next() {
		if err = rows.Scan(&dpt, &val); err != nil {
			return nil, fmt.Errorf("maxDptRef scan %v", err)
		}
		maxDptRef[dpt] = val
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("maxDptRef rows %v", err)
	}
	return maxDptRef, nil
}

type composedHousing struct {
	SummaryRef string
	HousingRef string
	ZipCode    int64
	Address    string
	PLS        int64
	PLUS       int64
	PLAI       int64
	ANRU       bool
}

type refsIRIS struct {
	SummaryRef string
	HousingRef string
}

// Save import a housing summary batch, validates it and process it to create
// new housing lines
func (h *HousingSummary) Save(db *sql.DB) error {
	if err := h.validate(); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_housing_summary", "insee_code",
		"address", "pls", "plai", "plus", "iris_code", "reference_code", "anru"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i, l := range h.Lines {
		if _, err = stmt.Exec(l.InseeCode, l.Address, l.PLS, l.PLAI, l.PLUS,
			l.IRISCode, l.ReferenceCode, l.ANRU); err != nil {
			tx.Rollback()
			return fmt.Errorf("line %d %v", i, err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	year := time.Now().Year()
	decade := year % 100
	maxDptRef, err := fetchMaxDptRef(decade, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows, err := tx.Query(`SELECT reference_code,insee_code,max(address),SUM(pls),
		SUM(plai),SUM(plus),bool_or(anru) FROM temp_housing_summary 
	WHERE reference_code NOT IN
		(SELECT reference_code FROM housing_summary WHERE year=$1)
		GROUP BY 1,2`, year)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("select temp_housing_summary %v", err)
	}
	defer rows.Close()
	var hsg composedHousing
	var hh []composedHousing
	var dpt int
	for rows.Next() {
		if err = rows.Scan(&hsg.SummaryRef, &hsg.ZipCode, &hsg.Address,
			&hsg.PLS, &hsg.PLAI, &hsg.PLUS, &hsg.ANRU); err != nil {
			tx.Rollback()
			return fmt.Errorf("scan temp_housing_summary %v", err)
		}
		dpt = int(hsg.ZipCode / 1000)
		hsg.HousingRef = fmt.Sprintf("LLS%d%d%03d", dpt, decade, maxDptRef[dpt]+1)
		hh = append(hh, hsg)
		maxDptRef[dpt]++
	}
	for _, hsg = range hh {
		if _, err = tx.Exec(`INSERT INTO housing (reference,address,zip_code,plai,
			plus,pls,anru) VALUES($1,$2,$3,$4,$5,$6,$7)`, hsg.HousingRef, hsg.Address,
			hsg.ZipCode, hsg.PLAI, hsg.PLUS, hsg.PLS, hsg.ANRU); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert housing %v", err)
		}
		if _, err = tx.Exec(`INSERT INTO housing_summary (year,housing_ref,
			import_ref,iris_code) SELECT $1,$2::varchar,$3::varchar,iris_code FROM temp_housing_summary 
			WHERE reference_code=$3`, year, hsg.HousingRef, hsg.SummaryRef); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert housing_summary %v", err)
		}
	}
	queries := []string{
		`UPDATE commitment SET housing_id=q.housing_id, 
	copro_id=NULL, renew_project_id=NULL FROM
		(SELECT c.id AS commitment_id, h.id AS housing_id FROM housing h 
			JOIN housing_summary hs ON h.reference = hs.housing_ref
			JOIN commitment c ON hs.iris_code=c.iris_code) q 
	WHERE commitment.id=q.commitment_id`,
		`DELETE FROM temp_housing_summary`,
	}
	for i, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			tx.Rollback()
			return fmt.Errorf("query %d commitment query %v", i, err)
		}
	}
	tx.Commit()
	return nil
}
