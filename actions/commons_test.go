package actions

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/Iledant/PreLoRUGo/config"
	"github.com/Iledant/PreLoRUGo/models"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
)

// TestContext contains all items for units tests in API.
type TestContext struct {
	DB     *sql.DB
	App    *iris.Application
	E      *httpexpect.Expect
	Config *config.PreLoRuGoConf
}

// TestCase is used as common structure for all request tests
type TestCase struct {
	Sent         []byte
	Token        string
	RespContains []string
	StatusCode   int
	ID           int
	Count        int
}

// TestAll embeddes all test functions and is the only test entry point
// It initializes a fresh new test database base and call test functions
// in the right order to avoid side effects
func TestAll(t *testing.T) {
	cfg := initializeTests(t)
	testUser(t, cfg)
	testCopro(t, cfg)
	testBudgetAction(t, cfg)
	testRenewProject(t, cfg)
	testHousing(t, cfg)
	testCommitment(t, cfg)
	testBeneficiary(t, cfg)

	testPayment(t, cfg)
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
	if _, err := db.Exec(`DROP TABLE IF EXISTS copro, users, imported_commitment, 
	commitment, imported_payment, payment, report, budget_action, beneficiary, 
	temp_copro, renew_project, temp_renew_project, housing, temp_housing, commitment , temp_commitment, beneficiary, payment , temp_payment `); err != nil {
		t.Error("Suppression des tables : " + err.Error())
		t.FailNow()
		return
	}
	queries := []string{`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name varchar(50) NOT NULL,
		email varchar(120) NOT NULL,
		password varchar(120) NOT NULL,
		role varchar(15) NOT NULL,
		active boolean NOT NULL
		);`, // 0 : users
		`CREATE TABLE copro (
			id SERIAL PRIMARY KEY,
			reference varchar(15) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint
		);`, // 1 : copro
		`CREATE TABLE temp_copro (
			reference varchar(150) NOT NULL,
			name varchar(150) NOT NULL,
			address varchar(200) NOT NULL,
			zip_code int NOT NULL,
			label_date date,
			budget bigint
		);`, // 2 : temp_copro
		`CREATE TABLE budget_action (
			id SERIAL PRIMARY KEY,
			code varchar(12) NOT NULL,
			name varchar(250) NOT NULL,
			sector_id int
		);`, // 3 : budget_action
		`CREATE TABLE renew_project (
			id SERIAL PRIMARY KEY,
			reference varchar(15) NOT NULL UNIQUE,
			name varchar(150) NOT NULL,
			budget bigint NOT NULL,
			population int,
			composite_index int
		);`, // 4 : renew_project
		`CREATE TABLE temp_renew_project (
			reference varchar(15) NOT NULL UNIQUE,
			name varchar(150) NOT NULL,
			budget bigint NOT NULL,	
			population int,
			composite_index int
		);`, // 5 : temp_renew_project
		`CREATE TABLE housing (
	    id SERIAL PRIMARY KEY,
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
	    anru boolean NOT NULL
		);`, // 7 : housing
		`CREATE TABLE temp_housing (
	    reference varchar(100) NOT NULL,
	    address varchar(150),
	    zip_code int,
	    plai int NOT NULL,
	    plus int NOT NULL,
	    pls int NOT NULL,
	    anru boolean NOT NULL
		);`, // 8 : temp_housing
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
	    iris_code varchar(20)
		);`, // 8 : commitment
		`CREATE TABLE temp_commitment (
	    year int NOT NULL,
	    code varchar(5) NOT NULL,
	    number int NOT NULL,
	    line int NOT NULL,
	    creation_date int NOT NULL,
	    modification_date int NOT NULL,
	    name varchar(150) NOT NULL,
	    value bigint NOT NULL,
	    beneficiary_code int NOT NULL,
	    beneficiary_name varchar(150) NOT NULL,
	    iris_code varchar(20)
		);`, // 9 : temp_commitment
		`CREATE TABLE beneficiary (
	    id SERIAL PRIMARY KEY,
	    code int NOT NULL,
	    name varchar(120) NOT NULL
		);`, // 10 : beneficiary
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
			value bigint NOT NULL,
			CONSTRAINT payment_commitment_id_fkey FOREIGN KEY (commitment_id)
			REFERENCES commitment (id) MATCH SIMPLE
			ON UPDATE NO ACTION ON DELETE NO ACTION
				);`, // 11 : payment
		`CREATE TABLE temp_payment (
	    commitment_year int NOT NULL,
	    commitment_code varchar(5) NOT NULL,
	    commitment_number int NOT NULL,
	    commitment_line int NOT NULL,
	    year int NOT NULL,
	    creation_date date NOT NULL,
	    modification_date date NOT NULL,
	    value bigint NOT NULL
		);`, // 12 : temp_payment
	}
	for i, q := range queries {
		if _, err := db.Exec(q); err != nil {
			t.Errorf("Création de table [%d] : "+err.Error(), i)
			t.FailNow()
			return
		}
	}
	admin := models.User{Name: "Christophe Saintillan",
		Email:    cfg.Users.Admin.Email,
		Password: cfg.Users.Admin.Password,
		Role:     models.AdminRole,
		Active:   true}
	if err := admin.CryptPwd(); err != nil {
		t.Error("Cryptage mot de passe admin : " + err.Error())
		t.FailNow()
		return
	}
	if err := admin.Create(db); err != nil {
		t.Error("Requête admin create : " + err.Error())
		t.FailNow()
		return
	}
	user := models.User{Name: "Utilisateur",
		Email:    cfg.Users.User.Email,
		Password: cfg.Users.User.Password,
		Role:     models.UserRole,
		Active:   true}
	if err := user.CryptPwd(); err != nil {
		t.Error("Cryptage mot de passe user : " + err.Error())
		t.FailNow()
		return
	}
	if err := user.Create(db); err != nil {
		t.Error("Requête user create : " + err.Error())
		t.FailNow()
		return
	}
}

// fetchTokens uses the login request to store an admin and an user token
func fetchTokens(t *testing.T, ctx *TestContext) {
	for _, u := range []*config.Credentials{
		&ctx.Config.Users.Admin,
		&ctx.Config.Users.User} {
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
