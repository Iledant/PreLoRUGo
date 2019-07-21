package models

import "database/sql"

// RPPerCommunityLine is used to decode one line of the report query to get
// a report on commitment and payment of a renew projects
type RPPerCommunityLine struct {
	CommunityID     int64  `json:"CommunityID"`
	CommunityName   string `json:"CommunityName"`
	CommunityBudget int64  `json:"CommunityBudget"`
	Commitment      int64  `json:"Commitment"`
	Payment         int64  `json:"Payment"`
}

// RPPerCommunityReport embeddes an array of RPPerCommunityLine for json export
type RPPerCommunityReport struct {
	Lines []RPPerCommunityLine `json:"RPPerCommunityReport"`
}

// Get fetches the report that calculates commitments and payments linked to
// a renew projet per community
func (r *RPPerCommunityReport) Get(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT c.id,c.name,bud.budget,
		COALESCE(q.cmt,0),COALESCE(q.pmt,0) FROM community c
	JOIN 
	(SELECT sum(q.budget) AS budget,q.id FROM
		(SELECT SUM(COALESCE(r.budget_city_1,0)) AS budget,co.id 
				FROM renew_project r, city ci, community co
				WHERE r.city_code1 = ci.insee_code AND ci.community_id=co.id
				GROUP BY 2
		UNION ALL
		SELECT SUM(COALESCE(r.budget_city_2,0)) AS budget,co.id 
				FROM renew_project r, city ci, community co
				WHERE r.city_code2 = ci.insee_code AND ci.community_id=co.id
				GROUP BY 2
		UNION ALL
		SELECT SUM(COALESCE(r.budget_city_3,0)) AS budget,co.id 
				FROM renew_project r, city ci, community co
				WHERE r.city_code3 = ci.insee_code AND ci.community_id=co.id
				GROUP BY 2
		) q
		GROUP BY 2
	) bud
	ON bud.id = c.id
	LEFT OUTER JOIN
	(SELECT SUM(cmt.value) AS cmt,SUM(pmt.value) AS pmt,co.id 
		FROM commitment cmt, rp_cmt_city_join r, city c, community co, payment pmt
		WHERE cmt.id=r.commitment_id AND r.city_code=c.insee_code 
			AND c.community_id=co.id AND pmt.commitment_id=cmt.id
		GROUP BY 3
	) q
	ON c.id=q.id
	`)
	if err != nil {
		return err
	}
	var l RPPerCommunityLine
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&l.CommunityID, &l.CommunityName, &l.CommunityBudget,
			&l.Commitment, &l.Payment); err != nil {
			return err
		}
		r.Lines = append(r.Lines, l)
	}
	err = rows.Err()
	return err
}
