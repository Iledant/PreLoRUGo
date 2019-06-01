package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testPmtForecasts is the entry point for testing all renew projet requests
func testPmtForecasts(t *testing.T, c *TestContext) {
	t.Run("PmtForecasts", func(t *testing.T) {
		testGetPmtForecasts(t, c)
	})
}

// testGetPmtForecasts checks if route is admin protected and forecasts
// correctly sent back
func testGetPmtForecasts(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			Sent:         []byte(`Year=2017`),
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			Count:        1,
			Sent:         []byte(`Year=a`),
			StatusCode:   http.StatusUnauthorized}, // 1 : bad year parameter format
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Prévisions de paiements, décodage : `},
			Count:        1,
			Sent:         []byte(`Year=a`),
			StatusCode:   http.StatusInternalServerError}, // 2 : bad year parameter format
		{Token: c.Config.Users.Admin.Token,
			RespContains:  []string{`"PmtForecast":[]`},
			Count:         0,
			CountItemName: `"Index"`,
			Sent:          []byte(`Year=2009`),
			StatusCode:    http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/payments/forecasts").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetPmtForecasts")
}
