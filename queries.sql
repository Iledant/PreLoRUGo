-- Suppression des tables
DROP TABLE IF EXISTS copro, users, imported_commitment, 
	commitment, imported_payment, payment, report, budget_action, beneficiary, 
	temp_copro, renew_project, temp_renew_project, housing, temp_housing, commitment , 
	temp_commitment, beneficiary, payment , temp_payment, action, budget_sector, commission, 
	community , temp_community, city , temp_city, renew_project_forecast , 
	temp_renew_project_forecast, copro_forecast, temp_copro_forecast;

-- Création des tables
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name varchar(50) NOT NULL,
  email varchar(120) NOT NULL,
  password varchar(120) NOT NULL,
  rights int NOT NULL);
		 
CREATE TABLE copro (
  id SERIAL PRIMARY KEY,
  reference varchar(15) NOT NULL,
  name varchar(150) NOT NULL,
  address varchar(200) NOT NULL,
  zip_code int NOT NULL,
  label_date date,
  budget bigint);
  
CREATE TABLE temp_copro (
  reference varchar(150) NOT NULL,
  name varchar(150) NOT NULL,
  address varchar(200) NOT NULL,
  zip_code int NOT NULL,
  label_date date,
  budget bigint);
  
CREATE TABLE budget_sector (
  id SERIAL PRIMARY KEY,
  name varchar(20) NOT NULL,
  full_name varchar(150));
  
CREATE TABLE budget_action (
  id SERIAL PRIMARY KEY,
  code bigint NOT NULL,
  name varchar(250) NOT NULL,
  sector_id int,
  FOREIGN KEY (sector_id) REFERENCES budget_sector(id) 
  ON UPDATE NO ACTION ON DELETE NO ACTION
  );
  
CREATE TABLE renew_project (
  id SERIAL PRIMARY KEY,
  reference varchar(15) NOT NULL UNIQUE,
  name varchar(150) NOT NULL,
  budget bigint NOT NULL,
  population int,
  composite_index int);
  
CREATE TABLE temp_renew_project (
  reference varchar(15) NOT NULL UNIQUE,
  name varchar(150) NOT NULL,
  budget bigint NOT NULL,	
  population int,
  composite_index int);
  
CREATE TABLE housing (
  id SERIAL PRIMARY KEY,
  reference varchar(100) NOT NULL,
  address varchar(150),
  zip_code int,
  plai int NOT NULL,
  plus int NOT NULL,
  pls int NOT NULL,
  anru boolean NOT NULL);
  
CREATE TABLE temp_housing (
  reference varchar(100) NOT NULL,
  address varchar(150),
  zip_code int,
  plai int NOT NULL,
  plus int NOT NULL,
  pls int NOT NULL,
  anru boolean NOT NULL);
  
CREATE TABLE beneficiary (
  id SERIAL PRIMARY KEY,
  code int NOT NULL,
  name varchar(120) NOT NULL);
  
CREATE TABLE commitment (
  id SERIAL PRIMARY KEY,
  year int NOT NULL,
  code varchar(5) NOT NULL,
  number int NOT NULL,
  line int NOT NULL,
  creation_date date NOT NULL,
  modification_date date NOT NULL,
  name varchar(150) NOT NULL,
  value bigint NOT NULL,
  beneficiary_id int NOT NULL,
  iris_code varchar(20),
  action_id int,
  housing_id int,
  copro_id int,
  renew_project_id int,
  FOREIGN KEY (beneficiary_id) REFERENCES beneficiary(id) 
  MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
  FOREIGN KEY (housing_id) REFERENCES housing(id) 
  MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
  FOREIGN KEY (copro_id) REFERENCES copro(id) 
  MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
  FOREIGN KEY (renew_project_id) REFERENCES renew_project(id) 
  MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
  FOREIGN KEY (action_id) REFERENCES budget_action(id) 
  MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION);
  
