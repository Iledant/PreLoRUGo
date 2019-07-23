package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testRPCmtCiyJoin is the entry point for testing all renew projets commitments
// and cities joints requests
func testRPCmtCiyJoin(t *testing.T, c *TestContext) {
	t.Run("RPCmtCiyJoin", func(t *testing.T) {
		ID := testCreateRPCmtCiyJoin(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la liaison ville engagement RU")
			t.FailNow()
			return
		}
		testUpdateRPCmtCiyJoin(t, c, ID)
		testGetRPCmtCiyJoin(t, c, ID)
		testGetRPCmtCiyJoins(t, c)
		testDeleteRPCmtCiyJoin(t, c, ID)
	})
}

// testCreateRPCmtCiyJoin checks if route is admin protected and created budget action
// is properly filled
func testCreateRPCmtCiyJoin(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"CommitmentID":3,"CityCode":75101}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized!
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"CommitmentID":0,"CityCode":75101}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : CommitmentID empty
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"CommitmentID":3,"CityCode":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : CityCode empty
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"CommitmentID":3,"CityCode":75000}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Création de lien engagement ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : CityCode doesn't exist
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"CommitmentID":3,"CityCode":75101}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"RPCmtCiyJoin":{"ID":2,"CommitmentID":3,"CityCode":75101`},
			StatusCode:   http.StatusCreated}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/rp_cmt_city_join").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "CreateRPCmtCiyJoin", &ID)
	return ID
}

// testUpdateRPCmtCiyJoin checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateRPCmtCiyJoin(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":77001}}`),
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Droits sur les projets RU requis`},
			StatusCode:   http.StatusUnauthorized}, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville, décodage :`},
			StatusCode:   http.StatusInternalServerError}, // 1 : bad request
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":0,"CityCode":77001}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 2 : commitment ID null
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":0}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville : Champ CommitmentID ou CityCode incorrect`},
			StatusCode:   http.StatusBadRequest}, // 3 : city code null
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"ID":0,"CommitmentID":4,"CityCode":77001}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{
			Sent:         []byte(`{"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":77000}}`),
			Token:        c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`Modification de lien engagement ville, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 5 : bad city code
		{
			Sent:  []byte(`{"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) + `,"CommitmentID":4,"CityCode":77001}}`),
			Token: c.Config.Users.RenewProjectUser.Token,
			RespContains: []string{`"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) +
				`,"CommitmentID":4,"CityCode":77001}`},
			StatusCode: http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/rp_cmt_city_join").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "UpdateRPCmtCiyJoin")
}

// testGetRPCmtCiyJoin checks if route is user protected and RPCmtCiyJoin correctly sent back
func testGetRPCmtCiyJoin(t *testing.T, c *TestContext, ID int) {
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
			RespContains: []string{`{"RPCmtCiyJoin":{"ID":` + strconv.Itoa(ID) +
				`,"CommitmentID":4,"CityCode":77001}`},
			ID:         ID,
			StatusCode: http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_cmt_city_join/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPCmtCiyJoin")
}

// testGetRPCmtCiyJoins checks if route is user protected and RPCmtCiyJoins correctly sent back
func testGetRPCmtCiyJoins(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{
			Token:        "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`{"RPCmtCiyJoin":[{"ID":2,"CommitmentID":4,"CityCode":77001}]}`},
			Count:         1,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/rp_cmt_city_joins").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	chkFactory(t, tcc, f, "GetRPCmtCiyJoins")
}

// testDeleteRPCmtCiyJoin checks if route is user protected and rp_cmt_city_joins correctly sent back
func testDeleteRPCmtCiyJoin(t *testing.T, c *TestContext, ID int) {
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
	chkFactory(t, tcc, f, "DeleteRPCmtCiyJoin")
}
