package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Iledant/PreLoRUGo/config"
	"github.com/Iledant/PreLoRUGo/models"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
)

// TestContext contains all items for units tests in API.
type TestContext struct {
	DB             *sql.DB
	App            *iris.Application
	E              *httpexpect.Expect
	Config         *config.PreLoRuGoConf
	CommissionID   int64
	RenewProjectID int64
	CoproID        int64
}

// TestCase is used as common structure for all request tests
type TestCase struct {
	Sent          []byte
	Token         string
	RespContains  []string
	StatusCode    int
	ID            int
	Count         int
	CountItemName string
	IDName        string
}

// TestAll embeddes all test functions and is the only test entry point
// It initializes a fresh new test database base and call test functions
// in the right order to avoid side effects
func TestAll(t *testing.T) {
	cfg := initializeTests(t)
	testUser(t, cfg)
	testCommunity(t, cfg)
	testCity(t, cfg)
	testCopro(t, cfg)
	testBudgetAction(t, cfg)
	testRenewProject(t, cfg)
	testHousing(t, cfg)
	testCommitment(t, cfg)
	testBeneficiary(t, cfg)
	testPayment(t, cfg)
	testBudgetSector(t, cfg)
	testCommitmentLink(t, cfg)
	testCommission(t, cfg)
	testRenewProjectForecast(t, cfg)
	testCoproForecast(t, cfg)
	testSettings(t, cfg)
	testHome(t, cfg)
	testBeneficiaryDatas(t, cfg)
	testPmtRatio(t, cfg)
	testPmtForecasts(t, cfg)
}

func initializeTests(t *testing.T) *TestContext {
	testCtx := &TestContext{}
	cfg := &config.PreLoRuGoConf{}
	var err error
	testCtx.App = iris.New().Configure(iris.WithConfiguration(
		iris.Configuration{DisablePathCorrection: true}))
	if err = cfg.Get(); err != nil {
		t.Error("Configuration : " + err.Error())
		t.FailNow()
		return nil
	}
	testCtx.Config = cfg
	testCtx.DB, err = config.LaunchDB(&testCtx.Config.Databases.Test)
	if err != nil {
		t.Error("Erreur de connexion à postgres : " + err.Error())
		t.FailNow()
		return nil
	}
	initializeTestDB(t, testCtx.DB, testCtx.Config)
	SetRoutes(testCtx.App, testCtx.DB)
	testCtx.E = httptest.New(t, testCtx.App)
	// Fetch admin and user tokens
	fetchTokens(t, testCtx)
	return testCtx
}

