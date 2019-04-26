package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// testHousing is the entry point for testing all renew projet requests
func testHousing(t *testing.T, c *TestContext) {
	t.Run("Housing", func(t *testing.T) {
		ID := testCreateHousing(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le logement")
			t.FailNow()
			return
		}
		testUpdateHousing(t, c, ID)
		testGetHousings(t, c)
		testDeleteHousing(t, c, ID)
		testBatchHousings(t, c)
	})
}

// testCreateHousing checks if route is admin protected and created budget action
// is properly filled
func testCreateHousing(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Housing":{"Reference":"Essai","Address":"Essai","ZipCode":1000000,"PLAI":1000000,"PLUS":1000000,"PLS":1000000,"ANRU":true,"QPV":true}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Housing":{"Reference":"","Address":"Essai","ZipCode":1000000,"PLAI":1000000,"PLUS":1000000,"PLS":1000000,"ANRU":true,"QPV":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement : Champ Reference incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{Sent: []byte(`{"Housing":{"Reference":"Essai","Address":"Essai","ZipCode":1000000,"PLAI":1000000,"PLUS":1000000,"PLS":1000000,"ANRU":true,"QPV":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Housing":{"ID":1,"Reference":"Essai","Address":"Essai","ZipCode":1000000,"PLAI":1000000,"PLUS":1000000,"PLS":1000000,"ANRU":true,"QPV":true`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/housing").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateHousing[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateHousing[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"Housing":{"ID":%d`, &ID)
		}
	}
	return ID
}

// testUpdateHousing checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateHousing(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Housing":{"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":2000000,"PLUS":2000000,"PLS":2000000,"ANRU":false,"QPV":false}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de logement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Housing":{"ID":` + strconv.Itoa(ID) + `,"Reference":"","Address":null,"ZipCode":null,"PLAI":2000000,"PLUS":2000000,"PLS":2000000,"ANRU":false,"QPV":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de logement : Champ Reference incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{Sent: []byte(`{"Housing":{"ID":0,"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":2000000,"PLUS":2000000,"PLS":2000000,"ANRU":false,"QPV":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de logement, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{Sent: []byte(`{"Housing":{"ID":` + strconv.Itoa(ID) + `,"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":2000000,"PLUS":2000000,"PLS":2000000,"ANRU":false,"QPV":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Housing":{"ID":` + strconv.Itoa(ID) + `,"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":2000000,"PLUS":2000000,"PLS":2000000,"ANRU":false,"QPV":false}`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/housing").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateHousing[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateHousing[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetHousings checks if route is user protected and Housings correctly sent back
func testGetHousings(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"Housing":[{"ID":1,"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":2000000,"PLUS":2000000,"PLS":2000000,"ANRU":false,"QPV":false}]}`},
			Count:        1,
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/housings").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetHousings[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetHousings[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetHousings[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testDeleteHousing checks if route is user protected and housings correctly sent back
func testDeleteHousing(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de logement, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Logement supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.DELETE("/api/housing/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("DeleteHousing[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("DeleteHousing[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testBatchHousings check route is limited to admin and batch import succeeds
func testBatchHousings(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Housing":[{"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":1,"PLUS":2,"PLS":3,"ANRU":false,"QPV":false},
			{"Reference":"","Address":"Adresse","ZipCode":75005,"PLAI":4,"PLUS":5,"PLS":6,"ANRU":true,"QPV":true}]}`),
			RespContains: []string{"Batch de Logements, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : validation error
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Housing":[{"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":1,"PLUS":2,"PLS":3,"ANRU":false,"QPV":false},
			{"Reference":"Essai3","Address":"Adresse","ZipCode":75005,"PLAI":4,"PLUS":5,"PLS":6,"ANRU":true,"QPV":true}]}`),
			RespContains: []string{"Batch de Logements importé"},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/housings").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchHousing[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchHousing[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			response = c.E.GET("/api/housings").
				WithHeader("Authorization", "Bearer "+tc.Token).Expect()
			body = string(response.Content)
			for _, j := range []string{`"Reference":"Essai2","Address":null,"ZipCode":null,"PLAI":1,"PLUS":2,"PLS":3,"ANRU":false,"QPV":false`,
				`"Reference":"Essai3","Address":"Adresse","ZipCode":75005,"PLAI":4,"PLUS":5,"PLS":6,"ANRU":true,"QPV":true`} {
				if !strings.Contains(body, j) {
					t.Errorf("BatchHousing[all]\n  ->attendu %s\n  ->reçu: %s", j, body)
				}
			}
		}
	}
}
