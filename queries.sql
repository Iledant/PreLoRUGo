-- Suppression des tables
DROP TABLE imported_commitment;

CREATE TABLE imported_commitment (
YEAR int,
CODE varchar(5),
NUM int,
LIG int,
CREATION int,
MODIFICATION int,
NAME varchar(50),
VALUE bigint,
BENEFICIARY_CODE int,
BENEFICIARY_NAME varchar(40),
IRIS_CODE varchar(20)
);

CREATE TABLE commitment (
  ID SERIAL PRIMARY KEY,
  YEAR int NOT NULL,
  CODE varchar(5) NOT NULL,
  NUM int NOT NULL,
  LIG INT NOT NULL,
  CREATION date,
  MODIFICATION date,
  NAME varchar(60),
  VALUE bigint,
  BENEFICIARY_ID int,
  IRIS_CODE varchar(20),
  CONSTRAINT commitment_beneficiary_id_fkey FOREIGN KEY (beneficiary_id)
	  REFERENCES beneficiary (id) MATCH SIMPLE
	  ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE table imported_payment (
  coriolis_year int NOT NULL,
  coriolis_egt_code varchar(30),
	coriolis_egt_num varchar(8),
	coriolis_egt_line varchar(3),
  date date,
  number int,
  value bigint,
  cancelled_value bigint,
  beneficiary_code int
);

CREATE table payment (
  id SERIAL PRIMARY KEY,
  coriolis_year int NOT NULL,
  coriolis_egt_code varchar(30) NOT NULL,
	coriolis_egt_num varchar(8) NOT NULL,
	coriolis_egt_line varchar(3) NOT NULL,
  commitment_id integer,
  date date NOT NULL,
  number int NOT NULL,
  value bigint NOT NULL,
  cancelled_value bigint NOT NULL,
  beneficiary_code int NOT NULL,
  beneficiary_id integer,
  CONSTRAINT payment_commitment_id_fkey FOREIGN KEY (commitment_id)
	  REFERENCES commitment (id) MATCH SIMPLE
	  ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT payment_beneficiary_id_fkey FOREIGN KEY (beneficiary_id)
	  REFERENCES beneficiary (id) MATCH SIMPLE
	  ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE table report (
	id SERIAL PRIMARY KEY,
	name varchar(8),
	date date
);

CREATE table budget_action (
	id SERIAL PRIMARY KEY,
	code varchar(12) NOT NULL,
	name varchar(250) NOT NULL,
	sector_id int
  CONSTRAINT budget_action_sector_id_fkey FOREIGN KEY (sector_id)
	  REFERENCES sector (id) MATCH SIMPLE
	  ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE table beneficiary (
	id SERIAL PRIMARY KEY,
	code int NOT NULL,
	name varchar(120) NOT NULL
);

CREATE TABLE commitment (
id SERIAL PRIMARY KEY,
action_id integer NOT NULL,
date date NOT NULL,
value bigint NOT NULL,
name varchar(200) NOT NULL,
report_id integer NOT NULL,
lapse_date date,
	iris_code varchar(50) NOT NULL,
	coriolis_year integer NOT NULL,
	coriolis_egt_code varchar(30) NOT NULL,
	coriolis_egt_num varchar(8) NOT NULL,
	coriolis_egt_line varchar(3) NOT NULL,
beneficiary_id integer NOT NULL,
CONSTRAINT commitment_action_id_fkey FOREIGN KEY (action_id)
	REFERENCES budget_action (id) MATCH SIMPLE
	ON UPDATE NO ACTION ON DELETE NO ACTION,
CONSTRAINT commitment_beneficiary_id_fkey FOREIGN KEY (beneficiary_id)
	REFERENCES beneficiary (id) MATCH SIMPLE
	ON UPDATE NO ACTION ON DELETE NO ACTION,
CONSTRAINT commitment_report_id_fkey FOREIGN KEY (report_id)
	REFERENCES report (id) MATCH SIMPLE
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE users (
id SERIAL PRIMARY KEY,
name varchar(50) NOT NULL,
email varchar(120) NOT NULL,
password varchar(120) NOT NULL,
role varchar(15) NOT NULL,
active boolean NOT NULL
);

CREATE TABLE copros (
id SERIAL PRIMARY KEY,
reference varchar(150) NOT NULL,
name varchar(150) NOT NULL,
address varchar(200) NOT NULL,
zip_code int NOT NULL,
label_date date,
budget bigint
);

-- Import des engagements
COPY imported_commitment FROM E'C:\\Users\\chris\\go\\src\\github.com\\Iledant\\PreLoRUGo\\assets\\20190129 AP PreLoRU.csv' DELIMITER ';' CSV HEADER;

-- Import ou mise à jour des bénéficiaires à partir de l'import des engagements
insert into beneficiary (code, name) 
	select distinct beneficiary_code, beneficiary_name 
		from imported_commitment 
		where beneficiary_code not in (select code from beneficiary);
		
update beneficiary set name = t.beneficiary_name from
	(select distinct beneficiary_code, beneficiary_name
		from imported_commitment) t
	where t.beneficiary_code not in (select code from beneficiary);

-- Mise à jour de la table des engagements
INSERT INTO commitment (year,code,num,lig,creation,modification,name,value,beneficiary_id,iris_code)
  (SELECT ic.year,ic.code,ic.num,ic.lig,make_date(ic.creation/10000,(ic.creation/100)%100,ic.creation%100),
    make_date(ic.modification/10000,(ic.modification/100)%100,ic.modification%100),ic.name,ic.value,b.id,ic.iris_code
  FROM imported_commitment ic
  JOIN beneficiary b on ic.beneficiary_code=b.code
  WHERE (ic.year,ic.code,ic.num,ic.lig,make_date(ic.creation/10000,(ic.creation/100)%100,ic.creation%100),
    make_date(ic.modification/10000,(ic.modification/100)%100,ic.modification%100),ic.name, ic.value) 
    NOT IN (select year,code,num,lig,creation,modification,name,value FROM commitment));

-- Lecture du fichier des paiements
COPY imported_payment	FROM E'C:\\Users\\chris\\go\\src\\github.com\\Iledant\\PreLoRUGo\\assets\\20190123 Mandats PreLoRU.csv' DELIMITER ';' CSV HEADER;
