package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
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
			RespContains:  []string{`"PmtRatio":[{"Index":0,"SectorID":1,"SectorName":"LO","Ratio":0.8},{"Index":1,"SectorID":1,"SectorName":"LO","Ratio":0.2}]`},
			Count:         2,
			CountItemName: `"Index"`,
			Sent:          []byte(`Year=2010`),
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/ratios").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetPmtRatios")
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
			Sent:         []byte(`{"Year":2009,"Ratios":[{"Index":0,"SectorID":1,"Ratio":0.1},{"Index":1,"SectorID":1,"Ratio":0.2},{"Index":2,"SectorID":1,"Ratio":0.3}]}`),
			RespContains: []string{"Batch de ratios de paiement traité"},
			StatusCode:   http.StatusOK}, // 2 : OK
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/ratios").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	if chkFactory(t, tcc, f, "BatchRatios") {
		var count int
		if err := c.DB.QueryRow(`SELECT count(1) FROM ratio WHERE year=2009`).
			Scan(&count); err != nil {
			t.Errorf("BatchRatios[final]\n  ->impossible de vérifier %v", err)
			return
		}
		if count != 3 {
			t.Errorf("BatchRatios[final]  ->nombre attendu %d  ->trouvé: %d", 3, count)
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
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/ratios/years").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetPmtRatiosYears")
}
