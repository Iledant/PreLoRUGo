package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCmtForecasts is the entry point for testing all commitments forecast requests
func testCmtForecasts(t *testing.T, c *TestContext) {
	t.Run("CmtForecasts", func(t *testing.T) {
		testGetCmtForecasts(t, c)
	})
}

// testGetCmtForecasts checks if route is admin protected and forecasts
// correctly sent back
func testGetCmtForecasts(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			Count:        1,
			StatusCode:   http.StatusUnauthorized}, // 1 : bad year parameter format
		{Token: c.Config.Users.Admin.Token,
			RespContains:  []string{`{"CmtForecast":[]}`},
			Count:         0,
			CountItemName: `"Index"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/forecasts").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCmtForecasts")
}