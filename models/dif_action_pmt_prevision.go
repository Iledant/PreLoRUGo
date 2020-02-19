package models

import (
	"database/sql"
	"fmt"
	"time"
)

// DifActionPmtPrevision model
type DifActionPmtPrevision struct {
	Sector     string  `json:"Sector"`
	ActionID   int64   `json:"ActionID"`
	ActionCode string  `json:"ActionCode"`
	ActionName string  `json:"ActionName"`
	Y0         float64 `json:"Y0"`
	Y1         float64 `json:"Y1"`
	Y2         float64 `json:"Y2"`
	Y3         float64 `json:"Y3"`
	Y4         float64 `json:"Y4"`
}

// DifActionPmtPrevisions embeddes an array of DifActionPmtPrevision for json
// export and dedicated queries
type DifActionPmtPrevisions struct {
	Lines []DifActionPmtPrevision `json:"DifActionPmtPrevision"`
}

type yearActionVal struct {
	Year     int64
	ActionID int64
	Val      float64
}

type actionItem struct {
	Sector     string
	ActionID   int64
	ActionCode string
	ActionName string
}

type actionItems struct {
	Lines []actionItem
}

// getActionRAM launches the queries with all years and action IDs including
// the 4 comming years. The query using outer and cross joins to generate
// all value or zero in order for the algorithm to work without any further
// test
func getActionRAM(db *sql.DB) ([]yearActionVal, error) {
	q := `
	WITH
		cmt AS (SELECT extract(year FROM creation_date) y,action_id,sum(value)::bigint v 
			FROM commitment
			WHERE extract (year FROM creation_date)>=2009
			AND extract(year FROM creation_date)<extract(year FROM CURRENT_DATE)
				AND value > 0
			GROUP BY 1,2),
		pmt AS (SELECT extract(year FROM f.creation_date) y,f.action_id,sum(p.value) v
			FROM payment p
			JOIN commitment f ON p.commitment_id=f.id
			WHERE extract(year FROM f.creation_date)>=2009
				AND extract(year FROM p.creation_date)-extract(year FROM f.creation_date)>=0
				AND extract(year FROM p.creation_date)<extract(year FROM CURRENT_DATE)
			GROUP BY 1,2),
		prg AS (SELECT p.year y,p.action_id,sum(p.value)::bigint v
			FROM prog p
			WHERE year=extract(year FROM CURRENT_DATE)
			GROUP BY 1,2),
		prev AS (SELECT year y,action_id,v FROM
			(SELECT EXTRACT(year FROM c.date) AS year,hf.action_id,sum(hf.value)::bigint v
        FROM housing_forecast hf
        JOIN commission c on hf.commission_id=c.id
        WHERE extract(year FROM c.date)>extract(year FROM CURRENT_DATE)
          AND extract(year FROM c.date)<extract(year FROM CURRENT_DATE)+5
        GROUP BY 1,2
      UNION ALL
      SELECT EXTRACT(year FROM c.date) AS year,cf.action_id,sum(cf.value)::bigint v
        FROM copro_forecast cf
        JOIN commission c on cf.commission_id=c.id
        WHERE extract(year FROM c.date)>extract(year FROM CURRENT_DATE)
          AND extract(year FROM c.date)<extract(year FROM CURRENT_DATE)+5
        GROUP BY 1,2          
      UNION ALL
      SELECT EXTRACT(year FROM c.date) AS year,rf.action_id,sum(rf.value)::bigint v
        FROM renew_project_forecast rf
        JOIN commission c on rf.commission_id=c.id
        WHERE extract(year FROM c.date)>extract(year FROM CURRENT_DATE)
          AND extract(year FROM c.date)<extract(year FROM CURRENT_DATE)+5
        GROUP BY 1,2                    
      ) q),
		ram AS (SELECT cmt.y,cmt.action_id,(cmt.v-COALESCE(pmt.v,0)::bigint) v FROM cmt
			LEFT OUTER JOIN pmt ON cmt.y=pmt.y AND cmt.action_id=pmt.action_id
			UNION ALL
			SELECT y,action_id,v FROM prg
			UNION ALL
			SELECT y,action_id,v FROM prev
		),
		action_id AS (SELECT distinct action_id FROM ram),
		years AS (SELECT generate_series(2009,
			extract(year FROM current_date)::int+4)::int y)
	SELECT years.y,action_id.action_id,COALESCE(ram.v,0)::double precision*0.00000001
	FROM action_id
	CROSS JOIN years
	LEFT OUTER JOIN ram ON ram.action_id=action_id.action_id AND ram.y=years.y
	WHERE action_id.action_id NOTNULL
	ORDER BY 1,2`
	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("SELECT action ram %v", err)
	}
	var (
		r  yearActionVal
		rr []yearActionVal
	)
	for rows.Next() {
		if err = rows.Scan(&r.Year, &r.ActionID, &r.Val); err != nil {
			return nil, fmt.Errorf("scan action ram %v", err)
		}
		rr = append(rr, r)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err action ram %v", err)
	}
	return rr, nil
}

func (a *actionItems) Get(db *sql.DB) error {
	q := `SELECT a.id,a.code,a.name,s.name FROM budget_action a
	JOIN budget_sector s ON a.sector_id=s.id
	ORDER BY 1`
	rows, err := db.Query(q)
	if err != nil {
		return fmt.Errorf("SELECT action datas %v", err)
	}
	var line actionItem
	for rows.Next() {
		if err = rows.Scan(&line.ActionID, &line.ActionCode, &line.ActionName,
			&line.Sector); err != nil {
			return fmt.Errorf("scan action ram %v", err)
		}
		a.Lines = append(a.Lines, line)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows err action ram %v", err)
	}
	return nil
}

