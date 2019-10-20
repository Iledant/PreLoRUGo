package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCoproReport is the entry point for testing all renew projet requests
func testCoproReport(t *testing.T, c *TestContext) {
	t.Run("CoproReport", func(t *testing.T) {
		testGetCoproReport(t, c)
	})
}

// testGetCoproReports checks if route is user protected and CoproReports correctly sent back
func testGetCoproReport(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"CoproReport":[`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro/report").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoproReport")
}
