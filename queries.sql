CREATE TABLE imported_commitment (
	report varchar(12),
	"action" varchar(255),
	iris_code varchar(50),
	coriolis_year integer,
	coriolis_egt_code varchar(30),
	coriolis_egt_num varchar(8),
	coriolis_egt_line varchar(3),
	name varchar(200),
	beneficiary varchar(120),
	beneficiary_code int,
	"date" date, 
	"value" bigint,
	lapse_date date
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

-- Création des actions budgétaires à partir de l'import des engagements
WITH 
 ba AS (SELECT DISTINCT substring(action from '^[0-9]+') AS code, ltrim(substring(action from '- .+'),'- ') AS name 
		FROM imported_commitment)
INSERT INTO budget_action (code,name) SELECT ba.code,ba.name FROM ba WHERE ba.code NOT IN (SELECT code FROM budget_action)

-- Création des bénéficiaires à partir de l'import des engagements
INSERT INTO beneficiary (code,name) 
  SELECT DISTINCT beneficiary_code,beneficiary FROM imported_commitment 
    WHERE beneficiary_code NOT IN (SELECT DISTINCT code from beneficiary)

-- Création des imports à partir de l'import des engagements	
INSERT INTO report (name) SELECT DISTINCT report from imported_commitment WHERE report NOT IN (SELECT DISTINCT name FROM report)

-- Requête de mise à jour de la table commitment à partir des imports
-- Mise à jour des lignes communes puis ajout des lignes non présentes
UPDATE commitment SET action_id=q.ba_id, value=q.value, name=q.name, report_id=q.re_id,
lapse_date=q.lapse_date, iris_code=q.iris_code, lapse_date=q.lapse_date, beneficiary_id=q.be_id FROM
(SELECT ic.date, ic.value, ic.name, ic.lapse_date,ic.iris_code,ic.coriolis_year,
ic.coriolis_egt_code, ic.coriolis_egt_num,ic.coriolis_egt_line, ba.id AS ba_id, be.id AS be_id, re.id AS re_id
FROM imported_commitment ic, budget_action ba, report re, beneficiary be
WHERE ic.action LIKE ba.code||' -%'  AND re.name = ic.report AND be.code = ic.beneficiary_code) q
WHERE commitment.coriolis_year=q.coriolis_year AND commitment.coriolis_egt_code=q.coriolis_egt_code AND 
commitment.coriolis_egt_line=q.coriolis_egt_line AND commitment.coriolis_egt_num=q.coriolis_egt_num;

INSERT INTO commitment (date, value, name, lapse_date, iris_code,coriolis_year, 
coriolis_egt_code, coriolis_egt_num,coriolis_egt_line, action_id, beneficiary_id, report_id)
SELECT ic.date, ic.value, ic.name, ic.lapse_date,ic.iris_code,ic.coriolis_year,
ic.coriolis_egt_code, ic.coriolis_egt_num,ic.coriolis_egt_line, ba.id, be.id, re.id
FROM imported_commitment ic, budget_action ba, report re, beneficiary be
WHERE ic.action LIKE ba.code||' -%'  AND re.name = ic.report AND be.code = ic.beneficiary_code;

-- Lecture du fichier des paiements
COPY imported_payment	FROM E'C:\\Users\\chris\\go\\src\\github.com\\Iledant\\PreLoRUGo\\assets\\20190123 Mandats PreLoRU.csv' DELIMITER ';' CSV HEADER;

-- Requête de mise à jour de la table payment à parti des imports
-- Mise à jour des lignes communes puis ajout des lignes non présentes

INSERT INTO payment (coriolis_year,coriolis_egt_code,coriolis_egt_num,coriolis_egt_line,commitment_id,
date,number,value,cancelled_value,beneficiary_code,beneficiary_id)
SELECT ip.coriolis_year, ip.coriolis_egt_code, ip.coriolis_egt_num, ip.coriolis_egt_line, co.id, 
  ip.date, ip.number,ip.value,ip.cancelled_value,ip.beneficiary_code,be.id
FROM imported_payment ip
LEFT OUTER JOIN commitment co ON ip.coriolis_year = co.coriolis_year, ip.coriolis_egt_code = co.coriolis_egt_code,
  ip.coriolis_egt_line = co.coriolis_egt_line, ip.coriolis_egt_num = co.coriolis_egt_num
LEFT OUTER JOIN beneficiary be ON ip.beneficiary_code = be.code