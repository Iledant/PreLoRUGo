package actions

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRenewProjectForecast is the entry point for testing all renew projet requests
func testRenewProjectForecast(t *testing.T, c *TestContext) {
	t.Run("RenewProjectForecast", func(t *testing.T) {
		ID := testCreateRenewProjectForecast(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la prévision RU")
			t.FailNow()
			return
		}
		testUpdateRenewProjectForecast(t, c, ID)
		testGetRenewProjectForecast(t, c, ID)
		testGetRenewProjectForecasts(t, c, ID)
		testDeleteRenewProjectForecast(t, c, ID)
		testBatchRenewProjectForecasts(t, c)
	})
}

// testCreateRenewProjectForecast checks if route is admin protected and created budget action
// is properly filled
func testCreateRenewProjectForecast(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":0,"Value":1000000,"Comment":"Essai","RenewProjectID":1000000}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":0,"Value":1000000,"Comment":"Essai","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : commission ID null
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","RenewProjectID":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : renew project ID null
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : value null
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 5 : actionID null
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Project":null,"Comment":"Essai","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":2}}`),
			Token:  c.Config.Users.RenewProjectUser.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"RenewProjectForecast":{"ID":1,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z","CommissionName":"Commission test",` +
				`"Value":1000000,"Project":null,"Comment":"Essai","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":2,"ActionCode":15400403,"ActionName":"Aide aux copropriétés en difficulté"`},
			StatusCode: http.StatusCreated}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/renew_project_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateRenewProjectForecast", &ID)
	return ID
}

// testUpdateRenewProjectForecast checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateRenewProjectForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":2000000,"Value":2000000,"Comment":null,"RenewProjectID":2000000}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : value nul
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":0,"Value":0,"Comment":"Essai2","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : commission ID nul
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","RenewProjectID":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : renew project ID nul
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","RenewProjectID":` + strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 5 : action ID nul
		{Sent: []byte(`{"RenewProjectForecast":{"ID":0,"CommissionID":2000000,"Value":2000000,"Comment":null,"RenewProjectID":2000000,"ActionID":3}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{Sent: []byte(`{"RenewProjectForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Project":"projet","Comment":"Essai2","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":3}}`),
			Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`"RenewProjectForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z","CommissionName":"Commission test",` +
				`"Value":2000000,"Project":"projet","Comment":"Essai2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":3,"ActionCode":15400202,"ActionName":"Aide à la création de logements locatifs sociaux"}`},
			StatusCode: http.StatusOK}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/renew_project_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateRenewProjectForecast")
}

// testGetRenewProjectForecast checks if route is user protected and RenewProjectForecast correctly sent back
func testGetRenewProjectForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de prévision RU, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"RenewProjectForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z","CommissionName":"Commission test",` +
				`"Value":2000000,"Project":"projet","Comment":"Essai2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":3,"ActionCode":15400202,"ActionName":"Aide à la création de logements locatifs sociaux"}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_project_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRenewProjectForecast")
}

// testGetRenewProjectForecasts checks if route is user protected and
// RenewProjectForecasts correctly sent back
func testGetRenewProjectForecasts(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"RenewProjectForecast":[{"ID":` + strconv.Itoa(ID) +
				`,"CommissionID":` + strconv.Itoa(int(c.CommissionID)) +
				`,"CommissionDate":"2018-03-01T00:00:00Z","CommissionName":"Commission test",` +
				`"Value":2000000,"Project":"projet","Comment":"Essai2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `,"ActionID":3,"ActionCode":15400202,` +
				`"ActionName":"Aide à la création de logements locatifs sociaux"}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/renew_project_forecasts").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRenewProjectForecasts")
}

// testDeleteRenewProjectForecast checks if route is user protected and renew_project_forecasts correctly sent back
func testDeleteRenewProjectForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Suppression de prévision RU, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Prévision RU supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/renew_project_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteRenewProjectForecast")
}

// testBatchRenewProjectForecasts check route is limited to admin and batch import succeeds
func testBatchRenewProjectForecasts(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"RenewProjectForecast":[{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Batch1","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Comment":"Batch2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `}]}`),
			RespContains: []string{"Batch de Prévision RUs, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : value nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"RenewProjectForecast":[{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":100,"Project":"projet","Comment":"Batch1","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `,"ActionCode":15400202},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Project":null,"Comment":"Batch2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `,"ActionCode":15400202}]}`),
			RespContains: []string{"Batch de Prévision RUs importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/renew_project_forecasts").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	if chkFactory(t, tcc, f, "BatchRenewProjectForecast") {
		response := c.E.GET("/api/renew_project_forecasts").
			WithHeader("Authorization", "Bearer "+c.Config.Users.Admin.Token).Expect()
		body := string(response.Content)
		for _, j := range []string{`"Value":100,"Project":"projet","Comment":"Batch1"`,
			`"Value":200,"Project":null,"Comment":"Batch2"`} {
			if !strings.Contains(body, j) {
				t.Errorf("BatchRenewProjectForecast[final]\n  ->attendu %s\n  ->reçu: %s", j, body)
			}
		}
	}
}
