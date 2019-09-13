package config

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kataras/iris"

	"github.com/Iledant/PreLoRUGo/models"
)

func getNames(db *sql.DB, tableType string) ([]string, error) {
	var tables []string
	var table string

	rows, err := db.Query(`select table_name from information_schema.tables 
	where table_schema='public' and table_type =$1;`, tableType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

// dropTriggersAndFunctions delete triggers and linked functions
func dropTriggersAndFunctions(db *sql.DB) error {
	queries := []string{
		`DROP TRIGGER IF EXISTS cmt_stamp ON commitment;`,
		`DROP FUNCTION IF EXISTS log_cmt();`,
		`DROP TRIGGER IF EXISTS pmt_stamp ON commitment;`,
		`DROP FUNCTION IF EXISTS log_pmt();`,
	}
	for i, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("drop trigger %d query %v", i, err)
		}
	}
	return nil
}

// dropAllTables delete all table from database for test purpose
func dropAllTables(db *sql.DB, app *iris.Application) error {
	views, err := getNames(db, "VIEW")
	if err != nil {
		return fmt.Errorf("get view names : %v", err)
	}
	if len(views) > 0 {
		if _, err = db.Exec("drop view " + strings.Join(views, ",")); err != nil {
			return fmt.Errorf("drop views : %v", err)
		}
		app.Logger().Infof("%d views dropped", len(views))
	}
	tables, err := getNames(db, "BASE TABLE")
	if err != nil {
		return fmt.Errorf("get table names : %v", err)
	}
	if len(tables) > 0 {
		if _, err = db.Exec("drop table " + strings.Join(tables, ",")); err != nil {
			return fmt.Errorf("drop tables : %v", err)
		}
		app.Logger().Infof("%d tables dropped", len(views))
	}
	return nil
}

