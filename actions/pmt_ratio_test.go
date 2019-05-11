package actions

import (
	"net/http"
	"strings"
	"testing"
)

// testPmtRatio is the entry point for testing all renew projet requests
func testPmtRatio(t *testing.T, c *TestContext) {
	t.Run("PmtRatio", func(t *testing.T) {
		testGetPmtRatios(t, c)
		testBatchPmtRatios(t, c)
		testGetPmtRatiosYears(t, c)
	})
}

// testGetPmtRatios checks if route is user protected and Ratios correctly sent back
func testGetPmtRatios(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			Sent:         []byte(`Year=2017`),
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Ratios de paiements, décodage : `},
			Count:        1,
			Sent:         []byte(`Year=a`),
			StatusCode:   http.StatusInternalServerError}, // 1 : bad year parameter format
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"PmtRatio":[]`},
			Count:        0,
			Sent:         []byte(`Year=2017`),
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/ratios").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetPmtRatios[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetPmtRatios[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"Index"`)
			if count != tc.Count {
				t.Errorf("GetPmtRatios[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testBatchPmtRatios check route is limited to admin and batch import succeeds
func testBatchPmtRatios(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`"Year":2009,"Ratios":[{"Index":0,"Ratio":0.1},{"Index":1,"Ratio":0.2},{"Index":2,"Ratio":0.3}]}`),
			RespContains: []string{"Batch de ratios de paiement, décodage :"},
			StatusCode:   http.StatusBadRequest}, // 1 : bad payload
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{"Year":2009,"Ratios":[{"Index":0,"Ratio":0.1},{"Index":1,"Ratio":0.2},{"Index":2,"Ratio":0.3}]}`),
			RespContains: []string{"Batch de ratios de paiement traité"},
			Count:        3,
			StatusCode:   http.StatusOK}, // 2 : OK
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/ratios").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchRatios[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchRatios[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK && tc.StatusCode == http.StatusOK {
			var count int
			if err := c.DB.QueryRow(`SELECT count(1) FROM ratio WHERE year=2009`).
				Scan(&count); err != nil {
				t.Errorf("BatchRatios[%d]\n  ->impossible de vérifier %v", i, err)
				return
			}
			if count != tc.Count {
				t.Errorf("BatchRatios[%d]  ->nombre attendu %d  ->trouvé: %d", i, tc.Count, count)
			}
		}
	}
}

// testGetPmtRatiosYears checks if route is user protected and Ratios Years are
// correctly sent back
func testGetPmtRatiosYears(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 1 : bad year parameter format
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"PmtRatiosYear":[2009]`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/ratios/years").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetPmtRatiosYears[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetPmtRatiosYears[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}
