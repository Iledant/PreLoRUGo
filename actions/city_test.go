package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// testCity is the entry point for testing all renew projet requests
func testCity(t *testing.T, c *TestContext) {
	t.Run("City", func(t *testing.T) {
		ID := testCreateCity(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la ville")
			t.FailNow()
			return
		}
		testUpdateCity(t, c, ID)
		testGetCity(t, c, ID)
		testGetCities(t, c)
		testDeleteCity(t, c, ID)
		testBatchCities(t, c)
		testGetPaginatedCities(t, c)
	})
}

// testCreateCity checks if route is admin protected and created budget action
// is properly filled
func testCreateCity(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"City":{"InseeCode":1000000,"Name":"Essai","CommunityID":1}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"City":{"InseeCode":0,"Name":"Essai","CommunityID":1}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : insee code nul
		{Sent: []byte(`{"City":{"InseeCode":100000,"Name":"","CommunityID":1}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"City":{"InseeCode":1000000,"Name":"Essai","CommunityID":2}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"City":{"InseeCode":1000000,"Name":"Essai","CommunityID":2`},
			StatusCode:   http.StatusCreated}, // 4 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/city").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("CreateCity[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("CreateCity[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if tc.StatusCode == http.StatusCreated {
			fmt.Sscanf(body, `{"City":{"InseeCode":%d`, &ID)
		}
	}
	return ID
}

// testUpdateCity checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateCity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"City":{"InseeCode":2000000,"Name":"Essai2","CommunityID":null}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"City":{"InseeCode":0,"Name":"Essai2","CommunityID":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code nul
		{Sent: []byte(`{"City":{"InseeCode":2000000,"Name":"","CommunityID":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"City":{"InseeCode":2000000,"Name":"Essai2","CommunityID":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{Sent: []byte(`{"City":{"InseeCode":` + strconv.Itoa(ID) + `,"Name":"Essai2","CommunityID":null}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"City":{"InseeCode":` + strconv.Itoa(ID) + `,"Name":"Essai2","CommunityID":null}`},
			StatusCode:   http.StatusOK}, // 5 : ok
	}
	for i, tc := range tcc {
		response := c.E.PUT("/api/city").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("UpdateCity[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("UpdateCity[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetCity checks if route is user protected and City correctly sent back
func testGetCity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de ville, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"City":{"InseeCode":` + strconv.Itoa(ID) + `,"Name":"Essai2","CommunityID":null}}`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/city/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetCity[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetCity[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testGetCities checks if route is user protected and Cities correctly sent back
func testGetCities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"City":[{"InseeCode":1000000,"Name":"Essai2","CommunityID":null}]}`},
			Count:        1,
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/cities").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetCities[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetCities[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"InseeCode"`)
			if count != tc.Count {
				t.Errorf("GetCities[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testDeleteCity checks if route is user protected and cities correctly sent back
func testDeleteCity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression de ville, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Ville supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.DELETE("/api/city/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("DeleteCity[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("DeleteCity[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
	}
}

// testBatchCities check route is limited to admin and batch import succeeds
func testBatchCities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"City":[{"InseeCode":75101,"Name":"","CommunityCode":"217500016"},
			{"InseeCode":77001,"Name":"ACHERES-LA-FORET","CommunityCode":"247700123"},
			{"InseeCode":78146,"Name":"CHATOU","CommunityCode":"200058519.78"}]}`),
			RespContains: []string{"Batch de Villes, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : name empty
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"City":[{"InseeCode":75101,"Name":"PARIS 1","CommunityCode":"217500016"},
			{"InseeCode":77001,"Name":"ACHERES-LA-FORET","CommunityCode":"247700123"},
			{"InseeCode":78146,"Name":"CHATOU","CommunityCode":"200058519.78"}]}`),
			Count:        3,
			RespContains: []string{`"InseeCode":75101,"Name":"PARIS 1"`, `"InseeCode":78146,"Name":"CHATOU","CommunityID":4`},
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	for i, tc := range tcc {
		response := c.E.POST("/api/cities").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("BatchCity[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("BatchCity[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"InseeCode"`)
			if count != tc.Count {
				t.Errorf("BatchCity[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}

// testGetPaginatedCities checks if route is user protected and Cities correctly sent back
func testGetPaginatedCities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Page=2&Search=acheres`),
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Page=2&Search=acheres`),
			RespContains: []string{`"City"`, `"Page"`, `"ItemsCount"`,
				// cSpell: disable
				`"InseeCode":77001,"Name":"ACHERES-LA-FORET","CommunityID":null`,
				//cSpell: enable
			},
			Count:      1,
			StatusCode: http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/cities/paginated").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetPaginatedCities[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetPaginatedCities[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"InseeCode"`)
			if count != tc.Count {
				t.Errorf("GetPaginatedCities[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