CREATE TABLE temp_commitment (
  year int NOT NULL,
  code varchar(5) NOT NULL,
  number int NOT NULL,
  line int NOT NULL,
  creation_date date NOT NULL,
  modification_date date NOT NULL,
  name varchar(150) NOT NULL,
  value bigint NOT NULL,
  beneficiary_code int NOT NULL,
  beneficiary_name varchar(150) NOT NULL,
  iris_code varchar(20),
  sector varchar(5) NOT NULL,
  action_code bigint,
  action_name varchar(150));
  
CREATE TABLE payment (
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
  FOREIGN KEY (commitment_id) REFERENCES commitment(id) MATCH SIMPLE
  ON UPDATE NO ACTION ON DELETE NO ACTION);
      
CREATE TABLE temp_payment (
  commitment_year int NOT NULL,
  commitment_code varchar(5) NOT NULL,
  commitment_number int NOT NULL,
  commitment_line int NOT NULL,
  year int NOT NULL,
  creation_date date NOT NULL,
  modification_date date NOT NULL,
  number int NOT NULL,
  value bigint NOT NULL);
  
CREATE TABLE commission (
  id SERIAL PRIMARY KEY,
  name varchar(140) NOT NULL,
  date date);
  
CREATE TABLE community (
  id SERIAL PRIMARY KEY,
  code varchar(15) NOT NULL,
  name varchar(150) NOT NULL);
  
CREATE TABLE temp_community (
  code varchar(15) NOT NULL,
  name varchar(150) NOT NULL);
  
CREATE TABLE city (
  insee_code int NOT NULL PRIMARY KEY,
  name varchar(50) NOT NULL,
  community_id int,
  qpv boolean NOT NULL,
  FOREIGN KEY (community_id) REFERENCES community(id) MATCH SIMPLE
  ON UPDATE NO ACTION ON DELETE NO ACTION);
  
CREATE TABLE temp_city (
  insee_code int NOT NULL UNIQUE,
  name varchar(50) NOT NULL,
  community_code varchar(15),
  qpv boolean NOT NULL);
  
CREATE TABLE renew_project_forecast (
  id SERIAL PRIMARY KEY,
  commission_id int NOT NULL,
  value bigint NOT NULL,
  comment text,
  renew_project_id int NOT NULL,
  FOREIGN KEY (renew_project_id) REFERENCES renew_project(id) MATCH SIMPLE
  ON UPDATE NO ACTION ON DELETE NO ACTION);
  
CREATE TABLE temp_renew_project_forecast (
  id int NOT NULL,
  commission_id int NOT NULL,
  value bigint NOT NULL,
  comment text,
  renew_project_id int NOT NULL);
  
CREATE TABLE copro_forecast (
  id SERIAL PRIMARY KEY,
  commission_id int NOT NULL,
  value bigint NOT NULL,
  comment text,
  copro_id int NOT NULL,
  FOREIGN KEY (copro_id) REFERENCES copro (id) MATCH SIMPLE
  ON UPDATE NO ACTION ON DELETE NO ACTION);
  
CREATE TABLE temp_copro_forecast (
  id int NOT NULL,
  commission_id int NOT NULL,
  value bigint NOT NULL,
  comment text,
  copro_id int NOT NULL);
		 
-- Insertion des utilisateurs
INSERT INTO users (name,email,password,rights) VALUES
('Administrateur','cs@if.fr','$2a$10$bQ1G2K8UeH8mTwERhwe2XerDeQtVLN02GDz4HD4WP/N9X/7S.MhbO',5),
-- cSpell: disable
('Utilisateur','user@if.fr','$2a$10$tMrZWq5yIgPI8tBwFge/B.aZ.4FyahEhd21Qdgwfc9TYCZhQbILAC',1);
-- cSpell: enable

-- Creation de la VIEW pour avoir des engagements cumulés
CREATE VIEW cumulated_commitment AS
  SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,q.value,
    c.beneficiary_id, c.iris_code,c.action_id,c.housing_id, c.copro_id,
    c.renew_project_id
  FROM commitment c
  JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
    min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
  ON c.id = q.id;

