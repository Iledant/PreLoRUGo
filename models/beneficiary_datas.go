package models

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

// BeneficiaryData is used to embed commitment and available commitment
type BeneficiaryData struct {
	ID        int64      `json:"ID"`
	Date      time.Time  `json:"Date"`
	Value     int64      `json:"Value"`
	Name      string     `json:"Name"`
	IRISCode  NullString `json:"IRISCode"`
	Available int64      `json:"Available"`
	Caducity  NullTime   `json:"Caducity"`
}

// BeneficiaryDatas is used to embed an array of BeneficiaryData for json export
type BeneficiaryDatas struct {
	BeneficiaryDatas []BeneficiaryData `json:"BeneficiaryData"`
}

// PaginatedBeneficiaryDatas is used for paginated request to fetch datas attached
// to a beneficiary that match a search pattern using PaginatedQuery
type PaginatedBeneficiaryDatas struct {
	BeneficiaryDatas []BeneficiaryData `json:"Datas"`
	Page             int64             `json:"Page"`
	ItemsCount       int64             `json:"ItemsCount"`
}

// Get fetches all paginated beneficiary datas from database that match the paginated query
func (p *PaginatedBeneficiaryDatas) Get(db *sql.DB, q *PaginatedQuery, ID int) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM cumulated_commitment c 
		WHERE c.year >= $1 AND c.name ILIKE $2 AND c.beneficiary_id=$3`, q.Year,
		"%"+q.Search+"%", ID).Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT c.id, c.value, c.creation_date, c.name, 
		c.iris_code,c.value-COALESCE(q.added,0), c.caducity_date
	FROM cumulated_commitment c
	LEFT JOIN (SELECT sum(value) AS added, commitment_id FROM payment GROUP BY 2) q
		ON q.commitment_id = c.id
	WHERE c.year >= $1 AND c.name ILIKE $2 AND c.beneficiary_id = $3
	ORDER BY 1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $4`,
		q.Year, "%"+q.Search+"%", ID, offset)
	if err != nil {
		return err
	}
	var row BeneficiaryData
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Value, &row.Date, &row.Name, &row.IRISCode,
			&row.Available, &row.Caducity); err != nil {
			return err
		}
		p.BeneficiaryDatas = append(p.BeneficiaryDatas, row)
	}
	err = rows.Err()
	if len(p.BeneficiaryDatas) == 0 {
		p.BeneficiaryDatas = []BeneficiaryData{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}

// GetAll fetches all beneficiary datas from database that match the paginated query
func (p *BeneficiaryDatas) GetAll(db *sql.DB, q *PaginatedQuery, ID int) error {
	rows, err := db.Query(`SELECT c.id, c.value, c.creation_date, c.name, 
		c.iris_code,c.value-COALESCE(q.added,0), c.caducity_date
	FROM cumulated_commitment c
	LEFT JOIN (SELECT sum(value) AS added, commitment_id FROM payment GROUP BY 2) q
		ON q.commitment_id = c.id
	WHERE c.year >= $1 AND c.name ILIKE $2 AND c.beneficiary_id = $3
	ORDER BY 1`, q.Year, "%"+q.Search+"%", ID)
	if err != nil {
		return err
	}
	var row BeneficiaryData
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.Value, &row.Date, &row.Name, &row.IRISCode,
			&row.Available, &row.Caducity); err != nil {
			return err
		}
		p.BeneficiaryDatas = append(p.BeneficiaryDatas, row)
	}
	err = rows.Err()
	if len(p.BeneficiaryDatas) == 0 {
		p.BeneficiaryDatas = []BeneficiaryData{}
	}
	return err
}
