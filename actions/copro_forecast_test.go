package actions

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCoproForecast is the entry point for testing all renew projet requests
func testCoproForecast(t *testing.T, c *TestContext) {
	t.Run("CoproForecast", func(t *testing.T) {
		ID := testCreateCoproForecast(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la prévision copro")
			t.FailNow()
			return
		}
		testUpdateCoproForecast(t, c, ID)
		testGetCoproForecast(t, c, ID)
		testGetCoproForecasts(t, c, ID)
		testDeleteCoproForecast(t, c, ID)
		testBatchCoproForecasts(t, c)
	})
}

// testCreateCoproForecast checks if route is admin protected and created copro forecast
// is properly filled
func testCreateCoproForecast(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"CoproForecast":{"CommissionID":0,"Value":1000000,"Comment":"Essai","CoproID":1000000}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les copropriétés requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création de prévision copro, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"CoproForecast":{"CommissionID":0,"Value":1000000,"Comment":"Essai","CoproID":` +
			strconv.Itoa(int(c.CoproID)) + "}}"),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création de prévision copro : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : commission ID nul
		{Sent: []byte(`{"CoproForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","CoproID":0}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création de prévision copro : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : renew project ID nul
		{Sent: []byte(`{"CoproForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai","CoproID":` +
			strconv.Itoa(int(c.CoproID)) + "}}"),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création de prévision copro : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : value nul
		{Sent: []byte(`{"CoproForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","CoproID":` +
			strconv.Itoa(int(c.CoproID)) + "}}"),
			Token: c.Config.Users.CoproUser.Token,
			RespContains: []string{`"CoproForecast":{"ID":1,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":1000000,"Comment":"Essai","CoproID":` +
				strconv.Itoa(int(c.CoproID))},
			StatusCode: http.StatusCreated}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateCoproForecast", &ID)
	return ID
}

// testUpdateCoproForecast checks if route is admin protected and updated copro forecast
// is properly filled
func testUpdateCoproForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"CoproForecast":{"CommissionID":2000000,"Value":2000000,"Comment":null,"CoproID":2000000}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les copropriétés requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de prévision copro, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"CoproForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","CoproID":` +
			strconv.Itoa(int(c.CoproID)) + "}}"),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de prévision copro : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : value nul
		{Sent: []byte(`{"CoproForecast":{"CommissionID":0,"Value":0,"Comment":"Essai2","CoproID":` +
			strconv.Itoa(int(c.CoproID)) + "}}"),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de prévision copro : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : commission ID nul
		{Sent: []byte(`{"CoproForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","CoproID":0}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de prévision copro : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : renew project ID nul
		{Sent: []byte(`{"CoproForecast":{"ID":0,"CommissionID":2000000,"Value":2000000,"Comment":null,"CoproID":2000000}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification de prévision copro, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{Sent: []byte(`{"CoproForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Comment":"Essai2","CoproID":` +
			strconv.Itoa(int(c.CoproID)) + "}}"),
			Token: c.Config.Users.CoproUser.Token,
			RespContains: []string{`"CoproForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":2000000,"Comment":"Essai2","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `}`},
			StatusCode: http.StatusOK}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/copro_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateCoproForecast")
}

// testGetCoproForecast checks if route is user protected and copro forecast correctly sent back
func testGetCoproForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de prévision copro, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"CoproForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":2000000,"Comment":"Essai2","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoproForecast")
}

// testGetCoproForecasts checks if route is user protected and CoproForecasts correctly sent back
func testGetCoproForecasts(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"CoproForecast":[{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":2000000,"Comment":"Essai2","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `}]}`},
			Count:      1,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro_forecasts").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCoproForecasts")
}

// testDeleteCoproForecast checks if route is user protected and CoproForecasts correctly sent back
func testDeleteCoproForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les copropriétés requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.CoproUser.Token,
			RespContains: []string{`Suppression de prévision copro, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.CoproUser.Token,
			RespContains: []string{`Prévision copro supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/copro_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteCoproForecast")
}

// testBatchCoproForecasts check route is limited to admin and batch import succeeds
func testBatchCoproForecasts(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"CoproForecast":[{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Batch1","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Comment":"Batch2","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `}]}`),
			RespContains: []string{"Batch de Prévision copros, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : value nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"CoproForecast":[{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":100,"Comment":"Batch1","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Comment":"Batch2","CoproID":` +
				strconv.Itoa(int(c.CoproID)) + `}]}`),
			RespContains: []string{"Batch de Prévision copros importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/copro_forecasts").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchCoproForecast[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchCoproForecast[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			response = c.E.GET("/api/copro_forecasts").
				WithHeader("Authorization", "Bearer "+tc.Token).Expect()
			body = string(response.Content)
			for _, j := range []string{`"Value":100,"Comment":"Batch1"`, `"Value":200,"Comment":"Batch2"`} {
				if !strings.Contains(body, j) {
					t.Errorf("BatchCoproForecast[all]\n  ->attendu %s\n  ->reçu: %s", j, body)
				}
			}
		}
	}
}
