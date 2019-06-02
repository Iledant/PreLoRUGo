package actions

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingForecast is the entry point for testing all renew projet requests
func testHousingForecast(t *testing.T, c *TestContext) {
	t.Run("HousingForecast", func(t *testing.T) {
		ID := testCreateHousingForecast(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la prévision logement")
			t.FailNow()
			return
		}
		testUpdateHousingForecast(t, c, ID)
		testGetHousingForecast(t, c, ID)
		testGetHousingForecasts(t, c, ID)
		testDeleteHousingForecast(t, c, ID)
		testBatchHousingForecasts(t, c)
	})
}

// testCreateHousingForecast checks if route is admin protected and created housing forecast
// is properly filled
func testCreateHousingForecast(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"HousingForecast":{"CommissionID":0,"Value":1000000,"Comment":"Essai","HousingID":1000000}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets logement requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Création de prévision logement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"HousingForecast":{"CommissionID":0,"Value":1000000,"Comment":"Essai","HousingID":` +
			strconv.Itoa(int(c.HousingID)) + "}}"),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Création de prévision logement : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : commission ID nul
		{Sent: []byte(`{"HousingForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","HousingID":0}}`),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Création de prévision logement : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : renew project ID nul
		{Sent: []byte(`{"HousingForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai","HousingID":` +
			strconv.Itoa(int(c.HousingID)) + "}}"),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Création de prévision logement : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : value nul
		{Sent: []byte(`{"HousingForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","HousingID":` +
			strconv.Itoa(int(c.HousingID)) + "}}"),
			Token:  c.Config.Users.HousingUser.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"HousingForecast":{"ID":1,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":1000000,"Comment":"Essai","HousingID":` +
				strconv.Itoa(int(c.HousingID))},
			StatusCode: http.StatusCreated}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateHousingForecast", &ID)
	return ID
}

// testUpdateHousingForecast checks if route is admin protected and updated housing forecast
// is properly filled
func testUpdateHousingForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"HousingForecast":{"CommissionID":2000000,"Value":2000000,"Comment":null,"HousingID":2000000}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets logement requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Modification de prévision logement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"HousingForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","HousingID":` +
			strconv.Itoa(int(c.HousingID)) + "}}"),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Modification de prévision logement : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : value nul
		{Sent: []byte(`{"HousingForecast":{"CommissionID":0,"Value":0,"Comment":"Essai2","HousingID":` +
			strconv.Itoa(int(c.HousingID)) + "}}"),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Modification de prévision logement : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : commission ID nul
		{Sent: []byte(`{"HousingForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai2","HousingID":0}}`),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Modification de prévision logement : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : renew project ID nul
		{Sent: []byte(`{"HousingForecast":{"ID":0,"CommissionID":2000000,"Value":2000000,"Comment":null,"HousingID":2000000}}`),
			Token:        c.Config.Users.HousingUser.Token,
			RespContains: []string{`Modification de prévision logement, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{Sent: []byte(`{"HousingForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Comment":"Essai2","HousingID":` +
			strconv.Itoa(int(c.HousingID)) + "}}"),
			Token: c.Config.Users.HousingUser.Token,
			RespContains: []string{`"HousingForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":2000000,"Comment":"Essai2","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `}`},
			StatusCode: http.StatusOK}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/housing_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateHousingForecast")
}

// testGetHousingForecast checks if route is user protected and housing forecast correctly sent back
func testGetHousingForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de prévision logement, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"HousingForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":2000000,"Comment":"Essai2","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetHousingForecast")
}

// testGetHousingForecasts checks if route is user protected and HousingForecasts correctly sent back
func testGetHousingForecasts(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"HousingForecast":[{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"CommissionDate":"2018-03-01T00:00:00Z",` +
				`"CommissionName":"Commission test","Value":2000000,"Comment":"Essai2","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `}]}`},
			Count:      1,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_forecasts").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetHousingForecasts")
}

// testDeleteHousingForecast checks if route is user protected and HousingForecasts correctly sent back
func testDeleteHousingForecast(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets logement requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.HousingUser.Token,
			RespContains: []string{`Suppression de prévision logement, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.HousingUser.Token,
			RespContains: []string{`Prévision logement supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/housing_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteHousingForecast")
}

// testBatchHousingForecasts check route is limited to admin and batch import succeeds
func testBatchHousingForecasts(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"HousingForecast":[{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Batch1","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Comment":"Batch2","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `}]}`),
			RespContains: []string{"Batch de Prévision logements, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : value nul
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"HousingForecast":[{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":100,"Comment":"Batch1","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Comment":"Batch2","HousingID":` +
				strconv.Itoa(int(c.HousingID)) + `}]}`),
			RespContains: []string{"Batch de Prévision logements importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_forecasts").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	if chkFactory(t, tcc, f, "BatchHousingForecast") {
		response := c.E.GET("/api/housing_forecasts").
			WithHeader("Authorization", "Bearer "+c.Config.Users.Admin.Token).Expect()
		body := string(response.Content)
		for _, j := range []string{`"Value":100,"Comment":"Batch1"`, `"Value":200,"Comment":"Batch2"`} {
			if !strings.Contains(body, j) {
				t.Errorf("BatchHousingForecast[all]\n  ->attendu %s\n  ->reçu: %s", j, body)
			}
		}
	}
}