func initializeTestDB(t *testing.T, db *sql.DB, cfg *config.PreLoRuGoConf) {
	dropQueries := []string{`DROP VIEW IF EXISTS cumulated_commitment, 
	cumulated_sold_commitment`,
		`DROP TABLE IF EXISTS copro, users, imported_commitment, 
	commitment, imported_payment, payment, report, budget_action, beneficiary, 
	temp_copro, renew_project, temp_renew_project, housing, temp_housing, commitment , 
	temp_commitment, beneficiary, payment , temp_payment, action, budget_sector, commission, 
	community , temp_community, city , temp_city, renew_project_forecast , 
	temp_renew_project_forecast, copro_forecast, temp_copro_forecast, ratio`,
	}
	for i, q := range dropQueries {
		if _, err := db.Exec(q); err != nil {
			t.Errorf("Suppression table/views[%d] : %s", i, err.Error())
			t.FailNow()
			return
		}
	}
	queries := []string{`CREATE EXTENSION IF NOT EXISTS tablefunc`,
		`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name varchar(50) NOT NULL,
		email varchar(120) NOT NULL,
		password varchar(120) NOT NULL,
		rights int NOT NULL
		);`, // 0 : users
		`CREATE TABLE community (
	    id SERIAL PRIMARY KEY,
	    code varchar(15) NOT NULL,
	    name varchar(150) NOT NULL
		);`, // 1 : community
		`CREATE TABLE temp_community (
	    code varchar(15) NOT NULL,
	    name varchar(150) NOT NULL
		);`, // 2 : temp_community
		`CREATE TABLE city (
	    insee_code int NOT NULL PRIMARY KEY,
	    name varchar(50) NOT NULL,
			community_id int,
			qpv boolean NOT NULL,
			FOREIGN KEY (community_id) REFERENCES community (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 3 : city
		`CREATE TABLE temp_city (
	    insee_code int NOT NULL UNIQUE,
	    name varchar(50) NOT NULL,
	    community_code varchar(15),
			qpv boolean NOT NULL
		);`, // 4 : temp_city
		`CREATE TABLE copro (
			id SERIAL PRIMARY KEY,
			reference varchar(15) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint,
			FOREIGN KEY (zip_code) REFERENCES city (insee_code) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 5 : copro
		`CREATE TABLE temp_copro (
			reference varchar(150) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint
		);`, // 6 : temp_copro
		`CREATE TABLE budget_sector (
			id SERIAL PRIMARY KEY,
			name varchar(20) NOT NULL,
	    full_name varchar(150)
		);`, // 7 : budget_sector
		`CREATE TABLE budget_action (
			id SERIAL PRIMARY KEY,
			code bigint NOT NULL,
			name varchar(250) NOT NULL,
			sector_id int
		);`, // 8 : budget_action
		`CREATE TABLE renew_project (
			id SERIAL PRIMARY KEY,
			reference varchar(15) NOT NULL UNIQUE,
			name varchar(150) NOT NULL,
			budget bigint NOT NULL,
			prin bool NOT NULL,
			city_code1 int NOT NULL,
			city_code2 int,
			city_code3 int,			
			population int,
			composite_index int,
			FOREIGN KEY (city_code1) REFERENCES city(insee_code) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
			FOREIGN KEY (city_code2) REFERENCES city(insee_code) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
			FOREIGN KEY (city_code3) REFERENCES city(insee_code) 
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 9 : renew_project
		`CREATE TABLE temp_renew_project (
			reference varchar(15) NOT NULL UNIQUE,
			name varchar(150) NOT NULL,
			budget bigint NOT NULL,	
			prin bool NOT NULL,
			city_code1 int NOT NULL,
			city_code2 int,
			city_code3 int,			
			population int,
			composite_index int
		);`, // 10 : temp_renew_project
		`CREATE TABLE housing (
	    id SERIAL PRIMARY KEY,
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
			anru boolean NOT NULL
		);`, // 11 : housing
		`CREATE TABLE temp_housing (
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
	    anru boolean NOT NULL
		);`, // 12 : temp_housing
		`CREATE TABLE beneficiary (
	    id SERIAL PRIMARY KEY,
	    code int NOT NULL,
	    name varchar(120) NOT NULL
		);`, // 13 : beneficiary
		`CREATE TABLE commitment (
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
			sold_out boolean NOT NULL,
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
			MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 14 : commitment
		`CREATE TABLE temp_commitment (
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
			sold_out boolean NOT NULL,
			sector varchar(5) NOT NULL,
			action_code bigint,
			action_name varchar(150)
		);`, // 15 : temp_commitment
		`CREATE TABLE payment (
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
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 16 : payment
		`CREATE TABLE temp_payment (
	    commitment_year int NOT NULL,
	    commitment_code varchar(5) NOT NULL,
	    commitment_number int NOT NULL,
	    commitment_line int NOT NULL,
	    year int NOT NULL,
	    creation_date date NOT NULL,
			modification_date date NOT NULL,
			number int NOT NULL,
	    value bigint NOT NULL
		);`, // 17 : temp_payment
		`CREATE TABLE commission (
	    id SERIAL PRIMARY KEY,
	    name varchar(140) NOT NULL,
	    date date
		);`, // 18 : commission
		`CREATE TABLE renew_project_forecast (
	    id SERIAL PRIMARY KEY,
	    commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
	    renew_project_id int NOT NULL,
			FOREIGN KEY (renew_project_id) REFERENCES renew_project (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 19 : renew_project_forecast
		`CREATE TABLE temp_renew_project_forecast (
			id int NOT NULL,
			commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
	    renew_project_id int NOT NULL
		);`, // 20 : temp_renew_project_forecast
		`CREATE TABLE copro_forecast (
	    id SERIAL PRIMARY KEY,
	    commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
	    copro_id int NOT NULL,
			FOREIGN KEY (copro_id) REFERENCES copro (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);`, // 21 : copro_forecast
		`CREATE TABLE temp_copro_forecast (
			id int NOT NULL,
			commission_id int NOT NULL,
	    value bigint NOT NULL,
	    comment text,
	    copro_id int NOT NULL
		);`, // 22 : temp_copro_forecast
		`CREATE VIEW cumulated_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,q.value,
			c.beneficiary_id, c.iris_code,c.action_id,c.housing_id, c.copro_id,
			c.renew_project_id
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`, // 23 : cumulated_commitment view
		`CREATE VIEW cumulated_sold_commitment AS
		SELECT c.id,c.year,c.code,c.number,c.creation_date,c.name,q.value, c.sold_out,
			c.beneficiary_id, c.iris_code,c.action_id,c.housing_id, c.copro_id,
			c.renew_project_id
		FROM commitment c
		JOIN (SELECT year,code,number,sum(value) as value,min(creation_date),
			min(id) as id FROM commitment GROUP BY 1,2,3 ORDER BY 1,2,3) q
		ON c.id = q.id;`, // 24 : cumulated_sold_commitment view
		`CREATE TABLE ratio (
			id SERIAL PRIMARY KEY,
			year int NOT NULL,
			sector_id int NOT NULL,
			index int NOT NULL,
			ratio double precision NOT NULL,
			FOREIGN KEY (sector_id) REFERENCES budget_sector (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION
			);`, // 25 : ratio
	}
	for i, q := range queries {
		if _, err := db.Exec(q); err != nil {
			t.Errorf("Création de table [%d] : "+err.Error(), i)
			t.FailNow()
			return
		}
	}
	users := []models.User{
		{Name: "Christophe Saintillan",
			Email:    cfg.Users.Admin.Email,
			Password: cfg.Users.Admin.Password,
			Rights:   models.AdminBit | models.ActiveBit},
		{Name: "Utilisateur",
			Email:    cfg.Users.User.Email,
			Password: cfg.Users.User.Password,
			Rights:   models.ActiveBit},
		{Name: "Utilisateur copro",
			Email:    cfg.Users.CoproUser.Email,
			Password: cfg.Users.CoproUser.Password,
			Rights:   models.ActiveCoproMask},
		{Name: "Utilisateur RU",
			Email:    cfg.Users.RenewProjectUser.Email,
			Password: cfg.Users.RenewProjectUser.Password,
			Rights:   models.ActiveRenewProjectMask},
	}
	for _, u := range users {
		if err := createUser(&u, db); err != nil {
			t.Error(err.Error())
			t.FailNow()
			return
		}
	}
}

// createUser creates an new user in the test database
func createUser(u *models.User, db *sql.DB) error {
	if err := u.CryptPwd(); err != nil {
		return fmt.Errorf("Cryptage du mot de passe de %s : %v", u.Name, err)
	}
	if err := u.Create(db); err != nil {
		return fmt.Errorf("Création en base de données de %s : %v", u.Name, err)
	}
	return nil
}

// fetchTokens uses the login request to store an admin and an user token
func fetchTokens(t *testing.T, ctx *TestContext) {
	for _, u := range []*config.Credentials{
		&ctx.Config.Users.Admin,
		&ctx.Config.Users.User,
		&ctx.Config.Users.CoproUser,
		&ctx.Config.Users.RenewProjectUser} {
		response := ctx.E.POST("/api/user/login").
			WithBytes([]byte(`{"Email":"` + u.Email + `","Password":"` + u.Password + `"}`)).
			Expect()
		lr := struct{ Token string }{}
		if err := json.Unmarshal(response.Content, &lr); err != nil {
			t.Errorf(err.Error())
			t.FailNow()
			return
		}
		u.Token = lr.Token
	}
}

type tcRespFunc func(TestCase) *httpexpect.Response

// chkFactory launch the test cases against the callback function and check the status
//  and the content of a response according. If test field CountItemName is filled,
// it also checks that the count of such elements is the one give in the Count field
func chkFactory(t *testing.T, tcc []TestCase, f tcRespFunc, name string, b ...*int) bool {
	ok := true
	for i, tc := range tcc {
		response := f(tc)
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				ok = false
				t.Errorf("%s[%d]\n  ->attendu %s\n  ->reçu: %s", name, i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			ok = false
			t.Errorf("%s[%d]  ->status attendu %d  ->reçu: %d", name, i, tc.StatusCode, status)
		}
		if status == http.StatusOK && tc.CountItemName != "" {
			count := strings.Count(body, tc.CountItemName)
			if count != tc.Count {
				ok = false
				t.Errorf("%s[%d]  ->nombre attendu %d  ->reçu: %d", name, i, tc.Count, count)
			}
		}
		if status == http.StatusCreated && tc.StatusCode == http.StatusCreated && len(b) > 0 {
			index := strings.Index(body, tc.IDName)
			if index > 0 {
				fmt.Sscanf(body[index:], tc.IDName+":%d", b[0])
			}
		}
	}
	return ok
}
