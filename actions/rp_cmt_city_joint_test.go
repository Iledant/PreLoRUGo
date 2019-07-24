package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPCmtCityJoin is the entry point for testing all renew projets commitments
// and cities joints requests
func testRPCmtCityJoin(t *testing.T, c *TestContext) {
	t.Run("RPCmtCityJoin", func(t *testing.T) {
		ID := testCreateRPCmtCityJoin(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la liaison ville engagement RU")
			t.FailNow()
			return
		}
		testUpdateRPCmtCityJoin(t, c, ID)
		testGetRPCmtCityJoin(t, c, ID)
		testGetRPCmtCityJoins(t, c)
		testDeleteRPCmtCityJoin(t, c, ID)
	})
}

// testCreateRPCmtCityJoin checks if route is admin protected and created budget action
// is properly filled
func testCreateRPCmtCityJoin(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"CommitmentID":3,"CityCode":75101}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized!
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"CommitmentID":0,"CityCode":75101}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : CommitmentID empty
		{
			Sent:         []byte(`{"CommitmentID":3,"CityCode":0}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : CityCode empty
		{
			Sent:         []byte(`{"CommitmentID":3,"CityCode":75000}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : CityCode doesn't exist
		{
			Sent:         []byte(`{"CommitmentID":3,"CityCode":75101}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"RPCmtCityJoin":{"ID":2,"CommitmentID":3,"CityCode":75101`},
			StatusCode:   http.StatusCreated}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/rp_cmt_city_join").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateRPCmtCityJoin", &ID)
	return ID
}

// testUpdateRPCmtCityJoin checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateRPCmtCityJoin(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":77001}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":0,"CityCode":77001}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : commitment ID null
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":0}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : city code null
		{
			Sent:         []byte(`{"ID":0,"CommitmentID":4,"CityCode":77001}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{
			Sent:         []byte(`{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":77000}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad city code
		{
			Sent:  []byte(`{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":77001}`),
			Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`"RPCmtCityJoin":{"ID":` + strconv.Itoa(ID) +
				`,"CommitmentID":4,"CityCode":77001}`},
			StatusCode: http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/rp_cmt_city_join").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateRPCmtCityJoin")
}

// testGetRPCmtCityJoin checks if route is user protected and RPCmtCityJoin correctly sent back
func testGetRPCmtCityJoin(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Token:        "",
			RespContains: []string{`Token absent`},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			StatusCode:   http.StatusInternalServerError,
			RespContains: []string{`Récupération de lien engagement ville, requête :`},
			ID:           0}, // 1 : bad ID
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`{"RPCmtCityJoin":{"ID":` + strconv.Itoa(ID) +
				`,"CommitmentID":4,"CityCode":77001}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_cmt_city_join/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPCmtCityJoin")
}

// testGetRPCmtCityJoins checks if route is user protected and RPCmtCityJoins correctly sent back
func testGetRPCmtCityJoins(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{
			Token:        "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`{"RPCmtCityJoin":[{"ID":2,"CommitmentID":4,"CityCode":77001}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_cmt_city_joins").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPCmtCityJoins")
}

// testDeleteRPCmtCityJoin checks if route is user protected and rp_cmt_city_joins correctly sent back
func testDeleteRPCmtCityJoin(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user token
		{
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Suppression de lien engagement ville, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Lien engagement ville supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/rp_cmt_city_join/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "DeleteRPCmtCityJoin")
}
