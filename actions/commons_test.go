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

// TestContext contains all items for units tests in API.
type TestContext struct {
	DB                          *sql.DB
	App                         *iris.Application
	E                           *httpexpect.Expect
	Config                      *config.PreLoRuGoConf
	CommissionID                int64
	RenewProjectID              int64
	CoproID                     int64
	HousingID                   int64
	RPEventTypeID               int64
	CoproEventTypeID            int64
	BeneficiaryGroupID          int
	HousingTypologyID           int
	HousingConventionID         int
	HousingCommentID            int
	HousingTransferID           int
	ConventionTypeID            int
	HousingTypeID               int
	AdminCheckTestCase          *TestCase
	UserCheckTestCase           *TestCase
	CoproCheckTestCase          *TestCase
	CoproPreProgCheckTestCase   *TestCase
	HousingCheckTestCase        *TestCase
	HousingPreProgCheckTestCase *TestCase
	RPCheckTestCase             *TestCase
	RPPreProgCheckTestCase      *TestCase
	ReservationFeeCheckTestCase *TestCase
}

// TestAll embeddes all test functions and is the only test entry point
// It initializes a fresh new test database base and call test functions
// in the right order to avoid side effects
func TestAll(t *testing.T) {
	cfg := initializeTests(t)
	testUser(t, cfg)
	testHomeMessage(t, cfg)
	testDepartment(t, cfg)
	testCommunity(t, cfg)
	testCity(t, cfg)
	testCopro(t, cfg)
	testBudgetAction(t, cfg)
	testRenewProject(t, cfg)
	testHousingType(t, cfg)
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
	testRPLS(t, cfg)
	testSummaries(t, cfg)
	testHousingSummary(t, cfg)
	testCoproEventType(t, cfg)
	testCoproEvent(t, cfg)
	testCoproDoc(t, cfg)
	testCoproReport(t, cfg)
	testRPMultiAnnualReport(t, cfg)
	testPaymentCredits(t, cfg)
	testPaymentCreditJournals(t, cfg)
	testPlacement(t, cfg)
	testBeneficiaryGroup(t, cfg)
	testBeneficiaryGroupDatas(t, cfg)
	testHousingTypology(t, cfg)
	testHousingConvention(t, cfg)
	testHousingComment(t, cfg)
	testHousingTransfer(t, cfg)
	testConventionType(t, cfg)
	testReservationFee(t, cfg)
	testGetDifActionPaymentPrevisions(t, cfg)
	testReservationReport(t, cfg)
	testSoldCommitment(t, cfg)
	testAvgPmtTime(t, cfg)
	testPaymentDemands(t, cfg)
	testPaymentDelays(t, cfg)
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
	fetchTokens(t, testCtx)
	createTestCases(testCtx)
	return testCtx
}

