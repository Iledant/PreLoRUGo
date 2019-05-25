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
	ADD CONSTRAINT ratio_sector_id_fkey FOREIGN KEY (sector_id) REFERENCES 
		budget_sector (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`CREATE VIEW cumulated_sold_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,q.value, 
			c.sold_out,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id,
			c.copro_id,c.renew_project_id
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id`,
	`ALTER TABLE renew_project
		ADD column prin bool NOT NULL,
		ADD column city_code1 int NOT NULL,
		ADD column city_code2 int,
		ADD column city_code3 int,
		ADD CONSTRAINT city_code1_city_insee_code_fkey FOREIGN KEY (city_code1) REFERENCES
		  city(insee_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
		ADD CONSTRAINT city_code2_city_insee_code_fkey FOREIGN KEY (city_code2) REFERENCES
		  city(insee_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
		ADD CONSTRAINT city_code3_city_insee_code_fkey FOREIGN KEY (city_code3) REFERENCES
		  city(insee_code) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION`,
	`ALTER TABLE temp_renew_project
		ADD column prin bool NOT NULL,
		ADD column city_code1 int NOT NULL,
		ADD column city_code2 int,
		ADD column city_code3 int`,
}

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
