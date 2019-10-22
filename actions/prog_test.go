package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testProg is the entry point for testing all programming tests
func testProg(t *testing.T, c *TestContext) {
	t.Run("Prog", func(t *testing.T) {
		testBatchProg(t, c)
		testGetProg(t, c)
		testGetProgDatas(t, c)
		testGetProgYears(t, c)
	})
}

// testBatchProg check route is admin user protected and batch import
// returns successfully
func testBatchProg(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent:         []byte(`{`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, décodage batch : "},
			StatusCode:   http.StatusBadRequest}, // 1 : bad JSON
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			Params:       "Year=a",
			RespContains: []string{"Fixation de la programmation d'une année, décodage année : "},
			StatusCode:   http.StatusBadRequest}, // 2 : bad year
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":0,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 3 : commision ID nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,` +
				`"Value":0,"KindID":5,"Comment":null,"ActionID":2}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 4 : value nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 5 : action ID nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 6 : kind nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":3,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2,"Kind":2}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 7 : bad commission ID
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":5,"Kind":2}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 8 : bad action ID
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Value":1000000,"KindID":5,"Comment":null,"ActionID":2,"Kind":2},
			{"CommissionID":2,"Value":2000000,"KindID":null,"Comment":null,"ActionID":3,"Kind":1},
			{"CommissionID":2,"Value":3000000,"KindID":3,"Comment":"commentaire RU","ActionID":4,"Kind":3}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Batch importé"},
			StatusCode:   http.StatusOK}, // 9 : OK
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/prog").WithQueryString(tc.Params).WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchProg")
	// Content will be checked by get test
}

// testGetProg checks if route is user protected and prog correctly sent back
func testGetProg(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: `fake`,
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Params:       `?Year=a`,
			RespContains: []string{`Programmation d'une année, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year param
		{Token: c.Config.Users.User.Token,
			Params: `Year=2019`,
			RespContains: []string{`"Prog":[`, `"Kind":1`, `"Kind":2`, `"Kind":3`,
				`"commentaire RU"`, `"Value":1000000`, `"ForecastValue"`, `"PreProgValue"`},
			Count:         5,
			CountItemName: `"Kind"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/prog").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetProg")
}

// testGetProgDatas checks if route is user protected and prog and others datas
// correctly sent back
func testGetProgDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: `fake`,
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Params:       `?Year=a`,
			RespContains: []string{`Données de programmation d'une année, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year param
		{Token: c.Config.Users.User.Token,
			Params: `Year=2019`,
			RespContains: []string{`"Prog":[`, `"Kind":1`, `"Kind":2`, `"Kind":3`,
				`"RenewProject":[`, `"Copro":[`, `"commentaire RU"`, `"Value":1000000`,
				`"BudgetAction":[`, `"Commission":[`},
			Count:         5,
			CountItemName: `"Kind"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/prog/datas").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetProgDatas")
}

// testGetProgYears checks if route is user protected and programmation years
// correctly sent back
func testGetProgYears(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: `fake`,
			RespContains: []string{`Token invalide`},
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"ProgYear":[2019]`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/prog/years").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetProgYears")
}
