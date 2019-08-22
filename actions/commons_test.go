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
	HousingID      int64
	RPEventTypeID  int64
}

// TestCase is used as common structure for all request tests
type TestCase struct {
	Sent          []byte
	Token         string
	Params        string
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
	testDepartment(t, cfg)
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
	testHousingForecast(t, cfg)
	testCoproForecast(t, cfg)
	testSettings(t, cfg)
	testHome(t, cfg)
	testBeneficiaryDatas(t, cfg)
	testBeneficiaryPayments(t, cfg)
	testPmtRatio(t, cfg)
	testPmtForecasts(t, cfg)
	testCmtForecasts(t, cfg)
	testLinkCommitmentsHousings(t, cfg)
	testCoproCommitmentLink(t, cfg)
	testRPEventType(t, cfg)
	testRPEvent(t, cfg)
	testRenewProjectReport(t, cfg)
	testRPPerCommunityReport(t, cfg)
	testRPCmtCityJoin(t, cfg)
	testDepartmentReport(t, cfg)
	testCityReport(t, cfg)
	testPreProg(t, cfg)
	testProg(t, cfg)
	testRPLS(t,cfg)
}

func initializeTests(t *testing.T) *TestContext {
	testCtx := &TestContext{}
	cfg := &config.PreLoRuGoConf{}
	var err error
	testCtx.App = iris.New().Configure(iris.WithConfiguration(
		iris.Configuration{DisablePathCorrection: true}))
	logFile, err := cfg.Get(testCtx.App)
	if err != nil {
		t.Errorf("Configuration : %v", err)
		t.FailNow()
	}
	cfg.App.Stage = config.TestStage
	if logFile != nil {
		defer logFile.Close()
	}
	testCtx.App.Logger().Infof("Lancement des tests\n")
	testCtx.Config = cfg
	testCtx.DB, err = config.InitDatabase(cfg, testCtx.App, true, false)
	if err != nil {
		t.Error("Erreur de connexion à postgres : " + err.Error())
		t.FailNow()
		return nil
	}
	createUsers(t, testCtx.DB, testCtx.Config)
	SetRoutes(testCtx.App, testCtx.Config.Users.SuperAdmin.Email, testCtx.DB)
	testCtx.E = httptest.New(t, testCtx.App)
	// Fetch admin and user tokens
	fetchTokens(t, testCtx)
	return testCtx
}

func createUsers(t *testing.T, db *sql.DB, cfg *config.PreLoRuGoConf) {
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
		{Name: "Utilisateur logement",
			Email:    cfg.Users.HousingUser.Email,
			Password: cfg.Users.HousingUser.Password,
			Rights:   models.ActiveHousingMask},
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
		&ctx.Config.Users.RenewProjectUser,
		&ctx.Config.Users.HousingUser} {
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
