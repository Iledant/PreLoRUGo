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
			RespContains: []string{"Fixation de la programmation d'une année, décodage : "},
			StatusCode:   http.StatusBadRequest}, // 1 : bad JSON
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":0,"Year":2019,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 2 : commision ID nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":2019,` +
				`"Value":0,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 3 : value nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":0,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 4 : year nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":2019,` +
				`"Value":1000000,"KindID":5,"Comment":null}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 5 : action ID nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":2019,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 6 : kind nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":3,"Year":2019,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":2,"Kind":"Copro"}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 7 : bad commission ID
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":2019,` +
				`"Value":1000000,"KindID":5,"Comment":null,"ActionID":5,"Kind":"Copro"}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 8 : bad action ID
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":2019,"Value":1000000,
			"KindID":5,"Comment":null,"ActionID":2,"Kind":"Copro"},
			{"CommissionID":2,"Year":2018,"Value":1000000,"KindID":5,"Comment":null,
			"ActionID":2,"Kind":"Copro"}]}`),
			RespContains: []string{"Fixation de la programmation d'une année, " +
				"requête : more than one year in batch"},
			StatusCode: http.StatusInternalServerError}, // 9 : two different years in one prog
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Prog":[{"CommissionID":2,"Year":2019,"Value":1000000,"KindID":5,"Comment":null,"ActionID":2,"Kind":"Copro"},
			{"CommissionID":2,"Year":2019,"Value":2000000,"KindID":null,"Comment":null,"ActionID":3,"Kind":"Housing"},
			{"CommissionID":2,"Year":2019,"Value":3000000,"KindID":3,"Comment":"commentaire RU","ActionID":4,"Kind":"RenewProject"}]}`),
			RespContains: []string{"Batch importé"},
			StatusCode:   http.StatusOK}, // 10 : OK
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/prog").WithBytes(tc.Sent).
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
			RespContains: []string{`"Prog":[`, `"Housing"`, `"RenewProject"`, `"Copro"`,
				`"commetaire RU"`, `"Value":1000000`},
			Count:         3,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/prog").WithQueryString(tc.Params).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetProg")
}