-- Mise à jour de la table payment pour que payment pointe vers les id de
-- la view cumulated_payment
UPDATE payment SET commitment_id = cumulated_commitment.id
  FROM cumulated_commitment
  WHERE payment.commitment_year = cumulated_commitment.year 
    AND payment.commitment_code = cumulated_commitment.code
    AND payment.commitment_number = cumulated_commitment.number;

-- Ajout de la colonne QPV dans la table housing
ALTER TABLE housing ADD COLUMN qpv boolean NOT NULL DEFAULT false;

ALTER TABLE temp_housing ADD COLUMN qpv boolean NOT NULL; 

-- Recopie des données d'engagements et de paiement pour alimenter la base de données
DELETE from payment;
DELETE from commitment;

COPY temp_commitment(year, code, number, line, creation_date, modification_date, 
  name, value, beneficiary_code, beneficiary_name, iris_code, sector, action_code, 
  action_name) FROM './assets/20190129 AP PreLoRu.csv' DELIMITER ';' CSV HEADER;

UPDATE temp_commitment set code=RTRIM(code), name=RTRIM(name), 
  beneficiary_name=RTRIM(beneficiary_name), sector=RTRIM(sector), 
  action_name=RTRIM(action_name);

COPY temp_payment(commitment_year, commitment_code, commitment_number, 
  commitment_line, year, number, creation_date, modification_date, value)
  FROM './assets/20190123 Mandats PreLoRU.csv' DELIMITER ';' CSV HEADER;

UPDATE temp_payment set commitment_code = RTRIM(commitment_code);

INSERT INTO beneficiary (code,name) SELECT DISTINCT beneficiary_code,beneficiary_name 
		FROM temp_commitment WHERE beneficiary_code not in (SELECT code from beneficiary);
		
INSERT INTO budget_sector (name) SELECT DISTINCT sector
	FROM temp_commitment WHERE sector not in (SELECT name from budget_sector);
	
INSERT INTO budget_action (code,name,sector_id) 
		SELECT DISTINCT ic.action_code,ic.action_name, s.id
		FROM temp_commitment ic
		LEFT JOIN budget_sector s ON ic.sector = s.name
		WHERE action_code not in (SELECT code from budget_action);
		

INSERT INTO commitment (year,code,number,line,creation_date,modification_date,
		name,value,beneficiary_id,iris_code,action_id)
  	(SELECT ic.year,ic.code,ic.number,ic.line,ic.creation_date,ic.modification_date,
			ic.name,ic.value,b.id,ic.iris_code,a.id
  	FROM temp_commitment ic
		JOIN beneficiary b on ic.beneficiary_code=b.code
		LEFT JOIN budget_action a on ic.action_code = a.code
  	WHERE (ic.year,ic.code,ic.number,ic.line,ic.creation_date,ic.modification_date,ic.name, ic.value) 
    NOT IN (select year,code,number,line,creation_date,modification_date,name,value FROM commitment));
	
select count(1) from temp_payment

INSERT INTO payment (commitment_id,commitment_year,commitment_code,
			commitment_number,commitment_line,year,creation_date,modification_date,
			number, value)
		SELECT c.id,t.commitment_year,t.commitment_code,t.commitment_number,
			t.commitment_line,t.year,t.creation_date,t.modification_date,t.number,t.value 
			FROM temp_payment t
			LEFT JOIN cumulated_commitment c 
				ON t.commitment_year=c.year AND t.commitment_code=c.code
				AND t.commitment_number=c.number
			WHERE (t.commitment_year,t.commitment_code,t.commitment_number,
				t.commitment_line,t.year,t.creation_date,t.modification_date) 
			NOT IN (SELECT DISTINCT commitment_year,commitment_code,commitment_number,
				commitment_line,year,creation_date,modification_date from payment);

DELETE from temp_payment;
DELETE from temp_commitment;
