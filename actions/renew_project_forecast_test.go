package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
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
			StatusCode:   http.StatusBadRequest}, // 2 : commission ID nul
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","RenewProjectID":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : renew project ID nul
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":0,"Comment":"Essai","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de prévision RU : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : value nul
		{Sent: []byte(`{"RenewProjectForecast":{"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`"RenewProjectForecast":{"ID":1,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":1000000,"Comment":"Essai","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID))},
			StatusCode: http.StatusCreated}, // 5 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/renew_project_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateRenewProjectForecast[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateRenewProjectForecast[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"RenewProjectForecast":{"ID":%d`, &ID)
		}
	}
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
		{Sent: []byte(`{"RenewProjectForecast":{"ID":0,"CommissionID":2000000,"Value":2000000,"Comment":null,"RenewProjectID":2000000}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de prévision RU, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{Sent: []byte(`{"RenewProjectForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
			strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Comment":"Essai2","RenewProjectID":` +
			strconv.Itoa(int(c.RenewProjectID)) + "}}"),
			Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`"RenewProjectForecast":{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Comment":"Essai2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `}`},
			StatusCode: http.StatusOK}, // 6 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/renew_project_forecast").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		chkBodyStatusAndCount(t, tc, i, response, "UpdateRenewProjectForecast")
	}
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
				strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Comment":"Essai2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/renew_project_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		chkBodyStatusAndCount(t, tc, i, response, "GetRenewProjectForecast")
	}
}

// testGetRenewProjectForecasts checks if route is user protected and RenewProjectForecasts correctly sent back
func testGetRenewProjectForecasts(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"RenewProjectForecast":[{"ID":` + strconv.Itoa(ID) + `,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":2000000,"Comment":"Essai2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/renew_project_forecasts").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		chkBodyStatusAndCount(t, tc, i, response, "GetRenewProjectForecasts")
	}
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
	for i, tc := range tcc {
		response := c.E.DELETE("/api/renew_project_forecast/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		chkBodyStatusAndCount(t, tc, i, response, "DeleteRenewProjectForecast")
	}
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
				strconv.Itoa(int(c.CommissionID)) + `,"Value":100,"Comment":"Batch1","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `},{"ID":0,"CommissionID":` +
				strconv.Itoa(int(c.CommissionID)) + `,"Value":200,"Comment":"Batch2","RenewProjectID":` +
				strconv.Itoa(int(c.RenewProjectID)) + `}]}`),
			RespContains: []string{"Batch de Prévision RUs importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/renew_project_forecasts").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchRenewProjectForecast[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchRenewProjectForecast[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			response = c.E.GET("/api/renew_project_forecasts").
				WithHeader("Authorization", "Bearer "+tc.Token).Expect()
			body = string(response.Content)
			for _, j := range []string{`"Value":100,"Comment":"Batch1"`, `"Value":200,"Comment":"Batch2"`} {
				if !strings.Contains(body, j) {
					t.Errorf("BatchRenewProjectForecast[all]\n  ->attendu %s\n  ->reçu: %s", j, body)
				}
			}
		}
	}
}
