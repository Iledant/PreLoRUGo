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

var migrations = []string{`ALTER TABLE copro ALTER COLUMN reference TYPE varchar(25)`,
	`ALTER TABLE reject_project_forecast ADD COLUMN action_id int NOT NULL,
		ADD CONSTRAINT renew_project_action_id_fkey FOREIGN KEY (action_id) 
		REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`ALTER TABLE temp_reject_project_forecast ADD COLUMN action_id int NOT NULL`,
	`ALTER TABLE copro_project_forecast ADD COLUMN action_id int NOT NULL,
		ADD CONSTRAINT copro_action_id_fkey FOREIGN KEY (action_id) 
		REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`ALTER TABLE temp_copro_forecast ADD COLUMN action_id int NOT NULL`}

// handleMigrations check if new migrations have been created and launches them
// against the database
func handleMigrations(db *sql.DB) error {
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
