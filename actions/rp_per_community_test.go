package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPPerCommunityReport is the entry point for testing all renew projet requests
func testRPPerCommunityReport(t *testing.T, c *TestContext) {
	t.Run("RPPerCommunityReport", func(t *testing.T) {
		testGetRPPerCommunityReport(t, c)
	})
}

// testGetRPPerCommunityReport checks if route is user protected and
// RPPerCommunityReport correctly sent back
func testGetRPPerCommunityReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"RPPerCommunityReport":[{"CommunityID":2,` +
				`"CommunityName":"CA SAINT GERMAIN BOUCLES DE SEINE (78-YVELINES)",` +
				`"CommunityBudget":0,"Commitment":0,"Payment":0},{"CommunityID":4,` +
				`"CommunityName":"VILLE DE PARIS (EPT1)","CommunityBudget":0,` +
				`"Commitment":0,"Payment":0}]`},
			StatusCode: http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_project/report_per_community").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPPerCommunityReport")
}