func createUsers(t *testing.T, db *sql.DB, cfg *config.PreLoRuGoConf) {
	users := []models.User{
		{
			Name:     "Christophe Saintillan",
			Email:    cfg.Users.Admin.Email,
			Password: cfg.Users.Admin.Password,
			Rights:   models.AdminBit | models.ActiveBit},
		{
			Name:     "Utilisateur",
			Email:    cfg.Users.User.Email,
			Password: cfg.Users.User.Password,
			Rights:   models.ActiveBit},
		{
			Name:     "Utilisateur copro",
			Email:    cfg.Users.CoproUser.Email,
			Password: cfg.Users.CoproUser.Password,
			Rights:   models.ActiveCoproMask},
		{
			Name:     "Utilisateur pre prog copro",
			Email:    cfg.Users.CoproPreProgUser.Email,
			Password: cfg.Users.CoproPreProgUser.Password,
			Rights:   models.ActiveCoproPreProgMask},
		{
			Name:     "Utilisateur RU",
			Email:    cfg.Users.RenewProjectUser.Email,
			Password: cfg.Users.RenewProjectUser.Password,
			Rights:   models.ActiveRenewProjectMask},
		{
			Name:     "Utilisateur pre prog RU",
			Email:    cfg.Users.RenewProjectPreProgUser.Email,
			Password: cfg.Users.RenewProjectPreProgUser.Password,
			Rights:   models.ActiveRenewProjectPreProgMask},
		{
			Name:     "Utilisateur logement",
			Email:    cfg.Users.HousingUser.Email,
			Password: cfg.Users.HousingUser.Password,
			Rights:   models.ActiveHousingMask},
		{
			Name:     "Utilisateur pre prog logement",
			Email:    cfg.Users.HousingPreProgUser.Email,
			Password: cfg.Users.HousingPreProgUser.Password,
			Rights:   models.ActiveHousingPreProgMask},
		{
			Name:     "Utilisateur réservation",
			Email:    cfg.Users.ReservationFeeUser.Email,
			Password: cfg.Users.ReservationFeeUser.Password,
			Rights:   models.ActiveReservationMask},
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
	var lr struct{ Token string }
	for _, u := range []*config.Credentials{
		&ctx.Config.Users.Admin,
		&ctx.Config.Users.User,
		&ctx.Config.Users.CoproUser,
		&ctx.Config.Users.CoproPreProgUser,
		&ctx.Config.Users.RenewProjectUser,
		&ctx.Config.Users.RenewProjectPreProgUser,
		&ctx.Config.Users.HousingUser,
		&ctx.Config.Users.HousingPreProgUser,
		&ctx.Config.Users.ReservationFeeUser,
	} {
		c := fmt.Sprintf(`{"Email":"%s","Password":"%s"}`, u.Email, u.Password)
		response := ctx.E.POST("/api/user/login").WithBytes([]byte(c)).Expect()
		if err := json.Unmarshal(response.Content, &lr); err != nil {
			t.Errorf(err.Error())
			t.FailNow()
			return
		}
		u.Token = lr.Token
	}
}

func createTestCases(ctx *TestContext) {
	ctx.AdminCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits administrateur requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.UserCheckTestCase = &TestCase{
		Token:        "",
		RespContains: []string{`Token absent`},
		StatusCode:   http.StatusInternalServerError,
	}
	ctx.CoproCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits sur les copropriétés requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.CoproPreProgCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits préprogrammation sur les copropriétés requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.HousingCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits sur les projets logement requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.HousingPreProgCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits préprogrammation sur les projets logement requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.RPCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits sur les projets RU requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.RPPreProgCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits préprogrammation sur les projets RU requis`},
		StatusCode:   http.StatusUnauthorized,
	}
	ctx.ReservationFeeCheckTestCase = &TestCase{
		Token:        ctx.Config.Users.User.Token,
		RespContains: []string{`Droits sur les réservations requis`},
		StatusCode:   http.StatusUnauthorized,
	}
}

type tcRespFunc func(TestCase) *httpexpect.Response

// chkFactory launch the test cases against the callback function and check the status
//  and the content of a response according. If test field CountItemName is filled,
// it also checks that the count of such elements is the one give in the Count field
func chkFactory(tcc []TestCase, f tcRespFunc, name string, b ...*int) []string {
	var resp []string
	for i, tc := range tcc {
		response := f(tc)
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				resp = append(resp,
					fmt.Sprintf("%s[%d]\n  ->attendu %s\n  ->reçu: %s", name, i, r, body))
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			resp = append(resp,
				fmt.Sprintf("%s[%d]  ->status attendu %d  ->reçu: %d", name, i, tc.StatusCode, status))
		}
		if status == http.StatusOK && tc.CountItemName != "" {
			count := strings.Count(body, tc.CountItemName)
			if count != tc.Count {
				resp = append(resp,
					fmt.Sprintf("%s[%d]  ->nombre attendu %d  ->reçu: %d", name, i, tc.Count, count))
			}
		}
		if status == http.StatusCreated && tc.StatusCode == http.StatusCreated && len(b) > 0 {
			index := strings.Index(body, tc.IDName)
			if index > 0 {
				fmt.Sscanf(body[index:], tc.IDName+":%d", b[0])
			}
		}
	}
	return resp
}
