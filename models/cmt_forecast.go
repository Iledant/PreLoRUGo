package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CmtForecast is used to calculate a commitment forecast per budget action
type CmtForecast struct {
	ActionID   int    `json:"ActionID"`
	ActionCode int64  `json:"ActionCode"`
	ActionName string `json:"ActionName"`
	Y0         int64  `json:"Y0"`
	Y1         int64  `json:"Y1"`
	Y2         int64  `json:"Y2"`
	Y3         int64  `json:"Y3"`
	Y4         int64  `json:"Y4"`
}

// CmtForecasts embeddes an array of CmtForecast
type CmtForecasts struct {
	CmtForecasts []CmtForecast `json:"CmtForecast"`
}

// GetAll fetches all commitment previsions for 5 years. Thefirst one is the
// actual year. The forecasts of the actual year are calculated using the
// commitments and the forecast of the actual year whose commission date is
// greater than the leatest commitment date
func (c *CmtForecasts) GetAll(db *sql.DB) error {
	actualYear := time.Now().Year()
	qry := fmt.Sprintf(`SELECT t.action_id, b.code,b.name,greatest(t.y0,0),
		greatest(t.y1,0),greatest(t.y2,0),greatest(t.y3,0),greatest(t.y4,0)
	FROM
	(SELECT * FROM 
		crosstab(
			'SELECT action_id, y, SUM(cmt)::bigint FROM
				(
				SELECT h.action_id,extract(year FROM c.date)::int y,SUM(h.value) cmt
					FROM housing_forecast h, commission c 
					WHERE h.commission_id=c.id 
						AND c.date>(select max(creation_date) FROM cumulated_commitment)
					GROUP BY 1,2
				UNION ALL
				SELECT r.action_id,extract(year FROM c.date)::int y,SUM(r.value) cmt 
					FROM renew_project_forecast r, commission c 
					WHERE r.commission_id=c.id 
						AND c.date>(select max(creation_date) FROM cumulated_commitment)
					GROUP BY 1,2
				UNION ALL
				SELECT co.action_id,extract(year FROM c.date)::int y,SUM(co.value) cmt 
					FROM copro_forecast co, commission c 
					WHERE co.commission_id=c.id 
						AND c.date>(select max(creation_date) FROM cumulated_commitment)
					GROUP BY 1,2
				UNION ALL
				SELECT action_id,year y,SUM(value) FROM cumulated_commitment 
					WHERE EXTRACT(year FROM creation_date)>=%d GROUP BY 1,2
				) qry
				WHERE qry.y>=%d AND qry.y<%d GROUP BY 1,2 ORDER BY 1,2',
			'SELECT m FROM GENERATE_SERIES(%d,%d) m')
		AS (action_id int,y0 bigint,y1 bigint,y2 bigint,y3 bigint,y4 bigint)
	) t
	JOIN budget_action b ON t.action_id=b.id`, actualYear, actualYear,
		actualYear+4, actualYear, actualYear+4)

	rows, err := db.Query(qry)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var r CmtForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ActionID, &r.ActionCode, &r.ActionName, &r.Y0, &r.Y1,
			&r.Y2, &r.Y3, &r.Y4); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		c.CmtForecasts = append(c.CmtForecasts, r)
	}
	err = rows.Err()
	if len(c.CmtForecasts) == 0 {
		c.CmtForecasts = []CmtForecast{}
	}
	return err
}
