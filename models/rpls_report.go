package models

import (
	"database/sql"
	"fmt"
	"time"
)

// RPLSReportParams is used to embeddes RPLSReport queries params
type RPLSReportParams struct {
	RPLSYear  int64   `json:"RPLSYear"`
	RPLSMin   float64 `json:"RPLSMin"`
	RPLSMax   float64 `json:"RPLSMax"`
	FirstYear int64   `json:"FirstYear"`
	LastYear  int64   `json:"LastYear"`
}

// RPLSReportLine is used to decode one line of the RPLS Report query
type RPLSReportLine struct {
	Value int64 `json:"Value"`
	Dpt   int64 `json:"Dpt"`
}

// RPLSReport embeddes an array of RPLSReportLine for json export
type RPLSReport struct {
	Lines []RPLSReportLine `json:"RPLSReport"`
}

// RPLSDetailedReportLine is used to decode one line of the detailed RPLS Report query
type RPLSDetailedReportLine struct {
	CreationDate time.Time `json:"CreationDate"`
	IrisCode     string    `json:"IrisCode"`
	Value        int64     `json:"Value"`
	Reference    string    `json:"Reference"`
	Address      string    `json:"Address"`
	PLAI         int64     `json:"PLAI"`
	PLUS         int64     `json:"PLUS"`
	PLS          int64     `json:"PLS"`
	InseeCode    int64     `json:"InseeCode"`
	CityName     string    `json:"CityName"`
	RPLS         float64   `json:"RPLS"`
}

// RPLSDetailedReport embeddes an array of RPLSDetailedReportLine for json export
type RPLSDetailedReport struct {
	Lines []RPLSDetailedReportLine `json:"RPLSDetailedReport"`
}

// GetAll fetches all lines of the RPLS report using the query params given in p
func (r *RPLSReport) GetAll(p *RPLSReportParams, db *sql.DB) error {
	rows, err := db.Query(`SELECT SUM(c.value) AS value, r.insee_code/1000 AS dpt
	FROM cumulated_commitment c
	JOIN housing h ON h.id=c.housing_id
	JOIN rpls r ON h.zip_code=r.insee_code
		WHERE EXTRACT(year FROM c.creation_date)>=$1 
			AND EXTRACT(year FROM c.creation_date)<=$2
			AND c.housing_id NOTNULL AND r.ratio >=$3 AND r.ratio <=$4 AND r.year=$5
	GROUP BY 2;`, p.FirstYear, p.LastYear, p.RPLSMin, p.RPLSMax, p.RPLSYear)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row RPLSReportLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.Value, &row.Dpt); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.Lines = append(r.Lines, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)

	}
	if len(r.Lines) == 0 {
		r.Lines = []RPLSReportLine{}
	}
	return nil
}

// GetAll fetches all lines of the detailed RPLS report using the query params
// given in p
func (r *RPLSDetailedReport) GetAll(p *RPLSReportParams, db *sql.DB) error {
	rows, err := db.Query(`SELECT c.creation_date, c.iris_code, c.value, 
		h.reference,h.address,h.plai,h.plus,h.pls,r.insee_code,city.name,r.ratio 
	FROM cumulated_commitment c
	JOIN housing h ON h.id=c.housing_id
	JOIN rpls r ON h.zip_code=r.insee_code
	JOIN city ON r.insee_code=city.insee_code
		WHERE EXTRACT(year FROM c.creation_date)>=$1
			AND EXTRACT(year FROM c.creation_date)<=$2 AND c.housing_id NOTNULL
			AND r.ratio >=$3 AND r.ratio <=$4 AND r.year=$5`, p.FirstYear, p.LastYear,
		p.RPLSMin, p.RPLSMax, p.RPLSYear)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var row RPLSDetailedReportLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.CreationDate, &row.IrisCode, &row.Value, &row.Reference,
			&row.Address, &row.PLAI, &row.PLUS, &row.PLS, &row.InseeCode, &row.CityName,
			&row.RPLS); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		r.Lines = append(r.Lines, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(r.Lines) == 0 {
		r.Lines = []RPLSDetailedReportLine{}
	}
	return nil
}
