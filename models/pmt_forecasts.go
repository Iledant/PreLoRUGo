package models

import (
	"database/sql"
	"fmt"
	"time"
)

// PmtForecast is used to calculate a payment forecast for an budget action
type PmtForecast struct {
	ActionID   int     `json:"ActionID"`
	ActionCode int64   `json:"ActionCode"`
	ActionName string  `json:"ActionName"`
	Y0         float64 `json:"Y0"`
	Y1         float64 `json:"Y1"`
	Y2         float64 `json:"Y2"`
	Y3         float64 `json:"Y3"`
	Y4         float64 `json:"Y4"`
}

// PmtForecasts embeddes an array of PmtForecast for json export
type PmtForecasts struct {
	PmtForecasts []PmtForecast `json:"PmtForecast"`
}

// Get calculates the payment forecasts per budget action by applying
// the payments ratios of the given year
func (p *PmtForecasts) Get(db *sql.DB, year int) error {
	actualYear := time.Now().Year()
	qry := fmt.Sprintf(`SELECT q.action_id, b.code, b.name, greatest(q.y0,0),
	greatest(q.y1,0),greatest(q.y2,0),greatest(q.y3,0),greatest(q.y4,0)
	FROM (select * FROM crosstab('SELECT action_id, year, 0.01*SUM(pmt) FROM ((SELECT c.action_id,
	extract(year FROM c.Creation_Date)::int+r.index AS year,
	SUM(c.value*r.ratio) AS pmt FROM cumulated_sold_commitment c, ratio r, budget_action a
	 WHERE r.year=%d AND c.action_id=a.id AND r.sector_id=a.sector_id AND c.sold_out = false 
	GROUP BY 1,2) 
	UNION ALL 
	(SELECT h.action_id, extract(year FROM c.date)::int+r.index as year, SUM(h.value*r.ratio) AS pmt 
	FROM housing_forecast h, commission c, ratio r , budget_action ba
	WHERE h.commission_id=c.id AND c.date > (select max(creation_date) FROM cumulated_commitment) 
	 AND r.year=%d AND r.sector_id = ba.sector_id AND h.action_id=ba.id
	GROUP BY 1,2) 
	UNION ALL 
	(SELECT co.action_id, extract(year FROM c.date)::int+r.index as year, SUM(co.value*r.ratio) AS pmt 
	FROM copro_forecast co, commission c, ratio r, budget_action ba 
	WHERE co.commission_id=c.id AND c.date > (select max(creation_date) FROM cumulated_commitment) 
	 AND r.year=%d AND r.sector_id = ba.sector_id AND co.action_id=ba.id
	GROUP BY 1,2) 
	UNION ALL 
	(SELECT rp.action_id, extract(year FROM c.date)::int+r.index as year, SUM(rp.value*r.ratio) AS pmt 
	FROM renew_project_forecast rp, commission c, ratio r, budget_action ba
	WHERE rp.commission_id=c.id AND 
	c.date > (select max(creation_date) FROM cumulated_commitment) 
	 AND r.year=%d AND r.sector_id=ba.sector_id AND rp.action_id=ba.id
	GROUP BY 1,2) 
	) qry 
	WHERE qry.year>=%d AND qry.year<%d GROUP BY 1,2 ORDER BY 1,2') 
	AS (action_id int, y0 double precision, 
	y1 double precision,y2 double precision, y3 double precision, 
	y4 double precision) ) q 
	JOIN budget_action b ON q.action_id=b.id ORDER BY 2`, year, year, year,
		year, actualYear, actualYear+5)
	rows, err := db.Query(qry)
	if err != nil {
		return fmt.Errorf("get request %v", err)
	}
	var r PmtForecast
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ActionID, &r.ActionCode, &r.ActionName, &r.Y0, &r.Y1,
			&r.Y2, &r.Y3, &r.Y4); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.PmtForecasts = append(p.PmtForecasts, r)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.PmtForecasts) == 0 {
		p.PmtForecasts = []PmtForecast{}
	}
	return nil
}
