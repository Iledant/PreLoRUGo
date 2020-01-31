package models

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

// BeneficiaryGroupData is used to embed commitment and available commitment
// linked to a beneficiary that belongs to a group
type BeneficiaryGroupData struct {
	ID              int64      `json:"ID"`
	BeneficiaryCode int64      `json:"BeneficiaryCode"`
	BeneficiaryName string     `json:"BeneficiaryName"`
	Date            time.Time  `json:"Date"`
	Value           int64      `json:"Value"`
	Name            string     `json:"Name"`
	IRISCode        NullString `json:"IRISCode"`
	Available       int64      `json:"Available"`
	Caducity        NullTime   `json:"Caducity"`
}

// BeneficiaryGroupDatas is used to embed an array of BeneficiaryData for json export
type BeneficiaryGroupDatas struct {
	BeneficiaryGroupDatas []BeneficiaryGroupData `json:"BeneficiaryGroupData"`
}

// PaginatedBeneficiaryGroupDatas is used for paginated request to fetch datas attached
// to a beneficiary that match a search pattern using PaginatedQuery
type PaginatedBeneficiaryGroupDatas struct {
	BeneficiaryGroupDatas []BeneficiaryGroupData `json:"Datas"`
	Page                  int64                  `json:"Page"`
	ItemsCount            int64                  `json:"ItemsCount"`
}

// Get fetches all paginated beneficiary group datas from database that match
// the paginated query
func (p *PaginatedBeneficiaryGroupDatas) Get(db *sql.DB, q *PaginatedQuery, ID int) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM cumulated_commitment c 
		WHERE c.year >= $1 AND c.name ILIKE $2 AND c.beneficiary_id IN 
			(SELECT beneficiary_id FROM beneficiary_belong WHERE group_id=$3)`, q.Year,
		"%"+q.Search+"%", ID).Scan(&count); err != nil {
		return errors.New("count query failed " + err.Error())
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT c.id,b.code,b.name,c.value,c.creation_date,c.name, 
		c.iris_code,c.value-COALESCE(q.added,0), c.caducity_date
	FROM cumulated_commitment c
	LEFT JOIN (SELECT sum(value) AS added, commitment_id FROM payment GROUP BY 2) q
		ON q.commitment_id = c.id
	LEFT JOIN beneficiary b ON c.beneficiary_id=b.id
	WHERE c.year >= $1 AND c.name ILIKE $2 AND c.beneficiary_id IN 
	(SELECT beneficiary_id FROM beneficiary_belong WHERE group_id=$3)
	ORDER BY 1 LIMIT `+strconv.Itoa(PageSize)+` OFFSET $4`,
		q.Year, "%"+q.Search+"%", ID, offset)
	if err != nil {
		return err
	}
	var row BeneficiaryGroupData
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.BeneficiaryCode, &row.BeneficiaryName,
			&row.Value, &row.Date, &row.Name, &row.IRISCode, &row.Available,
			&row.Caducity); err != nil {
			return err
		}
		p.BeneficiaryGroupDatas = append(p.BeneficiaryGroupDatas, row)
	}
	err = rows.Err()
	if len(p.BeneficiaryGroupDatas) == 0 {
		p.BeneficiaryGroupDatas = []BeneficiaryGroupData{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}

// GetAll fetches all beneficiary group datas from database that match the
//  paginated query
func (p *BeneficiaryGroupDatas) GetAll(db *sql.DB, q *PaginatedQuery, ID int) error {
	rows, err := db.Query(`SELECT c.id,b.code,b.name,c.value,c.creation_date,c.name, 
		c.iris_code,c.value-COALESCE(q.added,0),c.caducity_date
	FROM cumulated_commitment c
	LEFT JOIN (SELECT sum(value) AS added, commitment_id FROM payment GROUP BY 2) q
		ON q.commitment_id = c.id
	LEFT JOIN beneficiary b ON c.beneficiary_id=b.id
	WHERE c.year >= $1 AND c.name ILIKE $2 AND c.beneficiary_id IN 
	(SELECT beneficiary_id FROM beneficiary_belong WHERE group_id=$3)
	ORDER BY 1`, q.Year, "%"+q.Search+"%", ID)
	if err != nil {
		return err
	}
	var row BeneficiaryGroupData
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.BeneficiaryCode, &row.BeneficiaryName,
			&row.Value, &row.Date, &row.Name, &row.IRISCode, &row.Available,
			&row.Caducity); err != nil {
			return err
		}
		p.BeneficiaryGroupDatas = append(p.BeneficiaryGroupDatas, row)
	}
	err = rows.Err()
	if len(p.BeneficiaryGroupDatas) == 0 {
		p.BeneficiaryGroupDatas = []BeneficiaryGroupData{}
	}
	return err
}