func getDifRatios(db *sql.DB) ([]float64, error) {
	q := `
	WITH
  fcy AS (SELECT EXTRACT (year FROM creation_date) y,sum(value)::bigint v
    FROM commitment
    WHERE EXTRACT (year FROM creation_date)>=2009
      AND EXTRACT(year FROM creation_date)<EXTRACT(year FROM current_date)
    AND value>0
    GROUP BY 1 ORDER BY 1),
	pmy AS (SELECT EXTRACT(year FROM f.creation_date) y,
		EXTRACT(year FROM p.creation_date)-EXTRACT(year FROM f.creation_date) AS idx, 
    sum(p.value) v
	  FROM payment p
    JOIN commitment f ON p.commitment_id=f.id
    WHERE EXTRACT(year FROM f.creation_date)>=2009
      AND EXTRACT(year FROM p.creation_date)-EXTRACT(year FROM f.creation_date)>=0
      AND EXTRACT(year FROM p.creation_date)<EXTRACT(year FROM CURRENT_DATE)
    GROUP BY 1,2),
  y AS (SELECT generate_series(2009,EXTRACT(year FROM CURRENT_DATE)::int) y),
  idx AS (SELECT generate_series(0,max(idx)::int) idx FROM pmy idx),
  cpmy AS (SELECT y.y,idx.idx,COALESCE(v,0)::bigint v FROM y
    CROSS JOIN idx
    LEFT OUTER JOIN pmy on pmy.y=y.y AND idx.idx=pmy.idx
    WHERE y.y+idx.idx<EXTRACT(year FROM current_date)
    ORDER BY 1,2),
  spy AS (SELECT y,idx,sum(v) OVER (PARTITION by y ORDER BY y,idx) FROM cpmy),
  ry AS (SELECT y,0 AS idx,fcy.v FROM fcy
    UNION ALL
	  SELECT spy.y,spy.idx+1,fcy.v-spy.sum v FROM fcy JOIN spy on fcy.y=spy.y),
  r AS (SELECT ry.y,ry.idx,COALESCE(cpmy.v,0)::double precision/ry.v r
	  FROM ry JOIN cpmy on ry.y=cpmy.y AND ry.idx=cpmy.idx
	  WHERE ry.y<EXTRACT(year FROM current_date))
  SELECT idx,avg(r) FROM r WHERE idx+y>=EXTRACT(year FROM current_date) - 2
	GROUP BY 1 ORDER BY 1`
	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("select ratio %v", err)
	}
	var (
		idx    int64
		ratio  float64
		ratios []float64
	)
	for rows.Next() {
		if err = rows.Scan(&idx, &ratio); err != nil {
			return nil, fmt.Errorf("scan ratio %v", err)
		}
		ratios = append(ratios, ratio)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err ratio %v", err)
	}
	return ratios, nil
}

// Get calculates the DifActionPmtPrevision using the average differential
// ratios
func (m *DifActionPmtPrevisions) Get(db *sql.DB) error {
	ratios, err := getDifRatios(db)
	if err != nil {
		return err
	}
	ratioLen := len(ratios)
	ram, err := getActionRAM(db)
	if err != nil {
		return err
	}
	actualYear := time.Now().Year()
	var (
		actionLen, actualYearBegin, j int
		p                             DifActionPmtPrevision
	)
	for i := 1; i < len(ram); i++ {
		if ram[i].ActionID < ram[i-1].ActionID {
			actionLen = i
			break
		}
	}
	if actionLen == 0 {
		return fmt.Errorf("impossible de trouver la séquence d'action dans la requête")
	}
	for i, a := range ram {
		if a.Year == int64(actualYear) {
			actualYearBegin = i
			break
		}
	}
	if actualYearBegin == 0 {
		return fmt.Errorf("impossible de trouver l'année en cours dans la requête")
	}
	var prev float64
	m.Lines = make([]DifActionPmtPrevision, actionLen, actionLen)
	for y := 0; y < 5; y++ {
		for a := 0; a < actionLen; a++ {
			prev = 0
			j = actualYearBegin + a + y*actionLen
			p.ActionID = ram[j].ActionID
			for i := 0; i < ratioLen && j-i*actionLen >= 0; i++ {
				q := ratios[i] * ram[j-i*actionLen].Val
				if i+int(ram[j-i*actionLen].Year) != actualYear+y {
					fmt.Printf("différence de ratio+year : %+v Année : %d\n",
						ram[j-i*actionLen], actualYear+y)
				}
				prev += q
				ram[j-i*actionLen].Val -= q
			}
			m.Lines[a].ActionID = ram[j].ActionID
			switch y {
			case 0:
				m.Lines[a].Y0 = prev
			case 1:
				m.Lines[a].Y1 = prev
			case 2:
				m.Lines[a].Y2 = prev
			case 3:
				m.Lines[a].Y3 = prev
			case 4:
				m.Lines[a].Y4 = prev
			}
		}
	}
	var actions actionItems
	if err = actions.Get(db); err != nil {
		return err
	}
	var i int
	actionLen = len(actions.Lines)
	for x := 0; x < len(m.Lines); x++ {
		i = 0
		j = actionLen - 1
		for {
			if m.Lines[x].ActionID == actions.Lines[i].ActionID {
				break
			}
			if m.Lines[x].ActionID == actions.Lines[j].ActionID {
				i = j
				break
			}
			if m.Lines[x].ActionID < actions.Lines[(i+j)/2].ActionID {
				j = (i + j) / 2
			} else {
				i = (i + j) / 2
			}
		}
		m.Lines[x].Sector = actions.Lines[i].Sector
		m.Lines[x].ActionCode = actions.Lines[i].ActionCode
		m.Lines[x].ActionName = actions.Lines[i].ActionName
	}
	return nil
}
