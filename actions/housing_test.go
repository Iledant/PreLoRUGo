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
		// testGetHousings(t, c)
		// testDeleteHousing(t, c, ID)
	})
}

// testCreateHousing checks if route is admin protected and created budget action
// is properly filled
func testCreateHousing(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Housing":{"Reference":"LLS001","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Housing":{"Reference":"","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement : Champ reference incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{Sent: []byte(`{"Housing":{"Reference":"LLS001","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Housing":{"ID":1,"Reference":"LLS001","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}`},
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
		{Sent: []byte(`{"Housing":{"Reference":"LLS001","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Housing":{"Reference":"","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement : Champ reference incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : reference empty
		{Sent: []byte(`{"Housing":{"ID":0,"Reference":"LLS001","Address":"Adresse","ZipCode":75001,"PLAI":3,"PLUS":5,"PLS":7,"ANRU":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de logement, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 3 : bad ID
		{Sent: []byte(`{"Housing":{"ID":` + strconv.Itoa(ID) + `,"Reference":"LLS002","Address":"Adresse2","ZipCode":75002,"PLAI":4,"PLUS":6,"PLS":8,"ANRU":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Housing":{"ID":` + strconv.Itoa(ID) + `,"Reference":"LLS002","Address":"Adresse2","ZipCode":75002,"PLAI":4,"PLUS":6,"PLS":8,"ANRU":false}`},
			StatusCode:   http.StatusCreated}, // 4 : ok
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
