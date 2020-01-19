package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testCoproDoc is the entry point for testing all renew projet requests
func testCoproDoc(t *testing.T, c *TestContext) {
	t.Run("CoproDoc", func(t *testing.T) {
		ID := testCreateCoproDoc(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le document de copro")
			t.FailNow()
			return
		}
		testGetCoproDocs(t, c)
		testUpdateCoproDoc(t, c, ID)
		testDeleteCoproDoc(t, c, ID)
	})
}

// testCreateCoproDoc checks if route is admin protected and created budget action
// is properly filled
func testCreateCoproDoc(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création d'un document copro, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"CoproDoc":{"Name":null,"Link":"lien de document"}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création d'un document copro, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 2 : name null
		{
			Sent:         []byte(`{"CoproDoc":{"Name":"nom de document","Link":null}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Création d'un document copro, requête :`},
			StatusCode:   http.StatusInternalServerError}, // 3 : link null
		{
			Sent:   []byte(`{"CoproDoc":{"Name":"nom de document","Link":"lien de document"}}`),
			Token:  c.Config.Users.CoproUser.Token,
			IDName: `{"ID"`,
			RespContains: []string{`"CoproDoc":{"ID":1,"CoproID":` +
				strconv.FormatInt(c.CoproID, 10) + `,"Name":"nom de document","Link":"lien de document"}`},
			StatusCode: http.StatusCreated}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/copro/"+strconv.FormatInt(c.CoproID, 10)+"/copro_doc").
			WithBytes(tc.Sent).WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateCoproDoc", &ID) {
		t.Error(r)
	}
	return ID
}

// testUpdateCoproDoc checks if route is admin protected and Updated budget action
// is properly filled
func testUpdateCoproDoc(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'un document copro, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent: []byte(`{"CoproDoc":{"InseeCode":0,"Name":"Essai2","CommunityID":null,` +
				`"QPV":false}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'un document copro, requête : link vide`},
			StatusCode:   http.StatusInternalServerError}, // 2 : code nul
		{
			Sent: []byte(`{"CoproDoc":{"InseeCode":2000000,"Name":"","CommunityID":null,` +
				`"QPV":false}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'un document copro, requête : nom vide`},
			StatusCode:   http.StatusInternalServerError}, // 3 : name empty
		{
			Sent: []byte(`{"CoproDoc":{"InseeCode":2000000,"Name":"Essai2",` +
				`"CommunityID":null,"QPV":false}}`),
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Modification d'un document copro, requête : `},
			StatusCode:   http.StatusInternalServerError}, // 4 : bad ID
		{
			Sent: []byte(`{"CoproDoc":{"ID":` + strconv.Itoa(ID) +
				`,"CoproID":` + strconv.FormatInt(c.CoproID, 10) + `,"Name":"nom2 de doc","Link":"lien2 de doc"}}`),
			Token: c.Config.Users.CoproUser.Token,
			RespContains: []string{`"CoproDoc":{"ID":` + strconv.Itoa(ID) +
				`,"CoproID":` + strconv.FormatInt(c.CoproID, 10) + `,"Name":"nom2 de doc","Link":"lien2 de doc"}}`},
			StatusCode: http.StatusOK}, // 5 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/copro/"+strconv.FormatInt(c.CoproID, 10)+"/copro_doc").
			WithBytes(tc.Sent).WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateCoproDoc") {
		t.Error(r)
	}
}

// testDeleteCoproDoc checks if route is user protected and cities correctly sent back
func testDeleteCoproDoc(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.CoproCheckTestCase, // 0 : user token
		{
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Suppression d'un document copro, requête : `},
			ID:           0,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad ID
		{
			Token:        c.Config.Users.CoproUser.Token,
			RespContains: []string{`Document supprimé`},
			ID:           ID,
			StatusCode:   http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/copro/"+strconv.FormatInt(c.CoproID, 10)+"/copro_doc/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteCoproDoc") {
		t.Error(r)
	}
}

// testGetCoproDocs checks if route is user protected and CoproDocs correctly sent back
func testGetCoproDocs(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : token empty
		{
			Token:        c.Config.Users.User.Token,
			RespContains: []string{`Documents d'une copro, erreur CoproID : `},
			Count:        1,
			Params:       "a",
			StatusCode:   http.StatusBadRequest}, // 1 : bad CoproID
		{
			Token: c.Config.Users.User.Token,
			RespContains: []string{`"CoproDoc"`, `"Name":"nom de document"`,
				`"Link":"lien de document"`},
			Count:         1,
			CountItemName: `"ID"`,
			Params:        strconv.FormatInt(c.CoproID, 10),
			StatusCode:    http.StatusOK}, // 2 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/copro/"+tc.Params+"/copro_docs").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetCoproDocs") {
		t.Error(r)
	}
}
