package actions

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testPreProg is the entry point for testing all renew projet requests
func testPreProg(t *testing.T, c *TestContext) {
	t.Run("PreProg", func(t *testing.T) {
		testBatchCoproPreProgs(t, c)
		testGetCoproPreProgs(t, c)
		testBatchHousingPreProgs(t, c)
		testGetHousingPreProgs(t, c)
		testBatchRPPreProgs(t, c)
		testGetRPPreProgs(t, c)
		testGetPreProgs(t, c)
	})
}

// testBatchCoproPreProgs check route is copro user protected and batch import
// returns successfully
func testBatchCoproPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Params:       "Year=2019",
			Sent:         []byte(``),
			RespContains: []string{"Droits sur les copropriétés requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.CoproUser.Token,
			Params: "Year=a",
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, décodage année : "},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year
		{Token: c.Config.Users.CoproUser.Token,
			Params:       "Year=2019",
			Sent:         []byte(`{`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, décodage batch : "},
			StatusCode:   http.StatusBadRequest}, // 2 : bad JSON
		{Token: c.Config.Users.CoproUser.Token,
			Params: "Year=2019",
			Sent: []byte(`{"PreProg":[{"CommissionID":0,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 3 : commision ID nul
		{Token: c.Config.Users.CoproUser.Token,
			Params: "Year=2019",
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":0,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 4 : value nul
		{Token: c.Config.Users.CoproUser.Token,
			Params: "Year=2019",
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null}]}`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 5 : action ID nul
		{Token: c.Config.Users.CoproUser.Token,
			Params: "Year=2019",
			Sent: []byte(`{"PreProg":[{"CommissionID":3,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 6 : bad commission ID
		{Token: c.Config.Users.CoproUser.Token,
			Params: "Year=2019",
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":5}]}`),
			RespContains: []string{"Fixation de la préprogrammation copro d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 7 : bad action ID
		{Token: c.Config.Users.CoproUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Batch importé"},
			StatusCode:   http.StatusOK}, // 8 : OK
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/pre_prog/copro").WithQueryString(tc.Params).
			WithBytes(tc.Sent).WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchCoproPreProgs")
	// Content will be checked by get test
}

// testGetCoproPreProgs checks if route is copro user protected and pre prog
//  correctly sent back
func testGetCoproPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les copropriétés requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : token empty
		{Token: c.Config.Users.CoproUser.Token,
			Params:       `Year=a`,
			RespContains: []string{`Préprogrammation copro d'une année, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year param
		{Token: c.Config.Users.CoproUser.Token,
			Params: `Year=2019`,
			RespContains: []string{`"FcPreProg":[`, `"KindName":"copro4"`,
				`"ForecastValue"`, `"PreProgValue"`},
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/pre_prog/copro").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoproPreProgs")
}

// testBatchHousingPreProgs check route is housing user protected and batch
// import returns successfully
func testBatchHousingPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			Params:       "Year=2019",
			RespContains: []string{"Droits sur les projets logement requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.HousingUser.Token,
			Sent:         []byte(`{`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, décodage batch : "},
			StatusCode:   http.StatusBadRequest}, // 1 : bad JSON
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":null,"Comment":null,"ActionID":3}]}`),
			Params:       "Year=a",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, décodage année : "},
			StatusCode:   http.StatusBadRequest}, // 2 : year nul
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":0,` +
				`"Value":1000000,"KindID":null,"Comment":null,"ActionID":3}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 3 : commision ID nul
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":0,"KindID":null,"Comment":null,"ActionID":3}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 4 : value nul
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":null,"Comment":null}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 5 : action ID nul
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":3,` +
				`"Value":1000000,"KindID":null,"Comment":null,"ActionID":3}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 6 : bad commission ID
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":null,"Comment":null,"ActionID":5}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation logement d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 7 : bad action ID
		{Token: c.Config.Users.HousingUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":1000000,"KindID":null,"Comment":null,"ActionID":3}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Batch importé"},
			StatusCode:   http.StatusOK}, // 8 : OK
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/pre_prog/housing").WithQueryString(tc.Params).WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchHousingPreProgs")
	// Content will be checked by get test
}

// testGetHousingPreProgs checks if route is housing user protected and pre prog
//  correctly sent back
func testGetHousingPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets logement requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : token empty
		{Token: c.Config.Users.HousingUser.Token,
			Params:       `Year=a`,
			RespContains: []string{`Préprogrammation logement d'une année, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year param
		{Token: c.Config.Users.HousingUser.Token,
			Params: `Year=2019`,
			RespContains: []string{`"FcPreProg":[`, `"ActionCode":15400202`,
				`"KindName":null`, `"ForecastValue":`, `"PreProgValue":`},
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/pre_prog/housing").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetHousingPreProgs")
}

// testBatchRPPreProgs check route admin protected and batch import returns
// successfully
func testBatchRPPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			Params:       "Year=2019",
			RespContains: []string{"Droits sur les projets RU requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent:         []byte(`{`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, décodage batch : "},
			StatusCode:   http.StatusBadRequest}, // 1 : bad JSON
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":2000000,"KindID":2,"Comment":null,"ActionID":4}]}`),
			Params:       "Year=a",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, décodage année : "},
			StatusCode:   http.StatusBadRequest}, // 2 : year nul
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":0,` +
				`"Value":2000000,"KindID":2,"Comment":null,"ActionID":4}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 3 : commision ID nul
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":0,"KindID":2,"Comment":null,"ActionID":4}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 4 : value nul
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":2000000,"KindID":2,"Comment":null}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 5 : action ID nul
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":3,` +
				`"Value":2000000,"KindID":2,"Comment":null,"ActionID":4}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 6 : bad commission ID
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":2000000,"KindID":2,"Comment":null,"ActionID":5}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Fixation de la préprogrammation RU d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 7 : bad action ID
		{Token: c.Config.Users.RenewProjectUser.Token,
			Sent: []byte(`{"PreProg":[{"CommissionID":2,` +
				`"Value":2000000,"KindID":2,"Comment":null,"ActionID":4}]}`),
			Params:       "Year=2019",
			RespContains: []string{"Batch importé"},
			StatusCode:   http.StatusOK}, // 8 : OK
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/pre_prog/renew_project").WithQueryString(tc.Params).
			WithBytes(tc.Sent).WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchRPPreProgs")
	// Content will be checked by get test
}

// testGetRPPreProgs checks if route is renew projet user protected and pre prog
//  correctly sent back
func testGetRPPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : token empty
		{Token: c.Config.Users.RenewProjectUser.Token,
			Params:       `Year=a`,
			RespContains: []string{`Préprogrammation RU d'une année, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year param
		{Token: c.Config.Users.RenewProjectUser.Token,
			Params: `Year=2019`,
			RespContains: []string{`"FcPreProg":[`, `"KindName":"PARIS 1 - Site RU 1"`,
				`"ForecastValue":`, `"PreProgValue":2000000`},
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/pre_prog/renew_project").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPPreProgs")
}

// testGetPreProgs checks if route is admin user protected and pre prog
//  correctly sent back
func testGetPreProgs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : bad token
		{Token: c.Config.Users.Admin.Token,
			Params:       `Year=a`,
			RespContains: []string{`Préprogrammation d'une année, décodage : `},
			StatusCode:   http.StatusBadRequest}, // 1 : bad year param
		{Token: c.Config.Users.Admin.Token,
			Params:        `Year=2019`,
			RespContains:  []string{`"PreProg":[`, `"Kind":1`, `"Kind":2`, `"Kind":3`},
			Count:         3,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/pre_prog").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetPreProgs")
}
