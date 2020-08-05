package models

import (
	"database/sql"
	"fmt"
	"time"
)

// ImportLog model
type ImportLog struct {
	Kind int64     `json:"Kind"`
	Date time.Time `json:"Date"`
}

// ImportLogs embeddes an array of ImportLog for json export
type ImportLogs struct {
	Logs []ImportLog `json:"ImportLog"`
}

// GetAll fetches all import logs from database
func (i *ImportLogs) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT kind,date FROM import_logs`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	defer rows.Close()
	var row ImportLog
	for rows.Next() {
		if err = rows.Scan(&row.Kind, &row.Date); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		i.Logs = append(i.Logs, row)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(i.Logs) == 0 {
		i.Logs = []ImportLog{}
	}
	return nil
}