var initQueries = []string{`CREATE EXTENSION IF NOT EXISTS tablefunc`,
	`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name varchar(50) NOT NULL,
		email varchar(120) NOT NULL,
		password varchar(120) NOT NULL,
		rights int NOT NULL
		);`, // 0 : users
	`CREATE TABLE IF NOT EXISTS department (
			id SERIAL PRIMARY KEY,
			code int NOT NULL,
			name varchar(20)
		);`, // 1 department
	`CREATE TABLE IF NOT EXISTS community (
	    id SERIAL PRIMARY KEY,
	    code varchar(15) NOT NULL,
			name varchar(150) NOT NULL,
			department_id int,
			FOREIGN KEY (department_id) REFERENCES department (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 2 : community
	`CREATE TABLE IF NOT EXISTS temp_community (
	    code varchar(15) NOT NULL,
			name varchar(150) NOT NULL,
			department_code int
		);`, // 3 : temp_community
	`CREATE TABLE IF NOT EXISTS city (
	    insee_code int NOT NULL PRIMARY KEY,
	    name varchar(50) NOT NULL,
			community_id int,
			qpv boolean NOT NULL,
			FOREIGN KEY (community_id) REFERENCES community (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 4 : city
	`CREATE TABLE IF NOT EXISTS temp_city (
	    insee_code int NOT NULL UNIQUE,
	    name varchar(50) NOT NULL,
	    community_code varchar(15),
			qpv boolean NOT NULL
		);`, // 5 : temp_city
	`CREATE TABLE IF NOT EXISTS copro (
			id SERIAL PRIMARY KEY,
			reference varchar(60) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint,
			FOREIGN KEY (zip_code) REFERENCES city (insee_code) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 6 : copro
	`CREATE TABLE IF NOT EXISTS temp_copro (
			reference varchar(150) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint
		);`, // 7 : temp_copro
	`CREATE TABLE IF NOT EXISTS budget_sector (
			id SERIAL PRIMARY KEY,
			name varchar(20) NOT NULL,
	    full_name varchar(150)
		);`, // 8 : budget_sector
	`CREATE TABLE IF NOT EXISTS budget_action (
			id SERIAL PRIMARY KEY,
			code bigint NOT NULL,
			name varchar(250) NOT NULL,
			sector_id int
		);`, // 9 : budget_action
	`CREATE TABLE IF NOT EXISTS renew_project (
			id SERIAL PRIMARY KEY,
			reference varchar(15) NOT NULL UNIQUE,
			name varchar(150) NOT NULL,
			budget bigint NOT NULL,
			prin bool NOT NULL,
			city_code1 int NOT NULL,
			city_code2 int,
			city_code3 int,			
			budget_city_1 int,			
			budget_city_2 int,			
			budget_city_3 int,			
			population int,
			composite_index int,
			FOREIGN KEY (city_code1) REFERENCES city(insee_code) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (city_code2) REFERENCES city(insee_code) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (city_code3) REFERENCES city(insee_code) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 10 : renew_project
	`CREATE TABLE IF NOT EXISTS temp_renew_project (
			reference varchar(15) NOT NULL UNIQUE,
			name varchar(150) NOT NULL,
			budget bigint NOT NULL,	
			prin bool NOT NULL,
			city_code1 int NOT NULL,
			city_code2 int,
			city_code3 int,			
			budget_city_1 int,			
			budget_city_2 int,			
			budget_city_3 int,			
			population int,
			composite_index int
		);`, // 11 : temp_renew_project
	`CREATE TABLE IF NOT EXISTS housing (
	    id SERIAL PRIMARY KEY,
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
			anru boolean NOT NULL,
			FOREIGN KEY (zip_code) REFERENCES city (insee_code) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 12 : housing
	`CREATE TABLE IF NOT EXISTS temp_housing (
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
	    anru boolean NOT NULL
		);`, // 13 : temp_housing
	`CREATE TABLE IF NOT EXISTS beneficiary (
	    id SERIAL PRIMARY KEY,
	    code int NOT NULL,
	    name varchar(120) NOT NULL
		);`, // 14 : beneficiary
	`CREATE TABLE IF NOT EXISTS commitment (
	    id SERIAL PRIMARY KEY,
	    year int NOT NULL,
	    code varchar(5) NOT NULL,
	    number int NOT NULL,
	    line int NOT NULL,
	    creation_date date NOT NULL,
	    modification_date date NOT NULL,
			caducity_date date,
			name varchar(150) NOT NULL,
	    value bigint NOT NULL,
	    beneficiary_id int NOT NULL,
			iris_code varchar(20),
			sold_out boolean NOT NULL,
			action_id int,
			housing_id int,
			copro_id int,
			renew_project_id int,
			FOREIGN KEY (beneficiary_id) REFERENCES beneficiary(id) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (housing_id) REFERENCES housing(id) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (copro_id) REFERENCES copro(id) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (renew_project_id) REFERENCES renew_project(id) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (action_id) REFERENCES budget_action(id) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 15 : commitment
	`CREATE TABLE IF NOT EXISTS temp_commitment (
	    year int NOT NULL,
	    code varchar(5) NOT NULL,
	    number int NOT NULL,
	    line int NOT NULL,
	    creation_date date NOT NULL,
			modification_date date NOT NULL,
			caducity_date date,
	    name varchar(150) NOT NULL,
	    value bigint NOT NULL,
	    beneficiary_code int NOT NULL,
	    beneficiary_name varchar(150) NOT NULL,
			iris_code varchar(20),
			sold_out boolean NOT NULL,
			sector varchar(5) NOT NULL,
			action_code bigint,
			action_name varchar(150)
		);`, // 16 : temp_commitment
	`CREATE TABLE IF NOT EXISTS payment (
	    id SERIAL PRIMARY KEY,
	    commitment_id int,
	    commitment_year int NOT NULL,
	    commitment_code varchar(5) NOT NULL,
	    commitment_number int NOT NULL,
	    commitment_line int NOT NULL,
	    year int NOT NULL,
	    creation_date date NOT NULL,
			modification_date date NOT NULL,
			number int NOT NULL,
			value bigint NOT NULL,
			FOREIGN KEY (commitment_id) REFERENCES commitment (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 17 : payment
	`CREATE TABLE IF NOT EXISTS temp_payment (
	    commitment_year int NOT NULL,
	    commitment_code varchar(5) NOT NULL,
	    commitment_number int NOT NULL,
	    commitment_line int NOT NULL,
	    year int NOT NULL,
	    creation_date date NOT NULL,
			modification_date date NOT NULL,
			number int NOT NULL,
	    value bigint NOT NULL
		);`, // 18 : temp_payment
	`CREATE TABLE IF NOT EXISTS commission (
	    id SERIAL PRIMARY KEY,
	    name varchar(140) NOT NULL,
	    date date
		);`, // 19 : commission
	`CREATE TABLE IF NOT EXISTS renew_project_forecast (
	    id SERIAL PRIMARY KEY,
	    commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
			renew_project_id int NOT NULL,
			action_id int NOT NULL,
			FOREIGN KEY (renew_project_id) REFERENCES renew_project (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 20 : renew_project_forecast
	`CREATE TABLE IF NOT EXISTS temp_renew_project_forecast (
			id int NOT NULL,
			commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
			renew_project_id int NOT NULL,
			action_code bigint NOT NULL
		);`, // 21 : temp_renew_project_forecast
	`CREATE TABLE IF NOT EXISTS copro_forecast (
	    id SERIAL PRIMARY KEY,
	    commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
			copro_id int NOT NULL,
			action_id int NOT NULL,
			FOREIGN KEY (copro_id) REFERENCES copro (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 22 : copro_forecast
	`CREATE TABLE IF NOT EXISTS temp_copro_forecast (
			id int NOT NULL,
			commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
			copro_id int NOT NULL,
			action_code bigint NOT NULL
		);`, // 23 : temp_copro_forecast
	`CREATE OR REPLACE VIEW cumulated_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,
		  q.value,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id, c.copro_id,
			c.renew_project_id,c.caducity_date
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`, // 24 : cumulated_commitment view
	`CREATE OR REPLACE VIEW cumulated_sold_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,
			q.value, c.sold_out,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id,
			c.copro_id,c.renew_project_id,c.caducity_date
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`, // 25 : cumulated_sold_commitment view
	`CREATE TABLE IF NOT EXISTS ratio (
			id SERIAL PRIMARY KEY,
			year int NOT NULL,
			sector_id int NOT NULL,
			index int NOT NULL,
			ratio double precision NOT NULL,
			FOREIGN KEY (sector_id) REFERENCES budget_sector (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
			);`, // 26 : ratio
	`CREATE TABLE IF NOT EXISTS housing_forecast (
				id SERIAL PRIMARY KEY,
				commission_id int NOT NULL,
				value bigint NOT NULL,
				comment text,
				action_id int NOT NULL,
				FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
				ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
				FOREIGN KEY (commission_id) REFERENCES commission (id) MATCH SIMPLE
				ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
			);`, // 27 : renew_project_forecast
	`CREATE TABLE IF NOT EXISTS temp_housing_forecast (
				id int NOT NULL,
				commission_id int NOT NULL,
				value bigint NOT NULL,
				comment text,
				action_id int NOT NULL
			);`, // 28 : temp_renew_project_forecast
	`CREATE TABLE IF NOT EXISTS housing_commitment (
			iris_code varchar(20),
	    reference varchar(100)
			);`, // 29 housing_commitment
	`CREATE TABLE IF NOT EXISTS copro_commitment (
			iris_code varchar(20),
	    reference varchar(100)
			);`, // 30 copro_commitment
	`CREATE TABLE IF NOT EXISTS migration (
		id SERIAL PRIMARY KEY,
		created timestamp NOT NULL,
		index int NOT NULL,
		query text
	);`, // 31 migration
	`CREATE TABLE IF NOT EXISTS rp_event_type (
		id SERIAL PRIMARY KEY,
		name varchar(100) NOT NULL
	);`, // 32 rp_event_type
	`CREATE TABLE IF NOT EXISTS rp_event (
		id SERIAL PRIMARY KEY,
		renew_project_id int NOT NULL,
		rp_event_type_id int NOT NULL,
		date date NOT NULL,
		comment text,
		FOREIGN KEY (renew_project_id) REFERENCES renew_project (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
		FOREIGN KEY (rp_event_type_id) REFERENCES rp_event_type (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 33 rp_event
	`CREATE TABLE IF NOT EXISTS rp_cmt_city_join (
		id SERIAL PRIMARY KEY,
		commitment_id int NOT NULL UNIQUE,
		city_code int NOT NULL,
		FOREIGN KEY (commitment_id) REFERENCES commitment (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
		FOREIGN KEY (city_code) REFERENCES city (insee_code) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 34 rp_cmt_city_join
	`CREATE TABLE IF NOT EXISTS pre_prog (
		id SERIAL PRIMARY KEY,
		year int NOT NULL,
		commission_id int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		comment text,
		action_id int,
		FOREIGN KEY (commission_id) REFERENCES commission (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
		FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 35 pre_prog
	`CREATE TABLE IF NOT EXISTS temp_pre_prog (
		commission_id int NOT NULL,
		year int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		comment text,
		action_id int
	);`, // 36 temp_pre_prog
	`CREATE TABLE IF NOT EXISTS prog (
		id SERIAL PRIMARY KEY,
		year int NOT NULL,
		commission_id int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		comment text,
		action_id int,
		FOREIGN KEY (commission_id) REFERENCES commission (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
		FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 37 prog
	`CREATE TABLE IF NOT EXISTS temp_prog (
		commission_id int NOT NULL,
		year int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		comment text,
		action_id int
	);`, // 38 temp_prog
	`CREATE TABLE IF NOT EXISTS rpls(
		id SERIAL PRIMARY KEY,
		insee_code int NOT NULL,
		year int NOT NULL,
		ratio double precision NOT NULL,
		FOREIGN KEY (insee_code) REFERENCES city (insee_code) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 39 rpls		`
	`CREATE TABLE IF NOT EXISTS temp_rpls(
		insee_code int NOT NULL,
		year int NOT NULL,
		ratio double precision NOT NULL
	);`, // 40 temp_rpls
	`CREATE TABLE IF NOT EXISTS import_logs(
		kind int UNIQUE,
		date date
	);`, // 41
	`CREATE OR REPLACE FUNCTION log_cmt() RETURNS TRIGGER AS $log_cmt$
		BEGIN
			INSERT INTO import_logs (kind, date) VALUES (1, CURRENT_DATE)
			ON CONFLICT (kind) DO UPDATE SET date = CURRENT_DATE;
			RETURN NULL;
		END;
	$log_cmt$ LANGUAGE plpgsql;`, // 42
	`DROP TRIGGER IF EXISTS cmt_stamp ON commitment;`, // 43
	`CREATE TRIGGER cmt_stamp AFTER INSERT OR UPDATE ON commitment
	FOR EACH STATEMENT EXECUTE FUNCTION log_cmt();`, // 44
	`CREATE OR REPLACE FUNCTION log_pmt() RETURNS TRIGGER AS $log_cmt$
		BEGIN
			INSERT INTO import_logs (kind, date) VALUES (2, CURRENT_DATE)
			ON CONFLICT (kind) DO UPDATE SET date = CURRENT_DATE;
			RETURN NULL;
		END;
	$log_cmt$ LANGUAGE plpgsql;`,
	`DROP TRIGGER IF EXISTS pmt_stamp ON payment;`,
	`CREATE TRIGGER pmt_stamp AFTER INSERT OR UPDATE ON payment
	FOR EACH STATEMENT EXECUTE FUNCTION log_pmt();`,
}

// createTablesAndViews launches the queries against the database to create all
// tables or replace the views
func createTablesAndViews(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Transaction begin %v", err)
	}
	for i, q := range initQueries {
		if _, err = tx.Exec(q); err != nil {
			tx.Rollback()
			return fmt.Errorf("Query %d %v", i, err)
		}
	}
	tx.Commit()
	return nil
}

// createSuperAdmin check if the users table creates a super admin user if not exists
func createSuperAdmin(db *sql.DB, cfg *PreLoRuGoConf, app *iris.Application) error {
	pwd := cfg.Users.SuperAdmin.Password
	email := cfg.Users.SuperAdmin.Email

	if pwd == "" || email == "" {
		return fmt.Errorf("Impossible de récupérer les credentials super admin")
	}
	var count int64
	if err := db.QueryRow("SELECT count(1) FROM users WHERE email=$1", email).
		Scan(&count); err != nil {
		return fmt.Errorf("Requête vérification super admin %v", err)
	}
	if count > 0 {
		app.Logger().Infof("Super admin déjà présent dans la base de données")
		return nil
	}
	usr := models.User{
		Name:     "Super administrateur",
		Email:    email,
		Password: pwd,
		Rights:   models.SuperAdminBit | models.ActiveAdminMask,
	}

	if err := usr.CryptPwd(); err != nil {
		return fmt.Errorf("Codage du mot de passe super admin %v", err)
	}
	if err := usr.Create(db); err != nil {
		return fmt.Errorf("Création du super admin %v", err)
	}
	app.Logger().Infof("Super admin créé")
	return nil
}

// InitDatabase connect to database, create tables and view and launch migrations
func InitDatabase(cfg *PreLoRuGoConf, app *iris.Application, dropTables bool, migrate bool) (*sql.DB, error) {
	var dbCfg *DBConf
	switch cfg.App.Stage {
	case ProductionStage:
		dbCfg = &cfg.Databases.Prod
	case DevelopmentStage:
		dbCfg = &cfg.Databases.Development
	case TestStage:
		dbCfg = &cfg.Databases.Test
	}
	cfgStr := fmt.Sprintf("sslmode=disable host=%s port=%s user=%s dbname=%s password=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.UserName, dbCfg.Name, dbCfg.Password)
	db, err := sql.Open("postgres", cfgStr)
	if err != nil {
		return nil, fmt.Errorf("Database open %v", err)
	}
	if dropTables == true {
		if err = dropAllTables(db, app); err != nil {
			return nil, err
		}
	}
	if migrate {
		if err = handleMigrations(db); err != nil {
			return nil, fmt.Errorf("Migrations %v", err)
		}
	}
	if err = createTablesAndViews(db); err != nil {
		return nil, err
	}
	if err = createSuperAdmin(db, cfg, app); err != nil {
		return nil, err
	}
	return db, nil
}
