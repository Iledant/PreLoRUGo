package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testPlacement is the entry point for testing all placements requests
func testPlacement(t *testing.T, c *TestContext) {
	t.Run("Placement", func(t *testing.T) {
		testBatchPlacements(t, c)
		testGetPlacements(t, c)
	})
}

// testBatchPlacements check route is limited to admin and batch import succeeds
func testBatchPlacements(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Placement":[{"IrisCode":"","Count":1,"ContractYear":null},
			{"IrisCode":"14004240","Count":0,"ContractYear":2019}]}`),
			RespContains: []string{"Batch de stages, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : IrisCode empty
		{
			Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Placement":[{"IrisCode":"13021233","Count":1,"ContractYear":null},
			{"IrisCode":"14004240","Count":0,"ContractYear":2019}]}`),
			RespContains: []string{"Batch de stages importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/placements").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "BatchPlacements") {
		t.Error(r)
	}
	// GetPlacements checks if data are correctly analyzed
}

// testGetPlacements checks if route is user protected and Placements correctly sent back
func testGetPlacements(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"IrisCode":"13021233","Count":1,"ContractYear":null`,
				`"IrisCode":"14004240","Count":0,"ContractYear":2019`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/placements").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetPlacements") {
		t.Error(r)
	}
}

// testGetBeneficiaryPlacements checks if route is user protected and Placements correctly sent back
func testGetBeneficiaryPlacements(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:         c.Config.Users.User.Token,
			Params:        "2",
			RespContains:  []string{`"IrisCode":"14004240","Count":0,"ContractYear":2019`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/beneficiary/"+tc.Params+"/placements").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetBeneficiaryPlacements") {
		t.Error(r)
	}
}
