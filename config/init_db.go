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

var initQueries = []string{`CREATE EXTENSION IF NOT EXISTS tablefunc`, // 0 tablefunc
	`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name varchar(50) NOT NULL,
		email varchar(120) NOT NULL,
		password varchar(120) NOT NULL,
		rights int NOT NULL
		);`, // 1 : users
	`CREATE TABLE IF NOT EXISTS department (
			id SERIAL PRIMARY KEY,
			code int NOT NULL,
			name varchar(20)
		);`, // 2 department
	`CREATE TABLE IF NOT EXISTS community (
	    id SERIAL PRIMARY KEY,
	    code varchar(15) NOT NULL,
			name varchar(150) NOT NULL,
			department_id int,
			FOREIGN KEY (department_id) REFERENCES department (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 3 : community
	`CREATE TABLE IF NOT EXISTS temp_community (
	    code varchar(15) NOT NULL,
			name varchar(150) NOT NULL,
			department_code int
		);`, // 4 : temp_community
	`CREATE TABLE IF NOT EXISTS city (
	    insee_code int NOT NULL PRIMARY KEY,
	    name varchar(50) NOT NULL,
			community_id int,
			qpv boolean NOT NULL,
			FOREIGN KEY (community_id) REFERENCES community (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 5 : city
	`CREATE TABLE IF NOT EXISTS temp_city (
	    insee_code int NOT NULL UNIQUE,
	    name varchar(50) NOT NULL,
	    community_code varchar(15),
			qpv boolean NOT NULL
		);`, // 6 : temp_city
	`CREATE TABLE IF NOT EXISTS copro (
			id SERIAL PRIMARY KEY,
			reference varchar(60) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200),
			zip_code int,
			label_date date,
			budget bigint,
			FOREIGN KEY (zip_code) REFERENCES city (insee_code) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 7 : copro
	`CREATE TABLE IF NOT EXISTS temp_copro (
			reference varchar(150) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint
		);`, // 8 : temp_copro
	`CREATE TABLE IF NOT EXISTS budget_sector (
			id SERIAL PRIMARY KEY,
			name varchar(20) NOT NULL,
	    full_name varchar(150)
		);`, // 9 : budget_sector
	`CREATE TABLE IF NOT EXISTS budget_action (
			id SERIAL PRIMARY KEY,
			code bigint NOT NULL,
			name varchar(250) NOT NULL,
			sector_id int
		);`, // 10 : budget_action
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
		);`, // 11 : renew_project
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
		);`, // 12 : temp_renew_project
	`CREATE TABLE IF NOT EXISTS housing_type (
			id SERIAL PRIMARY KEY,
			short_name varchar(10) NOT NULL UNIQUE,
			long_name varchar(100)
		)`, // 13 housing_type
	`CREATE TABLE IF NOT EXISTS housing (
	    id SERIAL PRIMARY KEY,
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
			anru boolean NOT NULL,
			housing_type_id int,
			FOREIGN KEY (zip_code) REFERENCES city (insee_code) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (housing_type_id) REFERENCES housing_type(id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 14 : housing
	`CREATE TABLE IF NOT EXISTS temp_housing (
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
	    anru boolean NOT NULL
		);`, // 15 : temp_housing
	`CREATE TABLE IF NOT EXISTS beneficiary (
	    id SERIAL PRIMARY KEY,
	    code int NOT NULL UNIQUE,
	    name varchar(120) NOT NULL
		);`, // 16 : beneficiary
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
		);`, // 17 : commitment
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
		);`, // 18 : temp_commitment
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
			receipt_date date,
			FOREIGN KEY (commitment_id) REFERENCES commitment (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 19 : payment
	`CREATE TABLE IF NOT EXISTS temp_payment (
	    commitment_year int NOT NULL,
	    commitment_code varchar(5) NOT NULL,
	    commitment_number int NOT NULL,
	    commitment_line int NOT NULL,
	    year int NOT NULL,
	    creation_date date NOT NULL,
			modification_date date NOT NULL,
			number int NOT NULL,
			value bigint NOT NULL,
			receipt_date date
		);`, // 20 : temp_payment
	`CREATE TABLE IF NOT EXISTS commission (
	    id SERIAL PRIMARY KEY,
	    name varchar(140) NOT NULL,
	    date date
		);`, // 21 : commission
	`CREATE TABLE IF NOT EXISTS renew_project_forecast (
	    id SERIAL PRIMARY KEY,
	    commission_id int NOT NULL,
			value bigint NOT NULL,
			project varchar(150),
	    comment text,
			renew_project_id int NOT NULL,
			action_id int NOT NULL,
			FOREIGN KEY (renew_project_id) REFERENCES renew_project (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 22 : renew_project_forecast
	`CREATE TABLE IF NOT EXISTS temp_renew_project_forecast (
			id int NOT NULL,
			commission_id int NOT NULL,
	    value bigint NOT NULL,
			project varchar(150),
	    comment text,
			renew_project_id int NOT NULL,
			action_code bigint NOT NULL
		);`, // 23 : temp_renew_project_forecast
	`CREATE TABLE IF NOT EXISTS copro_forecast (
	    id SERIAL PRIMARY KEY,
	    commission_id int NOT NULL,
			value bigint NOT NULL,
			project varchar(150),
	    comment text,
			copro_id int NOT NULL,
			action_id int NOT NULL,
			FOREIGN KEY (copro_id) REFERENCES copro (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
			FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
		);`, // 24 : copro_forecast
	`CREATE TABLE IF NOT EXISTS temp_copro_forecast (
			id int NOT NULL,
			commission_id int NOT NULL,
	    value bigint NOT NULL,
			project varchar(150),
	    comment text,
			copro_id int NOT NULL,
			action_code bigint NOT NULL
		);`, // 25 : temp_copro_forecast
	`CREATE OR REPLACE VIEW cumulated_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,
		  q.value,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id, c.copro_id,
			c.renew_project_id,c.caducity_date
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`, // 26 : cumulated_commitment view
	`CREATE OR REPLACE VIEW cumulated_sold_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,
			q.value, c.sold_out,c.beneficiary_id, c.iris_code,c.action_id,c.housing_id,
			c.copro_id,c.renew_project_id,c.caducity_date
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`, // 27 : cumulated_sold_commitment view
	`CREATE TABLE IF NOT EXISTS ratio (
			id SERIAL PRIMARY KEY,
			year int NOT NULL,
			sector_id int NOT NULL,
			index int NOT NULL,
			ratio double precision NOT NULL,
			FOREIGN KEY (sector_id) REFERENCES budget_sector (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
			);`, // 28 : ratio
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
			);`, // 29 : renew_project_forecast
	`CREATE TABLE IF NOT EXISTS temp_housing_forecast (
				id int NOT NULL,
				commission_id int NOT NULL,
				value bigint NOT NULL,
				comment text,
				action_id int NOT NULL
			);`, // 30 : temp_renew_project_forecast
	`CREATE TABLE IF NOT EXISTS housing_commitment (
			iris_code varchar(20),
	    reference varchar(100)
			);`, // 31 housing_commitment
	`CREATE TABLE IF NOT EXISTS copro_commitment (
			iris_code varchar(20),
	    reference varchar(100)
			);`, // 32 copro_commitment
	`CREATE TABLE IF NOT EXISTS migration (
		id SERIAL PRIMARY KEY,
		created timestamp NOT NULL,
		index int NOT NULL,
		query text
	);`, // 33 migration
	`CREATE TABLE IF NOT EXISTS rp_event_type (
		id SERIAL PRIMARY KEY,
		name varchar(100) NOT NULL
	);`, // 34 rp_event_type
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
	);`, // 35 rp_event
	`CREATE TABLE IF NOT EXISTS rp_cmt_city_join (
		id SERIAL PRIMARY KEY,
		commitment_id int NOT NULL UNIQUE,
		city_code int NOT NULL,
		FOREIGN KEY (commitment_id) REFERENCES commitment (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
		FOREIGN KEY (city_code) REFERENCES city (insee_code) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 36 rp_cmt_city_join
	`CREATE TABLE IF NOT EXISTS pre_prog (
		id SERIAL PRIMARY KEY,
		year int NOT NULL,
		commission_id int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		project varchar(150),
		comment text,
		action_id int,
		FOREIGN KEY (commission_id) REFERENCES commission (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE,
		FOREIGN KEY (action_id) REFERENCES budget_action (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 37 pre_prog
	`CREATE TABLE IF NOT EXISTS temp_pre_prog (
		commission_id int NOT NULL,
		year int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		project varchar(150),
		comment text,
		action_id int
	);`, // 38 temp_pre_prog
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
	);`, // 39 prog
	`CREATE TABLE IF NOT EXISTS temp_prog (
		commission_id int NOT NULL,
		year int NOT NULL,
		value bigint NOT NULL,
		kind int CHECK (kind IN (1,2,3)),
		kind_id int,
		comment text,
		action_id int
	);`, // 40 temp_prog
	`CREATE TABLE IF NOT EXISTS rpls(
		id SERIAL PRIMARY KEY,
		insee_code int NOT NULL,
		year int NOT NULL,
		ratio double precision NOT NULL,
		FOREIGN KEY (insee_code) REFERENCES city (insee_code) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION DEFERRABLE
	);`, // 41 rpls		`
	`CREATE TABLE IF NOT EXISTS temp_rpls(
		insee_code int NOT NULL,
		year int NOT NULL,
		ratio double precision NOT NULL
	);`, // 42 temp_rpls
	`CREATE TABLE IF NOT EXISTS import_logs(
		kind int UNIQUE,
		date date
	);`, // 43 import_logs
	`CREATE TABLE IF NOT EXISTS temp_housing_summary(
		reference_code varchar(150),
		address varchar(150),
		iris_code varchar(20),
		pls int,
		plai int,
		plus int,
		anru boolean,
		insee_code int REFERENCES city(insee_code)
	);`, // 44 temp_housing_summary
	`CREATE TABLE IF NOT EXISTS housing_summary(
		id SERIAL PRIMARY KEY,
		year int NOT NULL,
		housing_ref varchar(100) NOT NULL,
		import_ref varchar(150) NOT NULL,
		iris_code varchar(20) NOT NULL
	);`, // 45 housing_summary
	`CREATE TABLE IF NOT EXISTS copro_event_type (
		id SERIAL PRIMARY KEY,
		name varchar(100) NOT NULL
	);`, // 46 copro_event_type
	`CREATE TABLE IF NOT EXISTS copro_event (
		id SERIAL PRIMARY KEY,
		copro_id int NOT NULL REFERENCES copro(id),
		copro_event_type_id int NOT NULL REFERENCES copro_event_type(id),
		date date NOT NULL,
		comment text
	);`, // 47 copro_event
	`CREATE TABLE IF NOT EXISTS copro_doc (
		id SERIAL PRIMARY KEY,
		copro_id int NOT NULL references copro(id),
		name varchar(150) NOT NULL,
		link varchar(250) NOT NULL
	);`, // 48 copro_doc
	`CREATE TABLE IF NOT EXISTS payment_credit (
		id SERIAL PRIMARY KEY,
		chapter int NOT NULL,
		function int NOT NULL,
		primitive bigint NOT NULL,
		reported bigint NOT NULL,
		added bigint NOT NULL,
		modified bigint NOT NULL,
		movement bigint NOT NULL,
		year int NOT NULL
	);`, // 49 payment_credit
	`CREATE TABLE IF NOT EXISTS payment_credit_journal (
		id SERIAL PRIMARY KEY,
		chapter int NOT NULL,
		function int NOT NULL,
		creation_date date NOT NULL,
		modification_date date NOT NULL,
		name varchar(150) NOT NULL,
		value bigint NOT NULL
	);`, // 50 payment_credit_journal
	`CREATE TABLE IF NOT EXISTS home_message (
		title varchar(255),
		body text
	);`, // 51 home_message
	`CREATE TABLE IF NOT EXISTS placement (
		id SERIAL PRIMARY KEY,
		iris_code varchar(20) NOT NULL UNIQUE,
		count int,
		contract_year int,
		comment varchar(150),
		commitment_id int REFERENCES commitment(id)
	);`, // 52 placement
	`CREATE TABLE IF NOT EXISTS temp_placement (
		iris_code varchar(20) NOT NULL,
		count int,
		contract_year int
	);`, // 53 temp_placement
	`CREATE TABLE IF NOT EXISTS beneficiary_group (
		id SERIAL PRIMARY KEY,
		name varchar(150) NOT NULL UNIQUE
	)`, // 54 beneficiary_group
	`CREATE TABLE IF NOT EXISTS beneficiary_belong (
		id SERIAL PRIMARY KEY,
		beneficiary_id int REFERENCES beneficiary(id) ON DELETE CASCADE,
		group_id int REFERENCES beneficiary_group(id) ON DELETE CASCADE
	)`, // 55 beneficiary_belong
	`CREATE OR REPLACE FUNCTION log_cmt() RETURNS TRIGGER AS $log_cmt$
		BEGIN
			INSERT INTO import_logs (kind, date) VALUES (1, CURRENT_DATE)
			ON CONFLICT (kind) DO UPDATE SET date = CURRENT_DATE;
			RETURN NULL;
		END;
	$log_cmt$ LANGUAGE plpgsql;`, // 56
	`DROP TRIGGER IF EXISTS cmt_stamp ON commitment;`, // 57
	`CREATE TRIGGER cmt_stamp AFTER INSERT OR UPDATE ON commitment
	FOR EACH STATEMENT EXECUTE FUNCTION log_cmt();`, // 58
	`CREATE OR REPLACE FUNCTION log_pmt() RETURNS TRIGGER AS $log_cmt$
		BEGIN
			INSERT INTO import_logs (kind, date) VALUES (2, CURRENT_DATE)
			ON CONFLICT (kind) DO UPDATE SET date = CURRENT_DATE;
			RETURN NULL;
		END;
	$log_cmt$ LANGUAGE plpgsql;`, // 59
	`DROP TRIGGER IF EXISTS pmt_stamp ON payment;`, // 60
	`CREATE TRIGGER pmt_stamp AFTER INSERT OR UPDATE ON payment
	FOR EACH STATEMENT EXECUTE FUNCTION log_pmt();`, // 61
	`CREATE TABLE IF NOT EXISTS housing_typology (
		id SERIAL PRIMARY KEY,
		name varchar(30) UNIQUE
	)`, // 62 housing_typology
	`CREATE TABLE IF NOT EXISTS housing_convention (
		id SERIAL PRIMARY KEY,
		name varchar(30) UNIQUE
	)`, // 63 housing_convention
	`CREATE TABLE IF NOT EXISTS convention_type (
		id SERIAL PRIMARY KEY,
		name varchar(15) UNIQUE
	)`, // 64 convention_type
	`CREATE TABLE IF NOT EXISTS housing_transfer (
		id SERIAL PRIMARY KEY,
		name varchar(50) UNIQUE
	)`, // 65 housing_tranfer
	`CREATE TABLE IF NOT EXISTS housing_comment (
		id SERIAL PRIMARY KEY,
		name varchar(150) UNIQUE
	)`, // 66 housing_comment
	`CREATE TABLE IF NOT EXISTS reservation_fee (
		id SERIAL PRIMARY KEY,
		current_beneficiary_id int NOT NULL REFERENCES beneficiary(id),
		first_beneficiary_id int REFERENCES beneficiary(id),
		city_code int REFERENCES city(insee_code),
		address_number varchar(20),
		address_street varchar(100),
		rpls varchar(15),
		convention varchar(80),
		convention_type_id int REFERENCES convention_type(id),
		count int,
		transfer_date date,
		transfer_id int REFERENCES housing_transfer(id),
		pmr boolean,
		comment_id int REFERENCES housing_comment(id),
		convention_date date,
		elise_ref varchar(30),
		area double precision,
		end_year int,
		loan double precision,
		charges double precision,
		typology_id int REFERENCES housing_typology(id)
	)`, // 67 reservation_fee
	`CREATE TABLE IF NOT EXISTS temp_reservation_fee (
		current_beneficiary varchar(80),
		first_beneficiary varchar(80),
		city varchar(50),
		address_number varchar(20),
		address_street varchar(100),
		convention varchar(80),
		typology varchar(30),
		rpls varchar(15),
		convention_type varchar(15),
		count int,
		transfer varchar(50),
		transfer_date date,
		pmr boolean,
		comment varchar(150),
		convention_date date,
		area double precision,
		end_year int,
		loan double precision,
		charges double precision
	)`, // 68 temp_reservation_fee
	`CREATE TABLE IF NOT EXISTS temp_iris_housing_type(
		iris_code varchar(20),
		housing_type_short_name varchar(10)
	)`, // 69 temp_iris_housing_type
	`CREATE TABLE IF NOT EXISTS reservation_report(
		id SERIAL PRIMARY KEY,
		beneficiary_id int NOT NULL REFERENCES beneficiary(id),
		area double precision NOT NULL,
		source_iris_code varchar(20) NOT NULL,
		dest_iris_code varchar(20),
		dest_date date
	)`, // 70 reservation_report
	`CREATE TABLE IF NOT EXISTS temp_payment_demands (
		iris_code varchar(32) NOT NULL,
		iris_name varchar(200) NOT NULL,
		commitment_date date NOT NULL,
		beneficiary_code int NOT NULL,
		demand_number int NOT NULL,
		demand_date	date NOT NULL,
		receipt_date date NOT NULL,
		demand_value bigint NOT NULL,
		csf_date date,
		csf_comment text,
		demand_status varchar(15),
		status_comment text
	)`, // 71 temp_payment_demands
	`CREATE OR REPLACE VIEW imported_payment_demands AS
		SELECT iris_code,iris_name,MAX(commitment_date),beneficiary_code,
			demand_number,demand_date,receipt_date,demand_value,csf_date,csf_comment,
			demand_status,status_comment FROM temp_payment_demands
			GROUP BY 1,2,4,5,6,7,8,9,10,11,12`, // 72 imported_payment_demands
	`CREATE TABLE IF NOT EXISTS payment_demands (
		id SERIAL PRIMARY KEY,
		import_date date NOT NULL,
		iris_code varchar(32) NOT NULL,
		iris_name varchar(200) NOT NULL,
		beneficiary_id int NOT NULL REFERENCES beneficiary(id),
		demand_number int NOT NULL,
		demand_date	date NOT NULL,
		receipt_date date NOT NULL,
		demand_value bigint NOT NULL,
		csf_date date,
		csf_comment text,
		demand_status varchar(15),
		status_comment text,
		excluded boolean,
		excluded_comment varchar(150),
		processed_date date
	)`, // 73 payment_demands
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
	if err = createTablesAndViews(db); err != nil {
		return nil, err
	}
	if migrate {
		if err = handleMigrations(db); err != nil {
			return nil, fmt.Errorf("Migrations %v", err)
		}
	}
	if err = createSuperAdmin(db, cfg, app); err != nil {
		return nil, err
	}
	return db, nil
}
