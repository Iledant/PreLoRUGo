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
	`ALTER TABLE renew_project_forecast ADD COLUMN action_id int NOT NULL,
		ADD CONSTRAINT renew_project_action_id_fkey FOREIGN KEY (action_id) 
		REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`ALTER TABLE temp_renew_project_forecast ADD COLUMN action_id int NOT NULL`,
	`ALTER TABLE copro_forecast ADD COLUMN action_id int NOT NULL,
		ADD CONSTRAINT copro_action_id_fkey FOREIGN KEY (action_id) 
		REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`ALTER TABLE temp_copro_forecast ADD COLUMN action_code bigint NOT NULL`,
	`ALTER TABLE community ADD COLUMN department_id int,
		ADD CONSTRAINT community_department_id_fkey FOREIGN KEY (department_id) 
		REFERENCES department (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`ALTER TABLE temp_community ADD COLUMN department_code int`,
	`ALTER TABLE renew_project
		ADD COLUMN budget_city_1 int,
		ADD COLUMN budget_city_2 int,
		ADD COLUMN budget_city_3 int`,
	`ALTER TABLE temp_renew_project
		ADD COLUMN budget_city_1 int,
		ADD COLUMN budget_city_2 int,
		ADD COLUMN budget_city_3 int`,
	`ALTER TABLE commitment
		ADD COLUMN cadicity_date date DEFAULT null`,
	`CREATE OR REPLACE VIEW cumulated_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.caducity_date,c.name,
		  q.value,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id, c.copro_id,
			c.renew_project_id
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`,
	`CREATE OR REPLACE VIEW cumulated_sold_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.caducity_date,c.name,
			q.value, c.sold_out,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id,
			c.copro_id,c.renew_project_id
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`,
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
