package models

import "database/sql"

// RenewProjectReportLine is used to decode the renew project query line
type RenewProjectReportLine struct {
	ID            int64      `json:"ID"`
	Reference     string     `json:"Reference"`
	Name          string     `json:"Name"`
	Budget        NullInt64  `json:"Budget"`
	Commitment    NullInt64  `json:"Commitment"`
	Payment       NullInt64  `json:"Payment"`
	LastEventName NullString `json:"LastEventName"`
	LastEventDate NullTime   `json:"LastEventDate"`
	City1Name     string     `json:"City1Name"`
	City1Cmt      NullInt64  `json:"City1Cmt"`
	City1Pmt      NullInt64  `json:"City1Pmt"`
	City2Name     NullString `json:"City2Name"`
	City2Cmt      NullInt64  `json:"City2Cmt"`
	City2Pmt      NullInt64  `json:"City2Pmt"`
	City3Name     NullString `json:"City3Name"`
	City3Cmt      NullInt64  `json:"City3Cmt"`
	City3Pmt      NullInt64  `json:"City3Pmt"`
}

// RenewProjectReport embeddes a array of RenewProjectLine fro json export
type RenewProjectReport struct {
	Lines []RenewProjectReportLine `json:"RenewProjectReport"`
}

// Get fetches all line of the renew project report
func (r *RenewProjectReport) Get(db *sql.DB) error {
	rows, err := db.Query(`WITH city_state AS (SELECT city.name,city.insee_code,
		SUM(c.value) AS cmt,SUM(p.value) AS pmt FROM city
	LEFT OUTER JOIN rp_cmt_city_join r
	ON r.city_code=city.insee_code
	LEFT OUTER JOIN cumulated_commitment c
	ON r.commitment_id = c.id
	LEFT OUTER JOIN payment p
	ON c.id=p.commitment_id GROUP BY 1,2)
	SELECT r.id,r.reference,r.name,r.budget,c.value,p.value,e.name,e.date,c1.name,
		c1.cmt,c1.pmt,c2.name,c2.cmt,c2.pmt,c3.name,c3.cmt,c3.pmt
	FROM renew_project r
	LEFT JOIN city_state c1 ON r.city_code1 = c1.insee_code
	LEFT OUTER JOIN city_state c2 ON r.city_code2 = c2.insee_code
	LEFT OUTER JOIN city_state c3 ON r.city_code3 = c3.insee_code
	LEFT OUTER JOIN 
	(SELECT renew_project_id AS renew_project_id, sum(value) AS value 
		FROM cumulated_commitment WHERE renew_project_id NOTNULL GROUP BY 1)c
	ON c.renew_project_id=r.id
	LEFT OUTER JOIN 
	(SELECT c.renew_project_id, sum(p.value) AS value 
		FROM payment p, commitment c 
		WHERE p.commitment_id=c.id AND c.renew_project_id NOTNULL GROUP BY 1) p
	ON p.renew_project_id=r.id
	LEFT OUTER JOIN
	(SELECT MAX(rp.date) AS date,rp.renew_project_id,rpt.name FROM rp_event rp
		JOIN rp_event_type rpt ON rp.rp_event_type_id=rpt.id GROUP BY 2,3) e
	ON e.renew_project_id=r.id`)
	if err != nil {
		return err
	}
	var l RenewProjectReportLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&l.ID, &l.Reference, &l.Name, &l.Budget, &l.Commitment,
			&l.Payment, &l.LastEventName, &l.LastEventDate, &l.City1Name, &l.City1Cmt,
			&l.City1Pmt, &l.City2Name, &l.City2Cmt, &l.City2Pmt, &l.City3Name,
			&l.City3Cmt, &l.City3Pmt); err != nil {
			return err
		}
		r.Lines = append(r.Lines, l)
	}
	err = rows.Err()
	return err
}
