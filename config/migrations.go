package config

import (
	"database/sql"
	"fmt"
	"time"
)

type migrationEntry struct {
	ID      int
	Created time.Time
	Index   int
	Query   string
}

var migrations = []string{`CREATE EXTENSION IF NOT EXISTS tablefunc;`,
	`DELETE FROM ratio`,
	`ALTER TABLE ratio
	ADD COLUMN sector_id int NOT NULL,
	ADD CONSTRAINT ratio_sector_id_fkey FOREIGN KEY (sector_id) REFERENCES budget_sector (id) MATCH SIMPLE
	ON UPDATE NO ACTION ON DELETE NO ACTION`}

// HandleMigrations check if new migrations have been created and launches them
// against the database
func HandleMigrations(db *sql.DB) error {
	var maxIdx *int
	err := db.QueryRow("SELECT max(index) FROM migration").Scan(&maxIdx)
	if err != nil {
		return fmt.Errorf("Migration select max index : %+v", err)
	}
	var i int
	if maxIdx != nil {
		i = *maxIdx + 1
	}
	for i < len(migrations) {
		if _, err = db.Exec(migrations[i]); err != nil {
			return fmt.Errorf("Migration %d : %+v", i, err)
		}
		if _, err = db.Exec(`INSERT INTO migration (created,index,query) 
		VALUES($1,$2,$3)`, time.Now(), i, migrations[i]); err != nil {
			return fmt.Errorf("Migration %d sauvegarde bdd : %+v", i, err)
		}
		i++
	}
	return nil
}
