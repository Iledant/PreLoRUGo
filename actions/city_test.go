package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
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
		{Sent: []byte(`{"City":{"InseeCode":1000000,"Name":"Essai",` +
			`"CommunityID":1,"QPV":true}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"City":{"InseeCode":0,"Name":"Essai","CommunityID":1,` +
			`"QPV":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : insee code nul
		{Sent: []byte(`{"City":{"InseeCode":100000,"Name":"","CommunityID":1,` +
			`"QPV":true}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"City":{"InseeCode":1000000,"Name":"Essai",` +
			`"CommunityID":2,"QPV":true}}`),
			Token:  c.Config.Users.Admin.Token,
			IDName: `{"InseeCode"`,
			RespContains: []string{`"City":{"InseeCode":1000000,"Name":"Essai",` +
				`"CommunityID":2,"CommunityName":"CA SAINT GERMAIN BOUCLES DE SEINE ` +
				`(78-YVELINES)","QPV":true`},
			StatusCode: http.StatusCreated}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/city").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateCity", &ID)
	return ID
}

// testUpdateCity checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateCity(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{Sent: []byte(`{"City":{"InseeCode":2000000,"Name":"Essai2",` +
			`"CommunityID":null,"QPV":false}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits administrateur requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Sent: []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{Sent: []byte(`{"City":{"InseeCode":0,"Name":"Essai2","CommunityID":null,` +
			`"QPV":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : code nul
		{Sent: []byte(`{"City":{"InseeCode":2000000,"Name":"","CommunityID":null,` +
			`"QPV":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville : Champ incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{Sent: []byte(`{"City":{"InseeCode":2000000,"Name":"Essai2",` +
			`"CommunityID":null,"QPV":false}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{Sent: []byte(`{"City":{"InseeCode":` + strconv.Itoa(ID) +
			`,"Name":"Essai2","CommunityID":3,"QPV":false}}`),
			Token: c.Config.Users.Admin.Token,
			RespContains: []string{`"City":{"InseeCode":1000000,"Name":"Essai2",` +
				`"CommunityID":3,"CommunityName":"(EX78) CC DES DEUX RIVES DE LA SEINE` +
				` (DISSOUTE AU 01/01/2016)","QPV":false`},
			StatusCode: http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/city").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateCity")
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
			RespContains: []string{`{"City":{"InseeCode":` + strconv.Itoa(ID) +
				`,"Name":"Essai2","CommunityID":3,"CommunityName":` +
				`"(EX78) CC DES DEUX RIVES DE LA SEINE ` +
				`(DISSOUTE AU 01/01/2016)","QPV":false}}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/city/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCity")
}

// testGetCities checks if route is user protected and Cities correctly sent back
func testGetCities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`{"City":[{"InseeCode":1000000,"Name":"Essai2",` +
				`"CommunityID":3,"CommunityName":"(EX78) CC DES DEUX` +
				` RIVES DE LA SEINE (DISSOUTE AU 01/01/2016)","QPV":false}]}`},
			Count:         1,
			CountItemName: `"InseeCode"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/cities").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetCities")
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
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/city/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteCity")
}

// testBatchCities check route is limited to admin and batch import succeeds
func testBatchCities(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(``),
			RespContains: []string{"Droits administrateur requis"},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"City":[{"InseeCode":75101,"Name":"","CommunityCode":` +
				`"217500016","QPV":false},
			{"InseeCode":77001,"Name":"ACHERES-LA-FORET","CommunityCode":` +
				`"247700123","QPV":true},
			{"InseeCode":78146,"Name":"CHATOU","CommunityCode":"200058519.78",` +
				`"QPV":false}]}`),
			RespContains: []string{"Batch de Villes, requête : "},
			StatusCode:   http.StatusInternalServerError}, // 1 : name empty
		{Token: c.Config.Users.Admin.Token,
			Sent: []byte(`{"City":[{"InseeCode":75101,"Name":"PARIS 1",` +
				`"CommunityCode":"217500016","QPV":false},
			{"InseeCode":77001,"Name":"ACHERES-LA-FORET","CommunityCode":` +
				`"247700123","QPV":true},
			{"InseeCode":78146,"Name":"CHATOU","CommunityCode":"200058519.78",` +
				`"QPV":false}]}`),
			Count:         3,
			CountItemName: `"InseeCode"`,
			RespContains: []string{`"InseeCode":75101,"Name":"PARIS 1"`,
				`"InseeCode":78146,"Name":"CHATOU","CommunityID":2,"CommunityName":` +
					`"CA SAINT GERMAIN BOUCLES DE SEINE (78-YVELINES)","QPV":false`},
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/cities").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "BatchCity")
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
				`"InseeCode":77001,"Name":"ACHERES-LA-FORET","QPV":true,` +
					`"CommunityID":null,"CommunityName":null`,
				//cSpell: enable
			},
			Count:         1,
			CountItemName: `"InseeCode"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/cities/paginated").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetPaginatedCities")
}
