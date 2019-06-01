package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCommunity is the entry point for testing all renew projet requests
func testCommunity(t *testing.T, c *TestContext) {
	t.Run("Community", func(t *testing.T) {
		ID := testCreateCommunity(t, c)
		if ID == 0 {
			t.Error("Impossible de créer l'interco")
			t.FailNow()
			return
		}
		testUpdateCommunity(t, c, ID)
		testGetCommunity(t, c, ID)
		testGetCommunities(t, c)
		testDeleteCommunity(t, c, ID)
		testBatchCommunities(t, c)
	})
}

// testCreateCommunity checks if route is admin protected and created budget action
// is properly filled
func testCreateCommunity(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Community":{"Code":"Essai","Name":"Essai"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'interco, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Community":{"Code":"","Name":"Essai"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'interco : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code empty
		{Sent: []byte(`{"Community":{"Code":"Essai","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création d'interco : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"Community":{"Code":"Essai","Name":"Essai"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `{"ID"`,
			RespContains: []string{`"Community":{"ID":1,"Code":"Essai","Name":"Essai"`},
			StatusCode:   http.StatusCreated}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/community").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateCommunity", &ID)
	return ID
}

// testUpdateCommunity checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateCommunity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"Community":{"Code":"Essai2","Name":"Essai2"}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'interco, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"Community":{"ID":` + strconv.Itoa(ID) + `,"Code":"","Name":"Essai2"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'interco : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : code empty
		{Sent: []byte(`{"Community":{"ID":` + strconv.Itoa(ID) + `,"Code":"Essai2","Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'interco : Champ code ou name incorrect`},
			StatusCode:   http.StatusBadRequest}, // 4 : name empty
		{Sent: []byte(`{"Community":{"ID":0,"Code":"Essai2","Name":"Essai2"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification d'interco, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad ID
		{Sent: []byte(`{"Community":{"ID":` + strconv.Itoa(ID) + `,"Code":"Essai2","Name":"Essai2"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"Community":{"ID":` + strconv.Itoa(ID) + `,"Code":"Essai2","Name":"Essai2"}`},
			StatusCode:   http.StatusOK}, // 6 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/community").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateCommunity")
}

// testGetCommunity checks if route is user protected and Community correctly sent back
func testGetCommunity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération d'interco, requête :`},
			ID:           0}, // 1 : bad ID
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"Community":{"ID":` + strconv.Itoa(ID) + `,"Code":"Essai2","Name":"Essai2"}}`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/community/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCommunity")
}

// testGetCommunities checks if route is user protected and Communities correctly sent back
func testGetCommunities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains:  []string{`{"Community":[{"ID":1,"Code":"Essai2","Name":"Essai2"}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/communities").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCommunities")
}

// testDeleteCommunity checks if route is user protected and communities correctly sent back
func testDeleteCommunity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Suppression d'interco, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{Token: c.Config.Users.Admin.Token,
			RespContains: []string{`Interco supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/community/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteCommunity")
}

// testBatchCommunities check route is limited to admin and batch import succeeds
func testBatchCommunities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Community":[{"Code":"200000321","Name":"(EX78) CC DES DEUX RIVES DE LA SEINE (DISSOUTE AU 01/01/2016)"},
			{"Code":"","Name":"VILLE DE PARIS (EPT1)"},{"Code":"200058519.78","Name":"CA SAINT GERMAIN BOUCLES DE SEINE (78-YVELINES)"}]}`),
			RespContains: []string{"Batch de Intercos, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : code empty
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"Community":[{"Code":"200000321","Name":"(EX78) CC DES DEUX RIVES DE LA SEINE (DISSOUTE AU 01/01/2016)"},
			{"Code":"217500016","Name":"VILLE DE PARIS (EPT1)"},{"Code":"200058519.78","Name":"CA SAINT GERMAIN BOUCLES DE SEINE (78-YVELINES)"}]}`),
			Count:         3,
			CountItemName: `"ID"`,
			RespContains:  []string{"Community", `"Code":"200000321","Name":"(EX78) CC DES DEUX RIVES DE LA SEINE (DISSOUTE AU 01/01/2016)"`},
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/communities").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchCommunity")
}
