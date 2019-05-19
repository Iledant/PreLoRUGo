package actions

import (
	"net/http"
	"strings"
	"testing"
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
			RespContains: []string{`"PmtForecast":[]`},
			Count:        0,
			Sent:         []byte(`Year=2009`),
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/payments/forecasts").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetPmtForecasts[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetPmtForecasts[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"Index"`)
			if count != tc.Count {
				t.Errorf("GetPmtForecasts[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
