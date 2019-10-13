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

var migrations = []string{`ALTER TABLE copro ALTER COLUMN reference TYPE varchar(25)`, // 0
	`ALTER TABLE renew_project_forecast ADD COLUMN action_id int NOT NULL,
		ADD CONSTRAINT renew_project_action_id_fkey FOREIGN KEY (action_id) 
		REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`, // 1
	`ALTER TABLE temp_renew_project_forecast ADD COLUMN action_id int NOT NULL`, // 2
	`ALTER TABLE copro_forecast ADD COLUMN action_id int NOT NULL,
		ADD CONSTRAINT copro_action_id_fkey FOREIGN KEY (action_id) 
		REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`, // 3
	`ALTER TABLE temp_copro_forecast ADD COLUMN action_code bigint NOT NULL`, // 4
	`ALTER TABLE community ADD COLUMN department_id int,
		ADD CONSTRAINT community_department_id_fkey FOREIGN KEY (department_id) 
		REFERENCES department (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`, // 5
	`ALTER TABLE temp_community ADD COLUMN department_code int`, // 6
	`ALTER TABLE renew_project
		ADD COLUMN budget_city_1 int,
		ADD COLUMN budget_city_2 int,
		ADD COLUMN budget_city_3 int`, // 7
	`ALTER TABLE temp_renew_project
		ADD COLUMN budget_city_1 int,
		ADD COLUMN budget_city_2 int,
		ADD COLUMN budget_city_3 int`, // 8
	`ALTER TABLE commitment
		ADD COLUMN caducity_date date DEFAULT null`, // 9
	`ALTER TABLE temp_commitment
		ADD COLUMN caducity_date date DEFAULT null`, // 10
	`ALTER TABLE copro
		ALTER COLUMN reference TYPE varchar(60)`, // 11
	`ALTER TABLE renew_project_forecast
		ADD COLUMN project varchar(150) DEFAULT NULL`, // 12
	`ALTER TABLE temp_renew_project_forecast
		ADD COLUMN project varchar(150) DEFAULT NULL`, // 13
}

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
