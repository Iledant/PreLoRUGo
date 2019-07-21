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
}

// RenewProjectReport embeddes a array of RenewProjectLine fro json export
type RenewProjectReport struct {
	Lines []RenewProjectReportLine `json:"RenewProjectReport"`
}

// Get fetches all line of the renew project report
func (r *RenewProjectReport) Get(db *sql.DB) error {
	rows, err := db.Query(`SELECT r.id,r.reference,r.name,r.budget,c.value,
		p.value,e.name,e.date
	FROM renew_project r
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
			&l.Payment, &l.LastEventName, &l.LastEventDate); err != nil {
			return err
		}
		r.Lines = append(r.Lines, l)
	}
	err = rows.Err()
	return err
}
