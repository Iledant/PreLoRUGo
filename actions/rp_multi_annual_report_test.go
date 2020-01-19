package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPMultiAnnualReport is the entry point for testing all renew projet requests
func testRPMultiAnnualReport(t *testing.T, c *TestContext) {
	t.Run("RPMultiAnnualReport", func(t *testing.T) {
		testGetRPMultiAnnualReport(t, c)
	})
}

// testGetRPMultiAnnualReports checks if route is user protected and RPMultiAnnualReports correctly sent back
func testGetRPMultiAnnualReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`"RPMultiAnnualReport":[`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_project/multi_annual_report").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetRPMultiAnnualReport") {
		t.Error(r)
	}
}
