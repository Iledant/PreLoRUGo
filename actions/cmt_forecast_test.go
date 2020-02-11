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
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`{"CmtForecast":[]}`},
			Count:         0,
			CountItemName: `"Index"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/commitments/forecasts").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetCmtForecasts") {
		t.Error(r)
	}
}
