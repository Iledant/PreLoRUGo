package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testHousingComment is the entry point for testing all housing comments
func testHousingComment(t *testing.T, c *TestContext) {
	t.Run("HousingComment", func(t *testing.T) {
		ID := testCreateHousingComment(t, c)
		if ID == 0 {
			t.Error("Impossible de créer la commentaire de logement")
			t.FailNow()
			return
		}
		testUpdateHousingComment(t, c, ID)
		testGetHousingComments(t, c)
		testDeleteHousingComment(t, c, ID)
	})
}

// testCreateHousingComment checks if route is admin protected and created housing
// comment is properly filled
func testCreateHousingComment(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de commentaire de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingComment":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de commentaire de logement, paramètre :`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"HousingComment":{"Name":"type de commentaire"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"HousingComment"`, `"Name":"type de commentaire"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/housing_comment").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingComment", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{{
		Sent:         []byte(`{"HousingComment":{"Name":"type de commentaire test"}}`),
		Token:        c.Config.Users.Admin.Token,
		IDName:       `"ID"`,
		RespContains: []string{`"HousingComment"`, `"Name":"type de commentaire test"`},
		StatusCode:   http.StatusCreated}, // created to store in config
	}
	for _, r := range chkFactory(tcc, f, "CreateHousingCommentForTests", &c.HousingCommentID) {
		t.Error(r)
	}
	return ID
}

// testUpdateHousingComment checks if route is admin protected and housing comment
// is properly filled
func testUpdateHousingComment(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commentaire de logement, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"HousingComment":{"ID":0,"Name":"type modifié de commentaire"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commentaire de logement, requête : Commentaire introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"HousingComment":{"ID":` + strconv.Itoa(ID) + `,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de commentaire de logement, paramètre : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"HousingComment":{"ID":` + strconv.Itoa(ID) + `,"Name":"type modifié de commentaire"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"HousingComment"`, `"Name":"type modifié de commentaire"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/housing_comment").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateHousingComment") {
		t.Error(r)
	}
}

// testGetHousingComments checks route is protected and all housing conventions
// are correctly sent back
func testGetHousingComments(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"HousingComment"`, `"Name"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/housing_comments").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetHousingComments") {
		t.Error(r)
	}
}

// testDeleteHousingComment checks that route is renew project protected and
// delete request sends ok back
func testDeleteHousingComment(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de commentaire de logement, requête : Commentaire introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Commentaire de logement supprimée`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/housing_comment/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteHousingComment") {
		t.Error(r)
	}
}
