package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingSummary is the entry point for testing housing summary requests
func testHousingSummary(t *testing.T, c *TestContext) {
	t.Run("HousingSummary", func(t *testing.T) {
		testBatchHousingSummary(t, c)
		testBatchIRISHousingType(t, c)
	})
}

var (
	tcSentBH1 = []byte(`{`)
	tcSentBH2 = []byte(`{"HousingSummary":[{"InseeCode":null,"Address":"32 rue de Rivoli","PLAI":2,"PLS":3,"PLUS":4,"IRISCode":"EX000001","ReferenceCode":"PARIS1LLS","ANRU":false}]}`)
	tcSentBH3 = []byte(`{"HousingSummary":[{"InseeCode":75101,"Address":null,"PLAI":2,"PLS":3,"PLUS":4,"IRISCode":"EX000001","ReferenceCode":"PARIS1LLS","ANRU":false}]}`)
	tcSentBH4 = []byte(`{"HousingSummary":[{"InseeCode":75101,"Address":"32 rue de Rivoli","IRISCode":"EX000001","ReferenceCode":"PARIS1LLS","ANRU":false}]}`)
	tcSentBH5 = []byte(`{"HousingSummary":[{"InseeCode":75101,"Address":"32 rue de Rivoli","PLAI":2,"PLS":3,"PLUS":4,"ReferenceCode":"PARIS1LLS","ANRU":false}]}`)
	tcSentBH6 = []byte(`{"HousingSummary":[{"InseeCode":75101,"Address":"32 rue de Rivoli","PLAI":2,"PLS":3,"PLUS":4,"IRISCode":"EX000001","ANRU":false}]}`)
	tcSentBH7 = []byte(`{"HousingSummary":[{"InseeCode":75000,"Address":"32 rue de Rivoli","PLAI":2,"PLS":3,"PLUS":4,"IRISCode":"EX000001","ReferenceCode":"PARIS1LLS","ANRU":false}]}`)
	tcSentBH8 = []byte(`{"HousingSummary":[{"InseeCode":75101,"Address":"32 rue de Rivoli","PLAI":2,"PLS":3,"PLUS":4,"IRISCode":"EX000001","ReferenceCode":"PARIS1LLS","ANRU":false}]}`)
)

// testBatchHousingSummary check route is limited to housing user and batch
// import succeeds
func testBatchHousingSummary(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.HousingCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH1,
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Batch de bilan logements, décodage :`},
		}, // 1 : token empty
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH2,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de bilan logements, requête : line 1 InseeCode nul`},
		}, // 2 : inseecode nul
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH3,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de bilan logements, requête : line 1 Address nul`},
		}, // 3 : address nul
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH4,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de bilan logements, requête : line 1 PLS PLAI and PLUS nul`},
		}, // 4 : address nul
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH5,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de bilan logements, requête : line 1 IrisCode nul`},
		}, // 5 : IrisCode nul
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH6,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de bilan logements, requête : line 1 ReferenceCode nul`},
		}, // 6 : ReferenceCode nul
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH7,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de bilan logements, requête :`},
		}, // 7 : insee_code not in database
		{
			Token:        c.Config.Users.HousingUser.Token,
			Sent:         tcSentBH8,
			StatusCode:   http.StatusOK,
			RespContains: []string{`Batch de bilan logements importé`},
		}, // 8 : insee_code not in database
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_summary").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "BatchHousingSummary") {
		t.Error(r)
	}
}

// testBatchIRISHousingType check route is admin protected and batch
// import succeeds
func testBatchIRISHousingType(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IRISHousingType":[{"IRISCode":"","HousingTypeShortName":"LF"}`),
			StatusCode:   http.StatusBadRequest,
			RespContains: []string{`Batch de lien IRIS / type de logement, requête :`},
		}, // 1 : bad payload
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IRISHousingType":[{"IRISCode":"","HousingTypeShortName":"LF"}]}`),
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de lien IRIS / type de logement, requête :`},
		}, // 2 : IRISCode empty
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IRISHousingType":[{"IRISCode":"EX000001","HousingTypeShortName":""}]}`),
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Batch de lien IRIS / type de logement, requête :`},
		}, // 3 : HousingTypeShortName empty
		{
			Token:        c.Config.Users.Admin.Token,
			Sent:         []byte(`{"IRISHousingType":[{"IRISCode":"EX000001","HousingTypeShortName":"LF"}]}`),
			StatusCode:   http.StatusOK,
			RespContains: []string{`Batch de lien IRIS / type de logement importé`},
		}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/iris_housing_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "BatchIRISHousingType") {
		t.Error(r)
	}
	tcc = []TestCase{
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusOK,
			RespContains: []string{`"HousingTypeID":` + strconv.Itoa(c.HousingTypeID)},
		}, // 1 : bad payload
	}
	f = func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housings").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "BatchIRISHousingTypeCheck") {
		t.Error(r)
	}
}
