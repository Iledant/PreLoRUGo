package actions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/iris-contrib/httpexpect"
)

// testConventionType is the entry point for testing all convention types
func testConventionType(t *testing.T, c *TestContext) {
	t.Run("ConventionType", func(t *testing.T) {
		ID := testCreateConventionType(t, c)
		if ID == 0 {
			t.Error("Impossible de créer le type de convention")
			t.FailNow()
			return
		}
		testUpdateConventionType(t, c, ID)
		testGetConventionTypes(t, c)
		testDeleteConventionType(t, c, ID)
	})
}

// testCreateConventionType checks if route is admin protected and created convention
// type is properly filled
func testCreateConventionType(t *testing.T, c *TestContext) (ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de type de convention, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"ConventionType":{"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Création de type de convention, paramètre :`},
			StatusCode:   http.StatusBadRequest}, // 2 : name empty
		{
			Sent:         []byte(`{"ConventionType":{"Name":"PLUS"}}`),
			Token:        c.Config.Users.Admin.Token,
			IDName:       `"ID"`,
			RespContains: []string{`"ConventionType"`, `"Name":"PLUS"`},
			StatusCode:   http.StatusCreated}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.POST("/api/convention_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "CreateConventionType", &ID) {
		t.Error(r)
	}
	tcc = []TestCase{{
		Sent:         []byte(`{"ConventionType":{"Name":"PLAI"}}`),
		Token:        c.Config.Users.Admin.Token,
		IDName:       `"ID"`,
		RespContains: []string{`"ConventionType"`, `"Name":"PLAI"`},
		StatusCode:   http.StatusCreated}, // created to store in config
	}
	for _, r := range chkFactory(tcc, f, "CreateConventionTypeForTests", &c.ConventionTypeID) {
		t.Error(r)
	}
	return ID
}

// testUpdateConventionType checks if route is admin protected and convention type
// is properly filled
func testUpdateConventionType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : user unauthorized
		{
			Sent:         []byte(`fake`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de type de convention, décodage :`},
			StatusCode:   http.StatusBadRequest}, // 1 : bad request
		{
			Sent:         []byte(`{"ConventionType":{"ID":0,"Name":"PLS"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de type de convention, requête : Type de convention introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Sent:         []byte(`{"ConventionType":{"ID":` + strconv.Itoa(ID) + `,"Name":""}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`Modification de type de convention, paramètre : Nom vide`},
			StatusCode:   http.StatusBadRequest}, // 3 : name empty
		{
			Sent:         []byte(`{"ConventionType":{"ID":` + strconv.Itoa(ID) + `,"Name":"PLS"}}`),
			Token:        c.Config.Users.Admin.Token,
			RespContains: []string{`"ConventionType"`, `"Name":"PLS"`},
			StatusCode:   http.StatusOK}, // 4 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.PUT("/api/convention_type").WithBytes(tc.Sent).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "UpdateConventionType") {
		t.Error(r)
	}
}

// testGetConventionTypes checks route is protected and all housing conventions
// are correctly sent back
func testGetConventionTypes(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		*c.UserCheckTestCase, // 0 : user unauthorized
		{
			Token:         c.Config.Users.User.Token,
			RespContains:  []string{`"ConventionType"`, `"Name"`},
			Count:         2,
			CountItemName: `"ID"`,
			StatusCode:    http.StatusOK}, // 1 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.GET("/api/convention_types").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "GetConventionTypes") {
		t.Error(r)
	}
}

// testDeleteConventionType checks that route is renew project protected and
// delete request sends ok back
func testDeleteConventionType(t *testing.T, c *TestContext, ID int) {
	tcc := []TestCase{
		*c.AdminCheckTestCase, // 0 : bad token
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           0,
			RespContains: []string{`Suppression de type de convention, requête : Type de convention introuvable`},
			StatusCode:   http.StatusInternalServerError}, // 2 : bad ID
		{
			Token:        c.Config.Users.Admin.Token,
			ID:           ID,
			RespContains: []string{`Type de convention supprimé`},
			StatusCode:   http.StatusOK}, // 3 : ok
	}
	f := func(tc TestCase) *httpexpect.Response {
		return c.E.DELETE("/api/convention_type/"+strconv.Itoa(tc.ID)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
	}
	for _, r := range chkFactory(tcc, f, "DeleteConventionType") {
		t.Error(r)
	}
}
