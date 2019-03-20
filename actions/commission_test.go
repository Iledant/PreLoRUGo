package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// testCommission is the entry point for testing all renew projet requests
func testCommission(t *testing.T, c *TestContext) {
	t.Run("Commission", func(t *testing.T) {
		ID := testCreateCommission(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la commission")
			t.FailNow()
			return
		}
		testUpdateCommission(t, c, ID)
		testGetCommission(t, c, ID)
		testGetCommissions(t, c)
		testDeleteCommission(t, c, ID)
	})
}

// testCreateCommission checks if route is admin protected and created budget action
// is properly filled
func testCreateCommission(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Commission":{"Name":"Essai","Date":2019-03-01T00:00:00}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de commission, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Commission":{}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de commission : Champ name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{Sent: []byte(`{"Commission":{"Name":"Essai","Date":"2019-03-01T00:00:00Z"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Commission":{"ID":1,"Name":"Essai","Date":"2019-03-01T00:00:00Z"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/commission").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateCommission[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateCommission[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"Commission":{"ID":%d`, &ID)
		}
	}
	return ID
}

// testUpdateCommission checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateCommission(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Commission":{"Name":"Essai2","Date":null}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commission, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Commission":{}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commission : Champ name vide`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{Sent: []byte(`{"Commission":{"ID":0,"Name":"Essai2","Date":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commission, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{Sent: []byte(`{"Commission":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","Date":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Commission":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","Date":null}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/commission").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateCommission[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateCommission[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetCommission checks if route is user protected and Commission correctly sent back
func testGetCommission(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de commission, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"Commission":{"ID":` + strconv.Itoa(ID) + `,"Name":"Essai2","Date":null}}`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/commission/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetCommission[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetCommission[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetCommissions checks if route is user protected and Commissions correctly sent back
func testGetCommissions(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"Commission":[{"ID":1,"Name":"Essai2","Date":null}]}`},
			Count:        1,
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/commissions").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetCommissions[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetCommissions[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetCommissions[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testDeleteCommission checks if route is user protected and commissions correctly sent back
func testDeleteCommission(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de commission, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Commission supprimée`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.DELETE("/api/commission/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("DeleteCommission[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("DeleteCommission[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}
